package netutil

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsIPV4(t *testing.T) {
	t.Run("should be ipv4", func(t *testing.T) {
		got := IsIPV4("192.168.0.102")
		assert.True(t, got)
	})
	t.Run("should not be ipv4", func(t *testing.T) {
		got := IsIPV4("fe80::7831:3c37:fc14:2348%20")
		assert.False(t, got)

		got = IsIPV4("fe80")
		assert.False(t, got)
	})
}

func TestLocalIP(t *testing.T) {
	got, err := LocalIP()
	assert.NoError(t, err)
	assert.True(t, len(got) > 0)
}

func TestClientIP(t *testing.T) {
	realIp := http.Header{}
	realIp.Set("X-Real-IP", " 10.10.10.10  ")
	forward := http.Header{}
	forward.Set("X-Forwarded-For", "  20.20.20.20, 30.30.30.30")
	var tests = []struct {
		req  *http.Request
		want string
	}{
		{&http.Request{Header: realIp}, "10.10.10.10"},
		{&http.Request{Header: forward}, "20.20.20.20"},
		{&http.Request{Header: http.Header{}, RemoteAddr: "  50.50.50.50:55555 "}, "50.50.50.50"},
	}

	for _, v := range tests {
		got := ClientIP(v.req)
		assert.Equal(t, v.want, got)
	}
}

func TestIPString2Uint(t *testing.T) {
	for _, v := range []struct {
		ip   string
		want uint
	}{
		{"127.0.0.1", 2130706433},
		{"0.0.0.0", 0},
		{"255.255.255.255", 4294967295},
		{"192.168.1.1", 3232235777},
		{"16.187.191.122", 280739706},
	} {
		got := IPString2Uint(v.ip)
		assert.Equal(t, v.want, got)
	}
}

func TestUint2IPString(t *testing.T) {
	for _, v := range []struct {
		num  uint
		want string
	}{
		{2130706433, "127.0.0.1"},
		{0, "0.0.0.0"},
		{4294967295, "255.255.255.255"},
		{3232235777, "192.168.1.1"},
		{280739706, "16.187.191.122"},
	} {
		got := Uint2IPString(v.num)
		assert.Equal(t, v.want, got)
	}
}
