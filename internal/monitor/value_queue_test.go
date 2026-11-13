package monitor

import (
	"testing"

	"git.ovh.im/yeqing/EdgeCommon/pkg/rpc/pb"
	"git.ovh.im/yeqing/EdgeNode/internal/rpc"
	"git.ovh.im/yeqing/EdgeNode/internal/utils/testutils"
	_ "github.com/iwind/TeaGo/bootstrap"
	"github.com/iwind/TeaGo/logs"
	"google.golang.org/grpc/status"
)

func TestValueQueue_RPC(t *testing.T) {
	if !testutils.IsSingleTesting() {
		return
	}

	rpcClient, err := rpc.SharedRPC()
	if err != nil {
		t.Fatal(err)
	}
	_, err = rpcClient.NodeValueRPC.CreateNodeValue(rpcClient.Context(), &pb.CreateNodeValueRequest{})
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {
			logs.Println(statusErr.Code())
		}
		t.Fatal(err)
	}
	t.Log("ok")
}
