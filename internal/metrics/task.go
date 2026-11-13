// Copyright 2024 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cloud .

package metrics

import (
	"time"

	"git.ovh.im/yeqing/EdgeCommon/pkg/serverconfigs"
)

type Task interface {
	Init() error
	Item() *serverconfigs.MetricItemConfig
	SetItem(item *serverconfigs.MetricItemConfig)
	Add(obj MetricInterface)
	InsertStat(stat *Stat) error
	Upload(pauseDuration time.Duration) error
	Start() error
	Stop() error
	Delete() error
	CleanExpired() error
}
