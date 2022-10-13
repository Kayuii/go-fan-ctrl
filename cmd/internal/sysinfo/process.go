package sysinfo

// Process is the generic interface that is implemented on every platform
// and provides common operations for processes.
type Process interface {
	// Pid is the process ID for this process.
	Pid() uint

	// Command name running this process. This is not a path to the
	// executable.
	Command() string

	RunTime() uint64

	UserName() string
}

// FindProcess looks up a single process by pid.
//
// Process will be nil and error will be nil if a matching process is
// not found.
func FindProcess(pid uint) (Process, error) {
	return findProcess(pid)
}
