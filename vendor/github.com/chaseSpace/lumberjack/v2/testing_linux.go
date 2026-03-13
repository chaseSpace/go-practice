package lumberjack

import (
	"os"
	"testing"
)

func sameFileMode(path string, mode FileMode, t testing.TB) {
	info, err := os.Stat(path)
	isNilUp(err, t, 1)
	equalsUp(mode.ToOS(), info.Mode(), t, 1)
}
