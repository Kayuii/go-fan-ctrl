//go:build linux
// +build linux

package sysinfo

import (
	"os"
	"testing"
)

func TestUnixProcess_impl(t *testing.T) {
	proc := &UnixProcess{
		pid: uint(os.Getpid()),
	}
	err := proc.Refresh()
	if err != nil {
		t.Fatalf("error should be nil but got: %v", err)
	}
	t.Logf("proc value: %+v", proc)
}
