package core

import (
	"crypto/md5"
	"fmt"
	"github.com/samber/lo"
	"net"
	"sort"
	"strings"
	"time"
)

func netDialTest(isUDP bool, addr string, retry int, intvl time.Duration) bool {
	network := "tcp"
	if isUDP {
		network = "udp"
	}
	for i := 0; i < retry; i++ {
		conn, err := net.DialTimeout(network, addr, time.Second*3)
		if err == nil {
			_ = conn.Close()
			return true
		}
		Sdlogger.Debug("dial tcp %s err:%v", addr, err)
		time.Sleep(intvl)
	}
	return false
}

func CalInstanceHash(instances []*ServiceInstance) string {
	if len(instances) == 0 {
		return ""
	}
	addrs := lo.Map(instances, func(item *ServiceInstance, index int) string {
		return item.Addr()
	})
	sort.Strings(addrs)

	hasher := md5.New()
	hasher.Write([]byte(strings.Join(addrs, "")))
	return fmt.Sprintf("%x", hasher.Sum(nil))
}
