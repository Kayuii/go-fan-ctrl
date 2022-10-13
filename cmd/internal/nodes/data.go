package nodes

import (
	"bytes"
	"fmt"
	"time"

	"github.com/kayuii/go-fan-ctrl/utils"
	"github.com/tomlazar/table"
)

type Memory struct {
	Used       int64 `json:"used"`
	Free       int64 `json:"free"`
	Total      int64 `json:"total"`
	Percentage int   `json:"percentage"`
}

type Process struct {
	Pid             int
	UsedGpuMemory   int64
	Name            string
	Username        string
	RunTime         uint64
	ExtendedCommand string
}

func (p *Process) GetPid() string {
	return fmt.Sprintf("%v", p.Pid)
}

func (p *Process) GetName() string {
	return p.Name
}

func (p *Process) GetGpuMemoryUsage() string {
	return fmt.Sprintf("%3dMiB", p.UsedGpuMemory/1024/1024)
}

func (p *Process) GetRuntime() string {
	return utils.SecTimeHuman(p.RunTime)
}

func (p *Process) GetRunUser() string {
	return p.Username
}
func (p *Process) GetExtendedCommand() string {
	return p.ExtendedCommand
}

func (p *Process) GetProcesses() []string {
	return []string{
		p.GetPid(),
		p.GetRunUser(),
		p.GetName(),
		p.GetGpuMemoryUsage(),
		p.GetRuntime(),
	}
}

type Device struct {
	Id                int    `json:"id"`
	Name              string `json:"name"`
	Utilization       int    `json:"utilization"`
	MemoryUtilization Memory `json:"memory"`
	FanSpeed          int    `json:"fan_speed"`
	Temperature       int    `json:"temperature"`
	PowerUsage        int    `json:"power_usage"`
	Processes         []Process
}

func (d *Device) GetDeviceName() string {
	return fmt.Sprintf("%d:%s", d.Id, d.Name)
}

func (d *Device) GetMemoryInfo() string {
	return fmt.Sprintf("%dMiB / %dMiB (%3d%%)",
		d.MemoryUtilization.Used/1024/1024,
		d.MemoryUtilization.Total/1024/1024,
		int(d.MemoryUtilization.Used*100/d.MemoryUtilization.Total))
}

func (d *Device) GetGpuUtilization() string {
	return fmt.Sprintf("%3d%%", d.Utilization)
}

func (d *Device) GetFanSpeed() string {
	return fmt.Sprintf("%3d%%", d.FanSpeed)
}

func (d *Device) GetTemp() string {
	return fmt.Sprintf("%3dC", d.Temperature)
}

func (d *Device) GetPowerUsage() string {
	return fmt.Sprintf("%3dW", d.PowerUsage)
}

func (d *Device) GetDevices() []string {
	return []string{
		d.GetDeviceName(),
		d.GetFanSpeed(),
		d.GetTemp(),
		d.GetPowerUsage(),
		d.GetMemoryInfo(),
		d.GetGpuUtilization(),
	}
}

func (d *Device) Print() {
	if len(d.Processes) > 0 {
		tableHeader := []string{"PID", "User", "Command", "GPU Mem Usage", ""}
		var tableRows [][]string

		for _, p := range d.Processes {
			tableRows = append(tableRows, p.GetProcesses())
		}

		gpuTable := table.Table{
			Headers: tableHeader,
			Rows:    tableRows,
		}

		var buf bytes.Buffer
		tableErr := gpuTable.WriteTable(&buf, &table.Config{})
		if tableErr != nil {
			utils.Fatal("Error printing table: %v", tableErr)
		}
		fmt.Printf("> %s\n", d.GetDeviceName())
		fmt.Println(buf.String())
	}
}

type Node struct {
	Name       string    `json:"name"`        // hostname
	Devices    []Device  `json:"devices"`     // devices
	Time       time.Time `json:"time"`        // current timestamp from message
	BootTime   int64     `json:"boot_time"`   // uptime of system
	ClockTicks int64     `json:"clock_ticks"` // cpu ticks per second
}

func (n *Node) Print() {
	if len(n.Devices) > 0 {
		tableHeader := []string{"Gpu", "Fan", "Temp", "Power", "Memory-Usage", "GPU-Util"}
		var tableRows [][]string
		for _, d := range n.Devices {
			tableRows = append(tableRows, d.GetDevices())

			d.Print()
		}
		gpuTable := table.Table{
			Headers: tableHeader,
			Rows:    tableRows,
		}
		var buf bytes.Buffer
		tableErr := gpuTable.WriteTable(&buf, &table.Config{})
		if tableErr != nil {
			utils.Fatal("Error printing table: %v", tableErr)
		}
		fmt.Printf("> %s\n", n.Name)
		fmt.Println(buf.String())
	}
}
