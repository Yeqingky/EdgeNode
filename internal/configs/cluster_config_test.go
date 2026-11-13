package configs_test

import (
	"testing"

	"git.ovh.im/yeqing/EdgeNode/internal/configs"
	"git.ovh.im/yeqing/EdgeNode/internal/utils/testutils"
	"gopkg.in/yaml.v3"
)

func TestLoadClusterConfig(t *testing.T) {
	if !testutils.IsSingleTesting() {
		return
	}

	config, err := configs.LoadClusterConfig()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", config)

	configData, err := yaml.Marshal(config)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(configData))
}
