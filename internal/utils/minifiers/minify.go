// Copyright 2024 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cloud .
//go:build !plus

package minifiers

import (
	"net/http"

	"git.ovh.im/yeqing/EdgeCommon/pkg/serverconfigs"
)

// MinifyResponse minify response body
func MinifyResponse(config *serverconfigs.HTTPPageOptimizationConfig, url string, resp *http.Response) error {
	// stub
	return nil
}
