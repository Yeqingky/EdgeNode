package checkpoints

import (
	"net/http"
	"testing"

	"git.ovh.im/yeqing/EdgeNode/internal/waf/requests"
)

func TestResponseHeaderCheckpoint_ResponseValue(t *testing.T) {
	resp := requests.NewResponse(new(http.Response))
	resp.StatusCode = 200
	resp.Header = http.Header{}
	resp.Header.Set("Hello", "World")

	checkpoint := new(ResponseHeaderCheckpoint)
	t.Log(checkpoint.ResponseValue(nil, resp, "Hello", nil, 1))
}
