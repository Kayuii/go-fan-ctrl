//go:build linux
// +build linux

package sysinfo

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"strings"
	"time"
)

// UnixProcess is an implementation of Process that contains Unix-specific
// fields and information.
type UnixProcess struct {
	pid       uint
	command   string
	utime     uint64
	stime     uint64
	starttime uint64

	cmdline  string
	username string
}

func (p *UnixProcess) Pid() uint {
	return p.pid
}

func (p *UnixProcess) UserName() string {
	return p.username
}

func (p *UnixProcess) Command() string {
	return strings.Trim(p.command, "()")
}

func (p *UnixProcess) fullCommand() string {
	return p.cmdline
}

func (p *UnixProcess) usedTime() uint64 {
	return p.utime + p.stime
}

func (p *UnixProcess) RunTime() uint64 {
	current_time := uint64(time.Now().Unix())
	boot_time, err := GetBootTime()
	if err != nil {
		return uint64(0)
	}
	return current_time - boot_time - p.starttime/100
}

func (p *UnixProcess) startTime() uint64 {
	return p.starttime
}

/*
   // https://kb.novaordis.com/index.php/Linux_Process_Information
   (1) pid           %d  The process ID
   (2) comm          %s  The filename of the executable, in parentheses.
   (3) state         %c
   (4) ppid          %d  The PID of the parent of this process.
   (5) pgrp          %d  The process group ID of the process.

   (6) session       %d  The session ID of the process.
   (7) tty_nr        %d  The controlling terminal of the process.
   (8) tpgid         %d  The ID of the foreground process group
   (9) flags         %u  The kernel flags word of the process.

   (10) minflt       %lu The number of minor faults the process has made
   (11) cminflt      %lu The number of minor faults that the process's waited-for children have made.
   (12) majflt       %lu The number of major faults the process has made
   (13) cmajflt      %lu The number of major faults that the process's
   (14) utime        %lu Amount of time that this process has been scheduled in user mode
   (15) stime        %lu Amount of time that this process has been scheduled in kernel mode

   (16) cutime       %ld Amount of time that this process's waited-for chil
   (17) cstime       %ld Amount of time that this process's waited-for chil
   (18) priority     %ld
   (19) nice         %ld The nice value (see setpriority(2))
   (20) num_threads  %ld Number of threads in this process
   (21) itrealvalue  %ld The time in jiffies before the next SIGALRM
   (22) starttime    %llu The time the process started after system boot (clock ticks (divide by sysconf(_SC_CLK_TCK)).

   ...
*/
// Refresh reloads all the data associated with this process.
func (p *UnixProcess) Refresh() error {
	p.loadExtendedCommand()
	err := p.loadProcStats()
	if err != nil {
		return err
	}

	// Read /proc/[pid]/status to get the uid, then lookup uid to get username.
	status, err := getProcStatus(p.pid)
	if err != nil {
		return fmt.Errorf("failed to read process status for pid %d: %v", p.pid, err)
	}
	uids, err := getUIDs(status)
	if err != nil {
		return fmt.Errorf("failed to read process status for pid %d: %v", p.pid, err)
	}
	user, err := user.LookupId(uids[0])
	if err == nil {
		p.username = user.Username
	} else {
		p.username = uids[0]
	}

	return nil
}

func (p *UnixProcess) loadProcStats() error {
	statPath := fmt.Sprintf("/proc/%d/stat", p.pid)

	f, err := os.OpenFile(statPath, os.O_RDONLY, 0444)
	if err != nil {
		return err
	}
	defer f.Close()

	var ignore struct {
		s string
		d int64
	}

	reader := bufio.NewReader(f)

	// extract                    1  2  3  4  5  6  7  8  9 10 11 12 13 14 15 16 17 18 19 20 21 22
	_, err = fmt.Fscanf(reader, "%d %s %s %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d",
		&ignore.d,
		&p.command,
		&ignore.s, &ignore.d, &ignore.d, &ignore.d, &ignore.d,
		&ignore.d, &ignore.d, &ignore.d, &ignore.d, &ignore.d, &ignore.d,
		&p.utime,
		&p.stime,
		&ignore.d, &ignore.d, &ignore.d, &ignore.d, &ignore.d, &ignore.d,
		&p.starttime)

	return err
}

func (p *UnixProcess) loadExtendedCommand() {
	statPath := fmt.Sprintf("/proc/%d/cmdline", p.pid)

	b, err := ioutil.ReadFile(statPath) // just pass the file name
	if err != nil {
		p.cmdline = "unkown"
	}
	p.cmdline = string(b)
}

func findProcess(pid uint) (Process, error) {
	dir := fmt.Sprintf("/proc/%d", pid)
	_, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, err
	}

	return newUnixProcess(pid)
}

func newUnixProcess(pid uint) (*UnixProcess, error) {
	p := &UnixProcess{pid: pid}
	return p, p.Refresh()
}

// getProcStatus reads /proc/[pid]/status which contains process status
// information in human readable form.
func getProcStatus(pid uint) (map[string]string, error) {
	status := make(map[string]string, 42)
	statPath := fmt.Sprintf("/proc/%d/status", pid)
	// path := filepath.Join(Procd, strconv.Itoa(pid), "status")
	err := readFile(statPath, func(line string) bool {
		fields := strings.SplitN(line, ":", 2)
		if len(fields) == 2 {
			status[fields[0]] = strings.TrimSpace(fields[1])
		}

		return true
	})
	return status, err
}

// getUIDs reads the "Uid" value from status and splits it into four values --
// real, effective, saved set, and  file system UIDs.
func getUIDs(status map[string]string) ([]string, error) {
	uidLine, ok := status["Uid"]
	if !ok {
		return nil, fmt.Errorf("Uid not found in proc status")
	}

	uidStrs := strings.Fields(uidLine)
	if len(uidStrs) != 4 {
		return nil, fmt.Errorf("Uid line ('%s') did not contain four values", uidLine)
	}

	return uidStrs, nil
}

func readFile(file string, handler func(string) bool) error {
	contents, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	reader := bufio.NewReader(bytes.NewBuffer(contents))

	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if !handler(string(line)) {
			break
		}
	}

	return nil
}
