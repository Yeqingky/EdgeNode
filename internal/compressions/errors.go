package compressions

import "errors"

var ErrIsBusy = errors.New("the system is busy for compression")

func CanIgnore(err error) bool {
	if err == nil {
		return true
	}
	return errors.Is(err, ErrIsBusy)
}
