package xnet_test

import (
	"testing"

	"github.com/byteweap/wukong/pkg/xnet"
)

func TestIP2Long(t *testing.T) {
	str1 := "218.108.212.34"
	ip := xnet.IP2Long(str1)
	str2 := xnet.Long2IP(ip)

	t.Logf("str format: %s", str1)
	t.Logf("long format: %d", ip)
	t.Logf("str format: %s", str2)
}

func TestExternalIP(t *testing.T) {
	ip, err := xnet.ExternalIP()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("External IP: %s", ip)
}
