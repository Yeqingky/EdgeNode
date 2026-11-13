package waf_test

import (
	"testing"

	"git.ovh.im/yeqing/EdgeCommon/pkg/serverconfigs/firewallconfigs"
	"git.ovh.im/yeqing/EdgeNode/internal/waf"
	"github.com/iwind/TeaGo/logs"
)

func TestWAFManager_convert(t *testing.T) {
	p := &firewallconfigs.HTTPFirewallPolicy{
		Id:   1,
		IsOn: true,
		Inbound: &firewallconfigs.HTTPFirewallInboundConfig{
			IsOn: true,
			Groups: []*firewallconfigs.HTTPFirewallRuleGroup{
				{
					Id: 1,
					Sets: []*firewallconfigs.HTTPFirewallRuleSet{
						{
							Id: 1,
						},
						{
							Id: 2,
							Rules: []*firewallconfigs.HTTPFirewallRule{
								{
									Id: 1,
								},
								{
									Id: 2,
								},
							},
						},
					},
				},
			},
		},
	}
	w, err := waf.SharedWAFManager.ConvertWAF(p)
	if err != nil {
		t.Fatal(err)
	}

	logs.PrintAsJSON(w, t)
}
