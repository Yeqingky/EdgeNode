// Copyright 2022 GoEdge goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cloud .

package clock_test

import (
	"testing"

	"git.ovh.im/yeqing/EdgeNode/internal/utils/clock"
	"git.ovh.im/yeqing/EdgeNode/internal/utils/testutils"
)

func TestReadServer(t *testing.T) {
	if !testutils.IsSingleTesting() {
		return
	}

	t.Log(clock.NewClockManager().ReadServer("pool.ntp.org"))
}
