//go:build windows
// +build windows

package sysinfo

import (
	"syscall"
	"time"
)

var getTickCount = syscall.NewLazyDLL("kernel32.dll").NewProc("GetTickCount64")

func getOsUptime() (uint64, error) {
	ret, _, err := getTickCount.Call()
	if errno, ok := err.(syscall.Errno); !ok || errno != 0 {
		return 0, err
	}
	return uint64(ret), nil
}
