package ip

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
)

func TestIp2region(t *testing.T) {
	f, _ := os.ReadFile("./ip2region.xdb")

	// 2、用全局的 cBuff 创建完全基于内存的查询对象。
	searcher, err := xdb.NewWithBuffer(f)
	if err != nil {
		fmt.Printf("failed to create searcher with content: %s\n", err)
		return
	}

	var ips = []string{
		"123.144.103.155",
		"123.144.51.35",
		"1.66.5.142",
	}

	for _, ip := range ips {
		var tStart = time.Now()

		// 单次查询效率在微秒级别，searcher可以并发使用
		region, err := searcher.SearchByStr(ip)
		if err != nil {
			fmt.Printf("failed to SearchIP(%s): %s\n", ip, err)
			return
		}
		/*
			region: 中国|0|重庆|重庆市|联通, took: 0ns
			region: 中国|0|重庆|重庆市|联通, took: 0ns
			region: 日本|0|0|0|0, took: 0ns
		*/
		fmt.Printf("region: %s, took: %dns\n", region, time.Since(tStart).Nanoseconds())
	}
}
