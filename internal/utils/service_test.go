package utils

import (
	"testing"

	teaconst "git.ovh.im/yeqing/EdgeNode/internal/const"
)

func TestServiceManager_Log(t *testing.T) {
	manager := NewServiceManager(teaconst.ProductName, teaconst.ProductName+" Server")
	manager.Log("Hello, World")
	manager.LogError("Hello, World")
}
