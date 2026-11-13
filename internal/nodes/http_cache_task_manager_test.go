package nodes_test

import (
	"testing"

	"git.ovh.im/yeqing/EdgeCommon/pkg/nodeconfigs"
	"git.ovh.im/yeqing/EdgeNode/internal/caches"
	"git.ovh.im/yeqing/EdgeNode/internal/nodes"
)

func TestHTTPCacheTaskManager_Loop(t *testing.T) {
	// initialize cache policies
	config, err := nodeconfigs.SharedNodeConfig()
	if err != nil {
		t.Fatal(err)
	}
	caches.SharedManager.UpdatePolicies(config.HTTPCachePolicies)

	var manager = nodes.NewHTTPCacheTaskManager()
	err = manager.Loop()
	if err != nil {
		t.Fatal(err)
	}
}
