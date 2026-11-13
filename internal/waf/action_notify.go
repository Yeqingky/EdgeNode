// Copyright 2021 GoEdge goedge.cdn@gmail.com. All rights reserved.

package waf

import (
	"net/http"
	"time"

	"git.ovh.im/yeqing/EdgeCommon/pkg/rpc/pb"
	teaconst "git.ovh.im/yeqing/EdgeNode/internal/const"
	"git.ovh.im/yeqing/EdgeNode/internal/events"
	"git.ovh.im/yeqing/EdgeNode/internal/remotelogs"
	"git.ovh.im/yeqing/EdgeNode/internal/rpc"
	"git.ovh.im/yeqing/EdgeNode/internal/utils/goman"
	"git.ovh.im/yeqing/EdgeNode/internal/waf/requests"
	"github.com/iwind/TeaGo/types"
)

type notifyTask struct {
	ServerId                int64
	HttpFirewallPolicyId    int64
	HttpFirewallRuleGroupId int64
	HttpFirewallRuleSetId   int64
	CreatedAt               int64
}

var notifyChan = make(chan *notifyTask, 128)

func init() {
	if !teaconst.IsMain {
		return
	}

	events.On(events.EventLoaded, func() {
		goman.New(func() {
			rpcClient, err := rpc.SharedRPC()
			if err != nil {
				remotelogs.Error("WAF_NOTIFY_ACTION", "create rpc client failed: "+err.Error())
				return
			}

			for task := range notifyChan {
				_, err = rpcClient.FirewallRPC.NotifyHTTPFirewallEvent(rpcClient.Context(), &pb.NotifyHTTPFirewallEventRequest{
					ServerId:                task.ServerId,
					HttpFirewallPolicyId:    task.HttpFirewallPolicyId,
					HttpFirewallRuleGroupId: task.HttpFirewallRuleGroupId,
					HttpFirewallRuleSetId:   task.HttpFirewallRuleSetId,
					CreatedAt:               task.CreatedAt,
				})
				if err != nil {
					remotelogs.Error("WAF_NOTIFY_ACTION", "notify failed: "+err.Error())
				}
			}
		})
	})
}

type NotifyAction struct {
	BaseAction
}

func (this *NotifyAction) Init(waf *WAF) error {
	return nil
}

func (this *NotifyAction) Code() string {
	return ActionNotify
}

func (this *NotifyAction) IsAttack() bool {
	return false
}

// WillChange determine if the action will change the request
func (this *NotifyAction) WillChange() bool {
	return false
}

// Perform the action
func (this *NotifyAction) Perform(waf *WAF, group *RuleGroup, set *RuleSet, request requests.Request, writer http.ResponseWriter) PerformResult {
	select {
	case notifyChan <- &notifyTask{
		ServerId:                request.WAFServerId(),
		HttpFirewallPolicyId:    types.Int64(waf.Id),
		HttpFirewallRuleGroupId: types.Int64(group.Id),
		HttpFirewallRuleSetId:   types.Int64(set.Id),
		CreatedAt:               time.Now().Unix(),
	}:
	default:

	}

	return PerformResult{
		ContinueRequest: true,
	}
}
