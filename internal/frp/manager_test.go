package frp

import (
	"encoding/json"
	"testing"

	"git.ovh.im/yeqing/EdgeCommon/pkg/nodeconfigs"
)

func TestBuildConfig(t *testing.T) {
	manager := &Manager{}
	data, err := manager.buildConfig(&nodeconfigs.FRPConfig{
		IsOn:     true,
		BindAddr: "127.0.0.1",
		BindPort: 7001,
		Token:    "secret",
	})
	if err != nil {
		t.Fatal(err)
	}

	var config frpsConfig
	if err := json.Unmarshal(data, &config); err != nil {
		t.Fatal(err)
	}
	if config.BindAddr != "127.0.0.1" || config.BindPort != 7001 {
		t.Fatalf("unexpected bind: %s:%d", config.BindAddr, config.BindPort)
	}
	if config.Auth.Method != "token" || config.Auth.Token != "secret" {
		t.Fatalf("unexpected auth: %#v", config.Auth)
	}
	if config.Log.To != "console" || !config.Log.DisablePrintColor {
		t.Fatalf("unexpected log: %#v", config.Log)
	}
}
