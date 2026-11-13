package iplibrary

import (
	"testing"
	"time"

	"git.ovh.im/yeqing/EdgeCommon/pkg/rpc/pb"
	"git.ovh.im/yeqing/EdgeCommon/pkg/serverconfigs/firewallconfigs"
	"git.ovh.im/yeqing/EdgeNode/internal/utils/testutils"
)

func TestScriptAction_AddItem(t *testing.T) {
	if !testutils.IsSingleTesting() {
		return
	}

	action := NewScriptAction()
	action.config = &firewallconfigs.FirewallActionScriptConfig{
		Path: "/tmp/ip-item.sh",
		Cwd:  "",
		Args: nil,
	}
	err := action.AddItem(IPListTypeBlack, &pb.IPItem{
		Type:      "ipv4",
		Id:        1,
		IpFrom:    "192.168.1.100",
		ExpiredAt: time.Now().Unix(),
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ok")
}

func TestScriptAction_DeleteItem(t *testing.T) {
	if !testutils.IsSingleTesting() {
		return
	}

	action := NewScriptAction()
	action.config = &firewallconfigs.FirewallActionScriptConfig{
		Path: "/tmp/ip-item.sh",
		Cwd:  "",
		Args: nil,
	}
	err := action.DeleteItem(IPListTypeBlack, &pb.IPItem{
		Type:      "ipv4",
		Id:        1,
		IpFrom:    "192.168.1.100",
		ExpiredAt: time.Now().Unix(),
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ok")
}
