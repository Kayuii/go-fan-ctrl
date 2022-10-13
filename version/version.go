package version

import (
	"fmt"
	"time"

	"github.com/alecthomas/kong"
)

var (
	BinaryVersion string = "0.0.0"
	GitLastLog    string
)

type VersionFlag string

func (v VersionFlag) Decode(ctx *kong.DecodeContext) error { return nil }
func (v VersionFlag) IsBool() bool                         { return true }
func (v VersionFlag) BeforeApply(app *kong.Kong, vars kong.Vars) error {
	fmt.Println("Version:", BinaryVersion)
	fmt.Println("GitLastLog:", GitLastLog)
	fmt.Println(time.Now().Format("Mon Jan 2 15:04:05 2006") + " (http://github.com/kayuii/go-fan-ctrl)")
	app.Exit(0)
	return nil
}
