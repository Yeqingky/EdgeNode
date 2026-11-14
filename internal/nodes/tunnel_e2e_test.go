package nodes

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"git.ovh.im/yeqing/EdgeCommon/pkg/nodeconfigs"
	"git.ovh.im/yeqing/EdgeCommon/pkg/serverconfigs"
	"git.ovh.im/yeqing/EdgeNode/internal/frp"
	"github.com/iwind/TeaGo/Tea"
)

// TestTunnelEndToEnd uses real frps/frpc binaries. Set both paths to run it.
func TestTunnelEndToEnd(t *testing.T) {
	frpsBinary, frpcBinary := os.Getenv("FRP_E2E_FRPS"), os.Getenv("FRP_E2E_FRPC")
	if frpsBinary == "" || frpcBinary == "" {
		t.Skip("FRP_E2E_FRPS and FRP_E2E_FRPC are required")
	}
	for _, name := range []string{frpsBinary, frpcBinary} {
		if !filepath.IsAbs(name) {
			t.Fatalf("binary path must be absolute: %s", name)
		}
	}

	root := t.TempDir()
	oldRoot, oldConfigDir := Tea.Root, Tea.ConfigDir()
	Tea.Root = root
	Tea.SetConfigDir(filepath.Join(root, "configs"))
	t.Cleanup(func() {
		Tea.Root = oldRoot
		Tea.SetConfigDir(oldConfigDir)
	})
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	for name, binary := range map[string]string{"frps": frpsBinary, "frpc": frpcBinary} {
		link := filepath.Join(filepath.Dir(exe), name)
		if err := os.Symlink(binary, link); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = os.Remove(link) })
	}

	freePort := func() int {
		listener, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			t.Fatal(err)
		}
		port := listener.Addr().(*net.TCPAddr).Port
		_ = listener.Close()
		return port
	}
	localPort, remotePort := freePort(), freePort()
	for localPort == remotePort {
		remotePort = freePort()
	}
	clusterToken, tunnelSecret := "cluster-test-token", "tunnel-test-secret"
	writeConfig := func(name string, config any) string {
		data, err := json.Marshal(config)
		if err != nil {
			t.Fatal(err)
		}
		path := filepath.Join(root, name)
		if err := os.WriteFile(path, data, 0600); err != nil {
			t.Fatal(err)
		}
		return path
	}
	start := func(binary string, config string) *exec.Cmd {
		cmd := exec.Command(binary, "-c", config)
		cmd.Stdout, cmd.Stderr = io.Discard, io.Discard
		if err := cmd.Start(); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = cmd.Process.Kill(); _ = cmd.Wait() })
		return cmd
	}
	remoteServer := writeConfig("remote-frps.json", map[string]any{
		"bindAddr": "127.0.0.1", "bindPort": remotePort,
		"auth": map[string]any{"method": "token", "token": clusterToken},
	})
	start(frpsBinary, remoteServer)
	frp.SharedManager.Apply(&nodeconfigs.FRPConfig{
		IsOn: true, BindAddr: "127.0.0.1", BindPort: localPort, Token: clusterToken,
		ServerAddrs: []nodeconfigs.FRPServerAddr{{Addr: "127.0.0.1", Port: remotePort}},
		Tunnels:     []nodeconfigs.FRPTunnelConfig{{ID: 7, Secret: tunnelSecret}},
	})
	t.Cleanup(frp.SharedManager.Stop)
	proxyConfig := writeConfig("remote-frpc.json", map[string]any{
		"serverAddr": "127.0.0.1", "serverPort": remotePort, "loginFailExit": false,
		"auth": map[string]any{"method": "token", "token": clusterToken},
		"proxies": []map[string]any{{
			"name": "tunnel-7", "type": "stcp", "secretKey": tunnelSecret,
			"plugin":    map[string]any{"type": "http_proxy"},
			"transport": map[string]any{"useEncryption": true},
		}},
	})
	start(frpcBinary, proxyConfig)

	originServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, "tunnel-origin:"+r.URL.Path)
	}))
	defer originServer.Close()
	host, port, err := net.SplitHostPort(strings.TrimPrefix(originServer.URL, "http://"))
	if err != nil {
		t.Fatal(err)
	}
	origin := &serverconfigs.OriginConfig{Id: 1, IsOn: true, Addr: &serverconfigs.NetworkAddressConfig{
		Protocol: serverconfigs.ProtocolHTTP, Host: host, PortRange: port, TunnelNodeId: 7,
	}}
	if err := origin.Init(context.Background()); err != nil {
		t.Fatal(err)
	}
	client := &http.Client{Timeout: 3 * time.Second, Transport: &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialOriginContext(ctx, origin, network, addr, 2*time.Second)
		},
	}}
	t.Cleanup(client.CloseIdleConnections)
	deadline := time.Now().Add(25 * time.Second)
	var lastErr error
	for time.Now().Before(deadline) {
		var resp *http.Response
		resp, lastErr = client.Get(originServer.URL + "/ok")
		if lastErr == nil {
			body, readErr := io.ReadAll(resp.Body)
			_ = resp.Body.Close()
			if readErr == nil && resp.StatusCode == http.StatusOK && string(body) == "tunnel-origin:/ok" {
				lastErr = nil
				break
			}
			lastErr = fmt.Errorf("unexpected origin response %d: %q (%v)", resp.StatusCode, body, readErr)
		}
		time.Sleep(200 * time.Millisecond)
	}
	if lastErr != nil {
		t.Fatalf("remote frps fallback failed: %v", lastErr)
	}
	if paths := frp.SharedManager.TunnelSocketPaths(7); len(paths) != 2 {
		t.Fatalf("expected two visitor sockets, got %v", paths)
	}

	missing := &serverconfigs.OriginConfig{Addr: &serverconfigs.NetworkAddressConfig{TunnelNodeId: 8}}
	if conn, err := dialOriginContext(context.Background(), missing, "tcp", net.JoinHostPort(host, port), time.Second); err == nil {
		conn.Close()
		t.Fatal("unknown tunnel unexpectedly connected directly to origin")
	}

	tlsOrigin := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusNoContent) }))
	defer tlsOrigin.Close()
	tlsHost, tlsPort, err := net.SplitHostPort(strings.TrimPrefix(tlsOrigin.URL, "https://"))
	if err != nil {
		t.Fatal(err)
	}
	tlsConfig := &serverconfigs.OriginConfig{Id: 2, IsOn: true, Addr: &serverconfigs.NetworkAddressConfig{
		Protocol: serverconfigs.ProtocolHTTPS, Host: tlsHost, PortRange: tlsPort, TunnelNodeId: 7,
	}}
	if err := tlsConfig.Init(context.Background()); err != nil {
		t.Fatal(err)
	}
	conn, _, err := OriginConnect(tlsConfig, 0, "", tlsHost)
	if err != nil {
		t.Fatalf("TLS origin through tunnel: %v", err)
	}
	if _, ok := conn.(*tls.Conn); !ok {
		t.Fatal("TLS origin did not return a TLS connection")
	}
	_ = conn.Close()
}
