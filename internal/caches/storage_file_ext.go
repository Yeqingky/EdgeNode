//go:build !plus

package caches

func (this *FileStorage) tryMMAPReader(isPartial bool, estimatedSize int64, path string) (Reader, error) {
	// stub
	return nil, nil
}
