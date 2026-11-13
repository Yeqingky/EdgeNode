package nodes

import (
	"sync/atomic"
	"time"

	"git.ovh.im/yeqing/EdgeCommon/pkg/nodeconfigs"
	teaconst "git.ovh.im/yeqing/EdgeNode/internal/const"
	"git.ovh.im/yeqing/EdgeNode/internal/events"
	"git.ovh.im/yeqing/EdgeNode/internal/monitor"
	"git.ovh.im/yeqing/EdgeNode/internal/utils/goman"
	"github.com/iwind/TeaGo/maps"
)

// 发送监控流量
func init() {
	if !teaconst.IsMain {
		return
	}

	events.On(events.EventStart, func() {
		var ticker = time.NewTicker(1 * time.Minute)
		goman.New(func() {
			for range ticker.C {
				// 加入到数据队列中
				var inBytes = atomic.LoadUint64(&teaconst.InTrafficBytes)
				atomic.StoreUint64(&teaconst.InTrafficBytes, 0) // 重置数据
				if inBytes > 0 {
					monitor.SharedValueQueue.Add(nodeconfigs.NodeValueItemTrafficIn, maps.Map{
						"total": inBytes,
					})
				}

				var outBytes = atomic.LoadUint64(&teaconst.OutTrafficBytes)
				atomic.StoreUint64(&teaconst.OutTrafficBytes, 0) // 重置数据
				if outBytes > 0 {
					monitor.SharedValueQueue.Add(nodeconfigs.NodeValueItemTrafficOut, maps.Map{
						"total": outBytes,
					})
				}
			}
		})
	})
}
