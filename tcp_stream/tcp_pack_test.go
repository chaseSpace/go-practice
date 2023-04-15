package tcp_stream

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestWrapParsePack(t *testing.T) {
	data := []byte("Hello, Tcp!")
	s, err := wrapPack(data)
	Convey("wrapPack", t, func() {
		So(err, ShouldBeNil)
	})

	Convey("parsePack", t, func() {
		data2, err := parsePack(s)
		So(err, ShouldBeNil)
		So(data2, ShouldResemble, data) // ShouldResemble 用于数组、切片、map和结构体相等
	})
}

// go test -benchmem -bench=^BenchmarkWrapPack$
func BenchmarkWrapPack(b *testing.B) {
	data := []byte("Hello, Tcp!")
	for i := 0; i < b.N; i++ {
		_, _ = parsePack(data)
	}
}

// go test -benchmem -bench=BenchmarkParsePack
func BenchmarkParsePack(b *testing.B) {
	data := []byte("Hello, Tcp!")
	for i := 0; i < b.N; i++ {
		_, _ = parsePack(data)
	}
}
