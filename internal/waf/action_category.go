// Copyright 2021 GoEdge goedge.cdn@gmail.com. All rights reserved.

package waf

import "git.ovh.im/yeqing/EdgeCommon/pkg/serverconfigs/firewallconfigs"

type ActionCategory = string

const (
	ActionCategoryAllow  ActionCategory = firewallconfigs.HTTPFirewallActionCategoryAllow
	ActionCategoryBlock  ActionCategory = firewallconfigs.HTTPFirewallActionCategoryBlock
	ActionCategoryVerify ActionCategory = firewallconfigs.HTTPFirewallActionCategoryVerify
)
