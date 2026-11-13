package fsutils_test

import (
	"testing"
	"time"

	fsutils "git.ovh.im/yeqing/EdgeNode/internal/utils/fs"
)

func TestWaitLoad(t *testing.T) {
	fsutils.WaitLoad(100, 5, 1*time.Minute)
}
