package nodes

import (
	"time"

	"git.ovh.im/yeqing/EdgeCommon/pkg/rpc/pb"
	teaconst "git.ovh.im/yeqing/EdgeNode/internal/const"
	"git.ovh.im/yeqing/EdgeNode/internal/events"
	"git.ovh.im/yeqing/EdgeNode/internal/remotelogs"
	"git.ovh.im/yeqing/EdgeNode/internal/rpc"
	"git.ovh.im/yeqing/EdgeNode/internal/utils/goman"
	"github.com/iwind/TeaGo/Tea"
)

var sharedOCSPTask = NewOCSPUpdateTask()

func init() {
	if !teaconst.IsMain {
		return
	}

	events.On(events.EventLoaded, func() {
		sharedOCSPTask.version = sharedNodeConfig.OCSPVersion

		goman.New(func() {
			sharedOCSPTask.Start()
		})
	})
	events.OnClose(func() {
		sharedOCSPTask.Stop()
	})
}

// OCSPUpdateTask 更新OCSP任务
type OCSPUpdateTask struct {
	version int64

	ticker *time.Ticker
}

func NewOCSPUpdateTask() *OCSPUpdateTask {
	var ticker = time.NewTicker(1 * time.Minute)
	if Tea.IsTesting() {
		ticker = time.NewTicker(10 * time.Second)
	}
	return &OCSPUpdateTask{
		ticker: ticker,
	}
}

func (this *OCSPUpdateTask) Start() {
	for range this.ticker.C {
		err := this.Loop()
		if err != nil {
			if rpc.IsConnError(err) {
				remotelogs.Debug("OCSPUpdateTask", "update ocsp failed: "+err.Error())
			} else {
				remotelogs.Warn("OCSPUpdateTask", "update ocsp failed: "+err.Error())
			}
		}
	}
}

func (this *OCSPUpdateTask) Loop() error {
	rpcClient, err := rpc.SharedRPC()
	if err != nil {
		return err
	}

	resp, err := rpcClient.SSLCertRPC.ListUpdatedSSLCertOCSP(rpcClient.Context(), &pb.ListUpdatedSSLCertOCSPRequest{
		Version: this.version,
		Size:    100,
	})
	if err != nil {
		return err
	}

	for _, ocsp := range resp.SslCertOCSP {
		// 更新OCSP
		if sharedNodeConfig != nil {
			sharedNodeConfig.UpdateCertOCSP(ocsp.SslCertId, ocsp.Data, ocsp.ExpiresAt)
		}

		// 修改版本
		this.version = ocsp.Version
	}

	return nil
}

func (this *OCSPUpdateTask) Stop() {
	this.ticker.Stop()
}
