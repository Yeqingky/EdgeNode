package utils

import (
	"os"

	"git.ovh.im/yeqing/EdgeNode/internal/events"
)

func Exit() {
	events.Notify(events.EventTerminated)
	os.Exit(0)
}
