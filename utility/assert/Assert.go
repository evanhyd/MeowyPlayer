package assert

import "meowyplayer.com/utility/logger"

// ensure no error occurs
func NoErr(err error, message string) {
	if err != nil {
		logger.Error(err, message, 2)
	}
}

// program invariant assertion
func Ensure(condition func() bool) {
	if !condition() {
		logger.Error(nil, "assertion failed", 2)
	}
}
