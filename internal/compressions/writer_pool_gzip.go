// Copyright 2022 GoEdge goedge.cdn@gmail.com. All rights reserved.

package compressions

import (
	"io"

	teaconst "git.ovh.im/yeqing/EdgeNode/internal/const"
)

var sharedGzipWriterPool *WriterPool

func init() {
	if !teaconst.IsMain {
		return
	}

	sharedGzipWriterPool = NewWriterPool(CalculatePoolSize(), func(writer io.Writer, level int) (Writer, error) {
		return newGzipWriter(writer)
	})
}
