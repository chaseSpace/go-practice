package xpath_css

import (
	"testing"
	"time"
)

func TestRunXpath(t *testing.T) {
	//Xpath()
	ipArr := []string{
		// 中国
		//"119.131.198.248",
		//"66.249.65.123", "223.246.46.2", "182.116.244.164",
		//"111.201.146.243", "49.78.109.73",
		// 新疆
		//"202.107.128.1",
		//// 直辖市
		//"101.227.131.220",
		//// 英国
		//"5.100.159.255",
		//// 日本
		//"66.117.31.255",
		// 俄罗斯 IP
		"109.167.134.253",
	}
	for _, ip := range ipArr {
		//XpathIPToAddr(ip)
		//QueryAddrByIP_BaiduAPI(ip)
		QueryAddrByIP_IPCn(ip)
		time.Sleep(time.Millisecond * 300)
	}
}
