//go:build linux
// +build linux

package sysinfo

/*
#include <unistd.h>
*/
import "C"

func GetClockTick() uint64 {
	return uint64(C.sysconf(C._SC_CLK_TCK))
}
