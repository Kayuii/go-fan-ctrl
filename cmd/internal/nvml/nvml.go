package nvml

import (
	"errors"
	"log"

	"github.com/NVIDIA/go-nvml/pkg/nvml"
)

type Device struct {
	DeviceName string
	DeviceUUID string
	d          *nvml.Device
	i          int
}

func newDevice(nvmlDevice *nvml.Device, idx int) (dev Device, err nvml.Return) {
	dev = Device{
		d: nvmlDevice,
		i: idx,
	}

	if dev.DeviceUUID, err = dev.UUID(); err != nvml.SUCCESS {
		return
	}

	if dev.DeviceName, err = dev.Name(); err != nvml.SUCCESS {
		return
	}
	return
}

// UUID returns the Device's Unique ID
func (s *Device) UUID() (uuid string, err nvml.Return) {
	return s.d.GetUUID()
}

// Name returns the Device's Name and is not guaranteed to exceed 64 characters in length
func (s *Device) Name() (name string, err nvml.Return) {
	return s.d.GetName()
}

// InitNVML initializes NVML
func InitNVML() error {
	if result := nvml.Init(); result != nvml.SUCCESS {
		return errors.New(nvml.ErrorString(result))
	}
	return nil
}

// ShutdownNVML all resources that were created when we initialized
func ShutdownNVML() error {
	if result := nvml.Shutdown(); result != nvml.SUCCESS {
		return errors.New(nvml.ErrorString(result))
	}
	return nil
}

// ret := nvml.Init()
// if ret != nvml.SUCCESS {
// 	// log.Printf("Unable to initialize NVML: %v", nvml.ErrorString(ret))
// 	utils.Fatal("Unable to initialize NVML: %v", errors.New(nvml.ErrorString(ret)))
// }
// defer func() {
// 	ret := nvml.Shutdown()
// 	if ret != nvml.SUCCESS {
// 		// log.Printf("Unable to shutdown NVML: %v", nvml.ErrorString(ret))
// 		utils.Fatal("Unable to shutdown NVML: %v", errors.New(nvml.ErrorString(ret)))
// 	}
// }()

func GetDevices() ([]Device, nvml.Return) {
	devCount, err := nvml.DeviceGetCount()
	if err != nvml.SUCCESS {
		log.Fatalf("Unable to get device count: %v", nvml.ErrorString(err))
	}

	devices := make([]Device, devCount)
	for i := 0; i <= devCount-1; i++ {
		nvdev, err := nvml.DeviceGetHandleByIndex(0)
		if err != nvml.SUCCESS {
			return nil, err
		}
		if devices[i], err = newDevice(&nvdev, i); err != nvml.SUCCESS {
			return nil, err
		}
	}

	return devices, nvml.SUCCESS
}

// GetUtilization returns the GPU and memory usage returned as a percentage used of a given GPU device
func (s *Device) GetUtilization() (gpu, memory int, err nvml.Return) {
	utilRates, err := s.d.GetUtilizationRates()
	if err != nvml.SUCCESS {
		return
	}

	gpu = int(utilRates.Gpu)
	memory = int(utilRates.Memory)
	return
}

// GetPowerUsage returns the power consumption of the GPU in watts
func (s *Device) GetPowerUsage() (int, nvml.Return) {
	usage, err := s.d.GetPowerUsage()
	if err != nvml.SUCCESS {
		return 0, err
	}
	// nvmlDeviceGetPowerUsage returns milliwatts.. convert to watts
	return int(usage) / 1000, nvml.SUCCESS
}

// GetFanSpeed returns the fan speed in percent
func (s *Device) GetFanSpeed() (int, nvml.Return) {
	speed, err := s.d.GetFanSpeed()
	if err != nvml.SUCCESS {
		return 0, err
	}
	return int(speed), nvml.SUCCESS
}

// GetTemperature returns the Device's temperature in Farenheit and celsius
func (s *Device) GetTemperature() (int, int, nvml.Return) {
	tempc, err := s.d.GetTemperature(nvml.TEMPERATURE_GPU)
	if err != nvml.SUCCESS {
		return -1, -1, err
	}

	return int(tempc), int(tempc*9/5 + 32), nvml.SUCCESS
}

// GetMemoryInfo retrieves the amount of used, free and total memory available on the device, in bytes.
func (s *Device) GetMemoryInfo() (memInfo *nvml.Memory, err nvml.Return) {
	res, err := s.d.GetMemoryInfo()
	if err != nvml.SUCCESS {
		return nil, err
	}
	return &res, nvml.SUCCESS
}

// GetProcessInfo retrieves the active proccesses (pid, used gpu memory) running on the device
func (s *Device) GetProcessInfo() (procInfo []nvml.ProcessInfo) {

	// var res []nvml.ProcessInfo
	// var cnt C.uint = 0

	res, err := s.d.GetComputeRunningProcesses()
	if err != nvml.SUCCESS {
		return nil
	}
	return res
}
