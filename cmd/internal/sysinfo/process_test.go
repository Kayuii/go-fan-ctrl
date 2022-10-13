package sysinfo

import (
	"os"
	"testing"

	"github.com/kayuii/go-fan-ctrl/utils"
)

func TestProcess(t *testing.T) {
	proc := &UnixProcess{
		pid: uint(os.Getpid()),
		// pid: uint(1),
	}
	err := proc.Refresh()
	if err != nil {
		t.Fatalf("error should be nil but got: %v", err)
	}

	t.Logf("Pid value: %+v", proc.Pid())
	t.Logf("UserName value: %+v", proc.UserName())
	t.Logf("Command value: %+v", proc.Command())
	t.Logf("UsedTime value: %+v", utils.SecTimeHuman(proc.usedTime()))
	t.Logf("RunTime value: %+v", utils.SecTimeHuman(proc.RunTime()))
	t.Logf("StartTime value: %+v", proc.startTime())
	t.Logf("cmdline value: %+v", proc.fullCommand())

}
