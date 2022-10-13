package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/kayuii/go-fan-ctrl/cmd"
	"github.com/kayuii/go-fan-ctrl/cmd/internal/nvml"
	"github.com/kayuii/go-fan-ctrl/utils"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func main() {
	// dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	// utils.Fatal("Unable to find current program execution directory", err)
	// // log.Println(dir)
	// log.Info().Msg(dir)

	if err := nvml.InitNVML(); err != nil {
		utils.Fatal("Failed initializing NVML:", err)
	}
	defer nvml.ShutdownNVML()

	cmd.Execute()
}
