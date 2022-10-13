package main

import (
	"fmt"
	"log"

	"github.com/NVIDIA/go-nvml/pkg/nvml"
)

func main() {
	ret := nvml.Init()
	if ret != nvml.SUCCESS {
		log.Fatalf("Unable to initialize NVML: %v", nvml.ErrorString(ret))
	}
	defer func() {
		ret := nvml.Shutdown()
		if ret != nvml.SUCCESS {
			log.Fatalf("Unable to shutdown NVML: %v", nvml.ErrorString(ret))
		}
	}()

	driverVersion, ret := nvml.SystemGetDriverVersion()
	if ret != nvml.SUCCESS {
		log.Fatalf("SystemGetDriverVersion: %v", ret)
	} else {
		fmt.Println("Driver Version:\t", driverVersion)
	}

	nvmlVersion, ret := nvml.SystemGetNVMLVersion()
	if ret != nvml.SUCCESS {
		log.Fatalf("SystemGetNVMLVersion: %v", ret)
	} else {
		fmt.Println("NVML Version:\t", nvmlVersion)
	}

	cudaDriverVersion, ret := nvml.SystemGetCudaDriverVersion_v2()
	if ret != nvml.SUCCESS {
		log.Fatalf("SystemGetCudaDriverVersion: %v", ret)
	} else {
		fmt.Println("Cuda Version:\t", cudaDriverVersion)
	}

	count, ret := nvml.DeviceGetCount()
	if ret != nvml.SUCCESS {
		log.Fatalf("Unable to get device count: %v", nvml.ErrorString(ret))
	}

	for di := 0; di < count; di++ {
		device, ret := nvml.DeviceGetHandleByIndex(di)
		if ret != nvml.SUCCESS {
			log.Fatalf("Unable to get device at index %d: %v", di, nvml.ErrorString(ret))
		}

		name, ret := nvml.DeviceGetName(device)
		if ret != nvml.SUCCESS {
			fmt.Printf("dev.DeviceGetName() error: %v\n", ret)
		} else {
			fmt.Printf("Product name:\t%s\n", name)
		}

		uuid, ret := nvml.DeviceGetUUID(device)
		if ret != nvml.SUCCESS {
			fmt.Printf("dev.DeviceGetUUID() error: %v\n", ret)
		} else {
			fmt.Printf("GPU UUID:\t%s\n", uuid)
		}

		memory, ret := device.GetMemoryInfo()
		if ret != nvml.SUCCESS {
			fmt.Printf("dev.MemoryInfo() error: %v\n", ret)
		} else {
			fmt.Printf("memory.total:\t%v, memory.used: %v\n", memory.Total/(1024*1024), memory.Used/(1024*1024))
		}

		rates, ret := device.GetUtilizationRates()
		if ret != nvml.SUCCESS {
			fmt.Printf("dev.UtilizationRates() error: %v\n", ret)
		} else {
			fmt.Printf("utilization.gpu:\t%v, utilization.memory: %v\n", rates.Gpu, rates.Memory)
		}

		powerDraw, ret := device.GetPowerUsage()
		if ret != nvml.SUCCESS {
			fmt.Printf("dev.PowerUsage() error: %v\n", ret)
		} else {
			fmt.Printf("power.draw:\t%v\n", powerDraw)
		}

		temperature, ret := device.GetTemperature(0)
		if ret != nvml.SUCCESS {
			fmt.Printf("dev.Temperature() error: %v\n", ret)
		} else {
			fmt.Printf("temperature.gpu: %v C\n", temperature)
		}

		fanSpeed, ret := device.GetFanSpeed()
		if ret != nvml.SUCCESS {
			fmt.Printf("dev.FanSpeed() error: %v\n", ret)
		} else {
			fmt.Printf("fan.speed:\t%v%%\n", fanSpeed)
		}

		fanNum, ret := device.GetNumFans()
		if ret != nvml.SUCCESS {
			fmt.Printf("dev.GetNumFans() error: %v\n", ret)
		} else {
			fmt.Printf("fan.num:\t%v\n", fanNum)
		}

		for fn := 0; fn < fanNum; fn++ {
			fanSpeed_v2, ret := device.GetFanSpeed_v2(fn)
			if ret != nvml.SUCCESS {
				fmt.Printf("dev.FanSpeed() error: %v\n", ret)
			} else {
				fmt.Printf("fan[%d].speed:\t%v%%\n", fn, fanSpeed_v2)
			}
		}

		encoderUtilization, _, ret := device.GetEncoderUtilization()
		if ret != nvml.SUCCESS {
			fmt.Printf("dev.EncoderUtilization() error: %v\n", ret)
		} else {
			fmt.Printf("utilization.encoder:\t%d\n", encoderUtilization)
		}

		decoderUtilization, _, ret := device.GetDecoderUtilization()
		if ret != nvml.SUCCESS {
			fmt.Printf("dev.DecoderUtilization() error: %v\n", ret)
		} else {
			fmt.Printf("utilization.decoder:\t%d\n", decoderUtilization)
		}

		processInfos, ret := device.GetComputeRunningProcesses()
		if ret != nvml.SUCCESS {
			log.Fatalf("Unable to get process info for device at index %d: %v", di, nvml.ErrorString(ret))
		}
		fmt.Printf("Found %d processes on device %d\n", len(processInfos), di)
		for pi, processInfo := range processInfos {
			fmt.Printf("[%2d] ProcessInfo: %+v\n", pi, processInfo)
		}
	}

}
