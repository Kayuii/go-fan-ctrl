package cmd

import (
	"github.com/kayuii/go-fan-ctrl/cmd/internal/nodes"
)

type detectCmd struct {
	Force     bool `help:"Force removal."`
	Recursive bool `help:"Recursively remove files."`
}

func (d *detectCmd) Run(ctx *Globals) error {
	node := nodes.Node{}
	InitNode(&node)
	FetchNode(&node)
	node.Print()
	return nil
}
