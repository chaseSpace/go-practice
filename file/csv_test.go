package file

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"runtime"
	"testing"

	"github.com/parnurzeal/gorequest"
)

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

// 测试从HTTP下载的CSV字节流读取数据
// - 若文件较大如几十M，建议先下载到文件，再从文件逐行读取
func TestCsvWrite(t *testing.T) {
	url := "http://xxx.csv"
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))

	_, bts, err := gorequest.New().Get(url).EndBytes()
	if err != nil {
		log.Fatalf("Failed to make HTTP request: %v", err)
	}
	log.Println("kb", float64(len(bts))/1000)

	// 创建一个CSV reader从响应体读取数据
	reader := csv.NewReader(bytes.NewReader(bts))
	reader.ReuseRecord = true
	reader.TrimLeadingSpace = true

	x := 0
	for {
		records, err2 := reader.Read()
		if err2 == io.EOF || len(records) == 0 {
			break
		}
		if err2 != nil {
			log.Fatalf("Failed to read CSV: %v", err2)
		}
		x++
		fmt.Printf("idx %v -- %+v\n", x, records)
	}

	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
}
