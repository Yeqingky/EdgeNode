package nodes

import (
	"testing"

	"git.ovh.im/yeqing/EdgeNode/internal/utils/testutils"
)

func TestAPIStream_Start(t *testing.T) {
	if !testutils.IsSingleTesting() {
		return
	}

	apiStream := NewAPIStream()
	apiStream.Start()
}
