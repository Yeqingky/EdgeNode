package frp

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"

	"git.ovh.im/yeqing/EdgeCommon/pkg/nodeconfigs"
	"git.ovh.im/yeqing/EdgeNode/internal/events"
	"git.ovh.im/yeqing/EdgeNode/internal/remotelogs"
	"github.com/iwind/TeaGo/Tea"
)

type frpsConfig struct {
	BindAddr string `json:"bindAddr"`
	BindPort int    `json:"bindPort"`
	Auth     struct {
		Method string `json:"method"`
		Token  string `json:"token"`
	} `json:"auth"`
	Log struct {
		To                string `json:"to"`
		Level             string `json:"level"`
		DisablePrintColor bool   `json:"disablePrintColor"`
	} `json:"log"`
}

type Manager struct {
	mu             sync.Mutex
	cmd            *exec.Cmd
	current        string
	missingLogged  bool
	visitorCmd     *exec.Cmd
	visitorCurrent string
	visitorDone    chan struct{}
	visitorSockets []string
	visitorPaths   map[int64][]string
}

var SharedManager = &Manager{}

func init() {
	events.OnKey(events.EventQuit, SharedManager, func() {
		SharedManager.Stop()
	})
}

// Apply starts, stops or restarts frps according to the cluster configuration.
func (m *Manager) Apply(config *nodeconfigs.FRPConfig) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if config == nil || !config.IsOn {
		m.stopVisitorLocked()
		m.stopLocked()
		return
	}

	configJSON, err := m.buildConfig(config)
	if err != nil {
		remotelogs.Error("FRP", "build frps config failed: "+err.Error())
		return
	}

	if m.cmd != nil && m.current == string(configJSON) {
		m.applyVisitorsLocked(config)
		return
	}

	binaryPath := m.binaryPath()
	if _, err := os.Stat(binaryPath); err != nil {
		if !m.missingLogged {
			remotelogs.Error("FRP", "frps binary not found: "+binaryPath)
			m.missingLogged = true
		}
		m.stopVisitorLocked()
		m.stopLocked()
		return
	}
	m.missingLogged = false

	configPath := Tea.ConfigFile("frps.json")
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		remotelogs.Error("FRP", "create frps config dir failed: "+err.Error())
		return
	}
	if err := os.WriteFile(configPath, configJSON, 0600); err != nil {
		remotelogs.Error("FRP", "write frps config failed: "+err.Error())
		return
	}

	m.stopVisitorLocked()
	m.stopLocked()

	cmd := exec.Command(binaryPath, "-c", configPath)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		remotelogs.Error("FRP", "open frps stdout failed: "+err.Error())
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		remotelogs.Error("FRP", "open frps stderr failed: "+err.Error())
		return
	}
	if err := cmd.Start(); err != nil {
		remotelogs.Error("FRP", "start frps failed: "+err.Error())
		return
	}

	m.cmd = cmd
	m.current = string(configJSON)
	go m.readLogs("FRPS", stdout)
	go m.readLogs("FRPS", stderr)
	go m.wait(cmd)
	m.applyVisitorsLocked(config)
	remotelogs.Println("FRP", "frps started, pid: "+strconv.Itoa(cmd.Process.Pid))
}

func (m *Manager) Stop() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.stopVisitorLocked()
	m.stopLocked()
}

func (m *Manager) stopLocked() {
	if m.cmd == nil {
		return
	}
	cmd := m.cmd
	m.cmd = nil
	m.current = ""
	_ = cmd.Process.Kill()
	remotelogs.Println("FRP", "frps stopped")
}

func (m *Manager) wait(cmd *exec.Cmd) {
	err := cmd.Wait()
	m.mu.Lock()
	expected := m.cmd != cmd
	if !expected {
		m.cmd = nil
		m.current = ""
	}
	m.mu.Unlock()
	if err != nil && !expected {
		remotelogs.Error("FRP", "frps exited: "+err.Error())
	}
}

func (m *Manager) readLogs(tag string, reader io.Reader) {
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 4096), 256*1024)
	for scanner.Scan() {
		remotelogs.Println(tag, scanner.Text())
	}
}

func (m *Manager) binaryPath() string {
	exe, err := os.Executable()
	if err != nil {
		return "frps"
	}
	name := "frps"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	return filepath.Join(filepath.Dir(exe), name)
}

func (m *Manager) buildConfig(config *nodeconfigs.FRPConfig) ([]byte, error) {
	if err := config.Init(); err != nil {
		return nil, err
	}
	var result frpsConfig
	result.BindAddr = config.BindAddr
	result.BindPort = config.BindPort
	result.Auth.Method = "token"
	result.Auth.Token = config.Token
	result.Log.To = "console"
	result.Log.Level = "info"
	result.Log.DisablePrintColor = true
	return json.Marshal(result)
}
