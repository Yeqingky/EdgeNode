package compressions

import (
	"io"

	teaconst "git.ovh.im/yeqing/EdgeNode/internal/const"
)

var sharedZSTDWriterPool *WriterPool

func init() {
	if !teaconst.IsMain {
		return
	}

	sharedZSTDWriterPool = NewWriterPool(CalculatePoolSize(), func(writer io.Writer, level int) (Writer, error) {
		return newZSTDWriter(writer)
	})
}
