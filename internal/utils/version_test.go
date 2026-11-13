package utils_test

import (
	"testing"

	"git.ovh.im/yeqing/EdgeNode/internal/utils"
)

func TestVersionToLong(t *testing.T) {
	for _, v := range []string{
		"",
		"a",
		"1",
		"1.2",
		"1.2.1",
		"1.2.1.4",
		"1.2.3.4.5",
	} {
		t.Log(v, "=>", utils.VersionToLong(v))
	}
}
