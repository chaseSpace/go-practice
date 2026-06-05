//go:build !linux
// +build !linux

package lumberjack

import (
	"testing"
)

// In Windows, the file mode is not checked is due to Windows's different file mode implements.
func sameFileMode(path string, mode FileMode, t testing.TB) {
}
