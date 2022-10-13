package cmd

import (
	"github.com/alecthomas/kong"
	"github.com/kayuii/go-fan-ctrl/utils"
	"github.com/kayuii/go-fan-ctrl/version"
)

type Globals struct {
	Debug   bool                `help:"print debug info" short:"d"`
	Version version.VersionFlag `help:"print version and exit" short:"V"`
}

type CLI struct {
	Globals

	Rm     RmCmd     `cmd:"" help:"Remove files."`
	Detect detectCmd `cmd:"" help:"Remove files. " default:"0"`
}

func Execute() {
	cli := CLI{}
	kong.ConfigureHelp(kong.HelpOptions{Compact: false, Summary: true})
	kong.Name("fanctrl")
	kong.Description("Control GPU power and fan speed")

	ctx := kong.Parse(&cli)
	utils.Fatal("Failed to parse input parameters/commands:", ctx.Error)
	err := ctx.Run(&cli.Globals)
	ctx.FatalIfErrorf(err)
}
