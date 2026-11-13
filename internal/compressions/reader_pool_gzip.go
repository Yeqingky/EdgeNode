// Copyright 2022 GoEdge goedge.cdn@gmail.com. All rights reserved.

package compressions

import (
	"io"

	teaconst "git.ovh.im/yeqing/EdgeNode/internal/const"
)

var sharedGzipReaderPool *ReaderPool

func init() {
	if !teaconst.IsMain {
		return
	}

	sharedGzipReaderPool = NewReaderPool(CalculatePoolSize(), func(reader io.Reader) (Reader, error) {
		return newGzipReader(reader)
	})
}
