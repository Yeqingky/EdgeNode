package memutils_test

import (
	"testing"

	"git.ovh.im/yeqing/EdgeNode/internal/utils/mem"
)

func TestSystemMemoryGB(t *testing.T) {
	t.Log(memutils.SystemMemoryGB())
	t.Log(memutils.SystemMemoryGB())
	t.Log(memutils.SystemMemoryGB())
	t.Log(memutils.SystemMemoryBytes())
	t.Log(memutils.SystemMemoryBytes())
	t.Log(memutils.SystemMemoryBytes()>>30, "GB")
	t.Log("available:", memutils.AvailableMemoryGB())
}
