package xnet_test

import (
	"net"
	"testing"

	"github.com/byteweap/meta/pkg/xnet"
)

func TestIP2Long(t *testing.T) {
	str1 := "218.108.212.34"
	ip := xnet.IP2Long(str1)
	str2 := xnet.Long2IP(ip)

	t.Logf("str format: %s", str1)
	t.Logf("long format: %d", ip)
	t.Logf("str format: %s", str2)

	// Test invalid IP
	invalidIP := xnet.IP2Long("invalid")
	if invalidIP != 0 {
		t.Errorf("IP2Long('invalid') = %d, want 0", invalidIP)
	}

	// Test empty IP
	emptyIP := xnet.IP2Long("")
	if emptyIP != 0 {
		t.Errorf("IP2Long('') = %d, want 0", emptyIP)
	}
}

func TestLong2IP(t *testing.T) {
	tests := []struct {
		input    uint32
		expected string
	}{
		{0, "0.0.0.0"},
		{3232235777, "192.168.1.1"},
		{2130706433, "127.0.0.1"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := xnet.Long2IP(tt.input)
			if result != tt.expected {
				t.Errorf("Long2IP(%d) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestExternalIP(t *testing.T) {
	ip, err := xnet.ExternalIP()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("External IP: %s", ip)
}

func TestInternalIP(t *testing.T) {
	ip, err := xnet.InternalIP()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Internal IP: %s", ip)
}

func TestExtractIP(t *testing.T) {
	addr := &net.TCPAddr{IP: net.ParseIP("192.168.1.1"), Port: 8080}
	ip, err := xnet.ExtractIP(addr)
	if err != nil {
		t.Fatal(err)
	}
	if ip != "192.168.1.1" {
		t.Errorf("ExtractIP() = %s, want 192.168.1.1", ip)
	}
}

func TestExtractPort(t *testing.T) {
	addr := &net.TCPAddr{IP: net.ParseIP("192.168.1.1"), Port: 8080}
	port, err := xnet.ExtractPort(addr)
	if err != nil {
		t.Fatal(err)
	}
	if port != 8080 {
		t.Errorf("ExtractPort() = %d, want 8080", port)
	}
}

func TestAssignRandPort(t *testing.T) {
	port, err := xnet.AssignRandPort()
	if err != nil {
		t.Fatal(err)
	}
	if port == 0 {
		t.Error("AssignRandPort() returned 0")
	}

	// Test with specific IP
	port2, err := xnet.AssignRandPort("127.0.0.1")
	if err != nil {
		t.Fatal(err)
	}
	if port2 == 0 {
		t.Error("AssignRandPort('127.0.0.1') returned 0")
	}
}

func TestFulfillAddr(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{":8080", "0.0.0.0:8080"},
		{"192.168.1.1:8080", "192.168.1.1:8080"},
		{"localhost:8080", "localhost:8080"},
		{"invalid", "invalid"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := xnet.FulfillAddr(tt.input)
			if result != tt.expected {
				t.Errorf("FulfillAddr(%s) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseAddr(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{"empty", "", false},
		{"with port", ":8080", false},
		{"with host and port", "192.168.1.1:8080", false},
		{"with 0 port", ":0", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			listenAddr, exposeAddr, err := xnet.ParseAddr(tt.input)
			if tt.expectError {
				if err == nil {
					t.Error("ParseAddr() should return error")
				}
			} else {
				if err != nil {
					t.Fatalf("ParseAddr() error: %v", err)
				}
				t.Logf("ParseAddr(%s) = %s, %s", tt.input, listenAddr, exposeAddr)
			}
		})
	}
}
