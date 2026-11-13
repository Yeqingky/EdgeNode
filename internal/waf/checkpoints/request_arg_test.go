package checkpoints

import (
	"net/http"
	"testing"

	"git.ovh.im/yeqing/EdgeNode/internal/waf/requests"
)

func TestArgParam_RequestValue(t *testing.T) {
	rawReq, err := http.NewRequest(http.MethodGet, "http://teaos.cn/?name=lu", nil)
	if err != nil {
		t.Fatal(err)
	}

	req := requests.NewTestRequest(rawReq)

	checkpoint := new(RequestArgCheckpoint)
	t.Log(checkpoint.RequestValue(req, "name", nil, 1))
	t.Log(checkpoint.ResponseValue(req, nil, "name", nil, 1))
	t.Log(checkpoint.RequestValue(req, "name2", nil, 1))
}
