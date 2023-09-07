package core

import (
	"crypto/sha256"
	"fmt"
	"github.com/samber/lo"
	"net"
	"sort"
	"strings"
	"time"
)

func netDialTest(addr string, retry int, intvl time.Duration) bool {
	for i := 0; i < retry; i++ {
		conn, err := net.DialTimeout("tcp", addr, time.Second*3)
		if err == nil {
			_ = conn.Close()
			return true
		}
		Sdlogger.Debug("dial tcp %s err:%v", addr, err)
		time.Sleep(intvl)
	}
	return false
}

func calInstanceHash(instances []ServiceInstance) string {
	if len(instances) == 0 {
		return ""
	}
	addrs := lo.Map(instances, func(item ServiceInstance, index int) string {
		return item.Addr()
	})
	sort.Strings(addrs)

	hasher := sha256.New()
	hasher.Write([]byte(strings.Join(addrs, "")))
	return fmt.Sprintf("%x", hasher.Sum(nil))
}
