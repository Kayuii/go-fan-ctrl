package cmd

import (
	"fmt"
	"os"

	"github.com/kayuii/go-fan-ctrl/version"
)

type RmCmd struct {
	Force     bool `help:"Force removal."`
	Recursive bool `help:"Recursively remove files."`

	Paths []string `arg:"" name:"path" help:"Paths to remove." type:"path"`
}

func (r *RmCmd) Run(ctx *Globals) error {

	fmt.Println("agent start")
	fmt.Println("pid: ", os.Getpid())
	fmt.Println("agent version: ", version.BinaryVersion)

	// currentUser, _ := user.Current()
	// fmt.Println("current user userName: ", currentUser.Username)

	return nil
}
