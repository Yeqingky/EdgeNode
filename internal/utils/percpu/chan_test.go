package percpu_test

import (
	"testing"

	"git.ovh.im/yeqing/EdgeNode/internal/utils/percpu"
	"git.ovh.im/yeqing/EdgeNode/internal/utils/zero"
)

func TestChan_C(t *testing.T) {
	var c = percpu.NewChan[zero.Zero](10)
	c.C() <- zero.Zero{}

	t.Log(<-c.C())

	select {
	case <-c.C():
		t.Fatal("should not return from here")
	default:
		t.Log("ok")
	}
}
