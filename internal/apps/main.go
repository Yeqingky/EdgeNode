package apps

import teaconst "git.ovh.im/yeqing/EdgeNode/internal/const"

func RunMain(f func()) {
	if teaconst.IsMain {
		f()
	}
}
