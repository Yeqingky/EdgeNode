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
