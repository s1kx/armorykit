package process

import (
	"fmt"
	"strings"
)

type Process interface {
	Start() error
	Stop() error

	PrintOutput()
	WaitForShutdown() error
}

const (
	SingleDash = "-"
	DoubleDash = "--"
)

var DefaultFlagDash = DoubleDash

type Flag struct {
	Name  string
	Value string
}

func (f Flag) Dashed(dash string) string {
	name := f.Name
	if !strings.HasPrefix(name, "-") {
		name = fmt.Sprintf("-%s", name)
	}
	return fmt.Sprintf("%s=%s", name, f.Value)
}

func (f Flag) String() string {
	return f.Dashed(DefaultFlagDash)
}

type FlagList []Flag

func (fl FlagList) StringList() []string {
	fs := []string{}
	for _, flag := range fl {
		fs = append(fs, flag.String())
	}
	return fs
}

func (fl FlagList) String() string {
	return strings.Join(fl.StringList(), " ")
}
