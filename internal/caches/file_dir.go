package caches

import "git.ovh.im/yeqing/EdgeCommon/pkg/serverconfigs/shared"

type FileDir struct {
	Path     string
	Capacity *shared.SizeCapacity
	IsFull   bool
}
