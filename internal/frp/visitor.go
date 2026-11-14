package frp

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"time"

	"git.ovh.im/yeqing/EdgeCommon/pkg/nodeconfigs"
	"git.ovh.im/yeqing/EdgeNode/internal/remotelogs"
	"github.com/iwind/TeaGo/Tea"
)

// TunnelSocketPath is the private local endpoint for a tunnel visitor.
func TunnelSocketPath(id int64) string {
	return filepath.Join(filepath.Dir(Tea.ConfigFile("frps.json")), "tunnels", strconv.FormatInt(id, 10)+".sock")
}

func (m *Manager) TunnelSocketPaths(id int64) []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return slices.Clone(m.visitorPaths[id])
}

func (m *Manager) applyVisitorsLocked(config *nodeconfigs.FRPConfig) {
	if len(config.Tunnels) == 0 {
		m.stopVisitorLocked()
		return
	}

	serverAddr := config.BindAddr
	switch serverAddr {
	case "", "0.0.0.0":
		serverAddr = "127.0.0.1"
	case "::", "[::]":
		serverAddr = "::1"
	}
	servers := []nodeconfigs.FRPServerAddr{{Addr: serverAddr, Port: config.BindPort}}
	seen := map[nodeconfigs.FRPServerAddr]bool{servers[0]: true}
	for _, server := range config.ServerAddrs {
		if server.Addr == "" || server.Addr == "0.0.0.0" || server.Addr == "::" || server.Port < 1 || server.Port > 65535 || seen[server] {
			continue
		}
		servers = append(servers, server)
		seen[server] = true
	}

	visitors := make([]map[string]any, 0, len(config.Tunnels))
	sockets := make([]string, 0, len(config.Tunnels)*len(servers))
	paths := make(map[int64][]string, len(config.Tunnels))
	for _, tunnel := range config.Tunnels {
		if tunnel.ID <= 0 || tunnel.Secret == "" {
			remotelogs.Error("FRP", "invalid tunnel visitor configuration")
			return
		}
		path := TunnelSocketPath(tunnel.ID)
		visitors = append(visitors, map[string]any{
			"name":       fmt.Sprintf("visitor-%d", tunnel.ID),
			"type":       "stcp",
			"serverName": fmt.Sprintf("tunnel-%d", tunnel.ID),
			"secretKey":  tunnel.Secret,
			"bindPort":   -1,
			"bindUnix":   path,
			"transport":  map[string]any{"useEncryption": true},
		})
		for index := range servers {
			visitorPath := path
			if len(servers) > 1 {
				visitorPath += "." + strconv.Itoa(index)
			}
			paths[tunnel.ID] = append(paths[tunnel.ID], visitorPath)
			sockets = append(sockets, visitorPath)
		}
	}
	configJSON, err := json.Marshal(map[string]any{
		"serverAddr":    serverAddr,
		"serverPort":    config.BindPort,
		"serverAddrs":   servers,
		"auth":          map[string]any{"method": "token", "token": config.Token},
		"log":           map[string]any{"to": "console", "level": "info", "disablePrintColor": true},
		"loginFailExit": false,
		"visitors":      visitors,
	})
	if err != nil {
		remotelogs.Error("FRP", "build visitor config failed: "+err.Error())
		return
	}
	if m.visitorCmd != nil && m.visitorCurrent == string(configJSON) {
		return
	}

	m.stopVisitorLocked()
	dir := filepath.Dir(sockets[0])
	if err := os.MkdirAll(dir, 0700); err != nil {
		remotelogs.Error("FRP", "create visitor socket directory failed: "+err.Error())
		return
	}
	if err := os.Chmod(dir, 0700); err != nil {
		remotelogs.Error("FRP", "secure visitor socket directory failed: "+err.Error())
		return
	}
	for _, path := range sockets {
		info, err := os.Lstat(path)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil || info.Mode()&os.ModeSocket == 0 {
			remotelogs.Error("FRP", "invalid visitor socket path: "+path)
			return
		}
		if err := os.Remove(path); err != nil {
			remotelogs.Error("FRP", "remove stale visitor socket failed: "+err.Error())
			return
		}
	}

	configPath := Tea.ConfigFile("frpc-visitors.json")
	if err := os.WriteFile(configPath, configJSON, 0600); err != nil {
		remotelogs.Error("FRP", "write visitor config failed: "+err.Error())
		return
	}
	binaryPath := filepath.Join(filepath.Dir(m.binaryPath()), "frpc")
	if filepath.Ext(m.binaryPath()) == ".exe" {
		binaryPath += ".exe"
	}
	cmd := exec.Command(binaryPath, "-c", configPath)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		remotelogs.Error("FRP", "open visitor stdout failed: "+err.Error())
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		remotelogs.Error("FRP", "open visitor stderr failed: "+err.Error())
		return
	}
	if err := cmd.Start(); err != nil {
		remotelogs.Error("FRP", "start visitor failed: "+err.Error())
		return
	}
	m.visitorCmd = cmd
	m.visitorCurrent = string(configJSON)
	m.visitorSockets = sockets
	m.visitorPaths = paths
	m.visitorDone = make(chan struct{})
	go m.readLogs("FRPC", stdout)
	go m.readLogs("FRPC", stderr)
	go m.waitVisitor(cmd, m.visitorDone)
}

func (m *Manager) stopVisitorLocked() {
	m.visitorPaths = nil
	if m.visitorCmd != nil {
		cmd, done := m.visitorCmd, m.visitorDone
		m.visitorCmd = nil
		m.visitorCurrent = ""
		_ = cmd.Process.Kill()
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			remotelogs.Error("FRP", "visitor process did not stop")
			return
		}
	}
	for _, path := range m.visitorSockets {
		if info, err := os.Lstat(path); err == nil && info.Mode()&os.ModeSocket != 0 {
			_ = os.Remove(path)
		}
	}
	m.visitorSockets = nil
}

func (m *Manager) waitVisitor(cmd *exec.Cmd, done chan struct{}) {
	err := cmd.Wait()
	close(done)
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.visitorCmd == cmd {
		m.visitorCmd = nil
		m.visitorCurrent = ""
		m.visitorPaths = nil
		if err != nil {
			remotelogs.Error("FRP", "visitor exited: "+err.Error())
		}
	}
}
