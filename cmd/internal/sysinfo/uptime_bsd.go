//go:build darwin || freebsd || netbsd
// +build darwin freebsd netbsd

package sysinfo

import (
	"fmt"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/unix"
)

func getOsUptime() (uint64, error) {
	tv, err := unix.SysctlTimeval("kern.boottime")
	if err != nil {
		return 0, err
	}
	return uint64(tv.Sec), nil
}
