package nodes

import (
	"testing"

	"git.ovh.im/yeqing/EdgeCommon/pkg/serverconfigs"
)

func TestListener_Listen(t *testing.T) {
	listener := NewListener()

	group := serverconfigs.NewServerAddressGroup("https://:1234")

	listener.Reload(group)
	err := listener.Listen()
	if err != nil {
		t.Fatal(err)
	}
}
