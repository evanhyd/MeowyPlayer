package assert

import "meowyplayer.com/utility/logger"

// ensure no error occurs
func NoErr(err error) {
	if err != nil {
		logger.Error("unexpected non-nil value", err, 2)
	}
}

// program invariant assertion
func Ensure(condition func() bool) {
	if !condition() {
		logger.Error("assertion failed", nil, 2)
	}
}
