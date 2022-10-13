//go:build linux
// +build linux

package sysinfo

import (
	"syscall"
)

func getOsUptime() (uint64, error) {
	info := &syscall.Sysinfo_t{}
	if err := syscall.Sysinfo(info); err != nil {
		return 0, err
	}
	return uint64(info.Uptime), nil
}
