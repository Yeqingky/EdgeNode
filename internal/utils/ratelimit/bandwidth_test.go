package ratelimit_test

import (
	"context"
	"testing"

	"git.ovh.im/yeqing/EdgeNode/internal/utils/ratelimit"
	"git.ovh.im/yeqing/EdgeNode/internal/utils/testutils"
)

func TestBandwidth(t *testing.T) {
	if !testutils.IsSingleTesting() {
		return
	}

	var bandwidth = ratelimit.NewBandwidth(32 << 10)
	bandwidth.Ack(context.Background(), 123)
	bandwidth.Ack(context.Background(), 16<<10)
	bandwidth.Ack(context.Background(), 32<<10)
}

func TestBandwidth_0(t *testing.T) {
	var bandwidth = ratelimit.NewBandwidth(0)
	bandwidth.Ack(context.Background(), 123)
	bandwidth.Ack(context.Background(), 123456)
}
