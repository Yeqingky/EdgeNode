package kvstore

import "github.com/cockroachdb/pebble"

var DefaultWriteOptions = &pebble.WriteOptions{
	Sync: false,
}

var DefaultWriteSyncOptions = &pebble.WriteOptions{
	Sync: true,
}
