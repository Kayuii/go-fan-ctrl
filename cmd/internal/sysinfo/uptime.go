package sysinfo

import (
	"time"
)

// Get uptime duration
func GetUptime() (uint64, error) {
	return getOsUptime()
}

func GetBootTime() (uint64, error) {
	uptime, err := GetUptime()
	if err != nil {
		return uint64(0), err
	}
	currentTime := uint64(time.Now().Unix())

	return currentTime - uptime, nil
}
