package fsutils_test

import (
	"testing"

	fsutils "git.ovh.im/yeqing/EdgeNode/internal/utils/fs"
	"github.com/iwind/TeaGo/assert"
)

func TestFileFlags(t *testing.T) {
	var a = assert.NewAssertion(t)
	a.IsTrue(fsutils.FlagRead&fsutils.FlagRead == fsutils.FlagRead)
	a.IsTrue(fsutils.FlagWrite&fsutils.FlagWrite != fsutils.FlagRead)
	a.IsTrue((fsutils.FlagWrite|fsutils.FlagRead)&fsutils.FlagRead == fsutils.FlagRead)
}
