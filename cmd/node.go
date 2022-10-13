package cmd

import (
	"os"
	"time"

	"github.com/kayuii/go-fan-ctrl/cmd/internal/nodes"
	"github.com/kayuii/go-fan-ctrl/cmd/internal/nvml"
	"github.com/kayuii/go-fan-ctrl/cmd/internal/sysinfo"
)

func InitNode(n *nodes.Node) {
	name, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	n.Name = name
	n.Time = time.Now()

	boot_time, _ := sysinfo.GetBootTime()
	n.BootTime = int64(boot_time)
	// n.ClockTicks = proc.ClockTicks()

	devices, _ := nvml.GetDevices()

	for i := 0; i < len(devices); i++ {
		n.Devices = append(n.Devices, nodes.Device{0, "", 0, nodes.Memory{0, 0, 0, 0}, 0, 0, 0, nil})
	}
}

func FetchNode(n *nodes.Node) {

	devices, _ := nvml.GetDevices()
	n.Time = time.Now()

	boot_time, _ := sysinfo.GetBootTime()
	n.BootTime = int64(boot_time)

	for idx, device := range devices {

		meminfo, _ := device.GetMemoryInfo()
		gpuPercent, _, _ := device.GetUtilization()
		memPercent := int(meminfo.Used / meminfo.Total)
		powerUsage, _ := device.GetPowerUsage()
		fanSpeed, _ := device.GetFanSpeed()
		tempc, _, _ := device.GetTemperature()

		// read processes
		deviceProcs := device.GetProcessInfo()

		// collect al proccess informations
		var processes []nodes.Process

		for i := 0; i < len(deviceProcs); i++ {

			if int(deviceProcs[i].Pid) == 0 {
				continue
			}

			PID := uint(deviceProcs[i].Pid)

			proc_info, _ := sysinfo.FindProcess(PID)

			// extendedCMD := proc.CmdFromPID(PID)
			if deviceProcs != nil {

				processes = append(processes, nodes.Process{
					Pid:           int(PID),
					UsedGpuMemory: int64(deviceProcs[i].UsedGpuMemory),
					Name:          proc_info.Command(),
					Username:      proc_info.UserName(),
					RunTime:       proc_info.RunTime(),
					// ExtendedCommand: extendedCMD,
				})

			}
		}

		n.Devices[idx].Id = idx
		n.Devices[idx].Name = device.DeviceName
		n.Devices[idx].Utilization = gpuPercent
		n.Devices[idx].MemoryUtilization = nodes.Memory{int64(meminfo.Used), int64(meminfo.Free), int64(meminfo.Total), memPercent}
		n.Devices[idx].FanSpeed = fanSpeed
		n.Devices[idx].PowerUsage = powerUsage
		n.Devices[idx].Temperature = tempc
		n.Devices[idx].Processes = processes

	}
}
