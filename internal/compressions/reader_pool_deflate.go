package compressions

import (
	"io"

	teaconst "git.ovh.im/yeqing/EdgeNode/internal/const"
)

var sharedDeflateReaderPool *ReaderPool

func init() {
	if !teaconst.IsMain {
		return
	}

	sharedDeflateReaderPool = NewReaderPool(CalculatePoolSize(), func(reader io.Reader) (Reader, error) {
		return newDeflateReader(reader)
	})
}
