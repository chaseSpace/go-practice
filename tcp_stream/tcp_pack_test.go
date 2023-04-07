package tcp_stream

import (
	"testing"
)

func TestParsePack(t *testing.T) {
	_, err := parsePack([]byte{
		0x1,                // version
		0x0, 0x0, 0x0, 0x3, // data len
		0x11, 0x11, // checksum
		0x1, 0x1, 0x1, // data
	})
	if err != nil {
		t.Error(err)
	}
}
