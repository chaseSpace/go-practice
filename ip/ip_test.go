package ip

import (
	"fmt"
	"os"
	"testing"
)

func TestGetLocalIPs(t *testing.T) {
	ips, err := getLocalPrivateIPs(true, "")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Println("Local IPs:")
	for _, ip := range ips {
		fmt.Println(ip.String())
	}
}
