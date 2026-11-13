package fsutils_test

import (
	"testing"

	fsutils "git.ovh.im/yeqing/EdgeNode/internal/utils/fs"
)

func TestCheckDiskWritingSpeed(t *testing.T) {
	t.Log(fsutils.CheckDiskWritingSpeed())
}

func TestCheckDiskIsFast(t *testing.T) {
	t.Log(fsutils.CheckDiskIsFast())
}
