package bfs

type BlockFileOptions struct {
	BytesPerSync int64
}

func (this *BlockFileOptions) EnsureDefaults() {
	if this.BytesPerSync <= 0 {
		this.BytesPerSync = 1 << 20
	}
}

var DefaultBlockFileOptions = &BlockFileOptions{
	BytesPerSync: 1 << 20,
}
