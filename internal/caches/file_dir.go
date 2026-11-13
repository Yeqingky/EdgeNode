// Copyright 2022 GoEdge goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cloud .

package caches

import "git.ovh.im/yeqing/EdgeCommon/pkg/serverconfigs/shared"

type FileDir struct {
	Path     string
	Capacity *shared.SizeCapacity
	IsFull   bool
}
