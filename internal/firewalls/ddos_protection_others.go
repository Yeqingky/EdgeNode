//go:build !linux
// +build !linux

package firewalls

import (
	"git.ovh.im/yeqing/EdgeCommon/pkg/serverconfigs/ddosconfigs"
)

var SharedDDoSProtectionManager = NewDDoSProtectionManager()

type DDoSProtectionManager struct {
}

func NewDDoSProtectionManager() *DDoSProtectionManager {
	return &DDoSProtectionManager{}
}

func (this *DDoSProtectionManager) Apply(config *ddosconfigs.ProtectionConfig) error {
	return nil
}
