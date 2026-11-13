package utils_test

import (
	"testing"
	"time"

	"git.ovh.im/yeqing/EdgeNode/internal/utils"
)

func TestGMTUnixTime(t *testing.T) {
	t.Log(utils.GMTUnixTime(time.Now().Unix()))
}

func TestGMTTime(t *testing.T) {
	t.Log(utils.GMTTime(time.Now()))
}
