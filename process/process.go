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
	var flag string
	if strings.HasPrefix(f.Name, "-") {
		flag = f.Name
	} else {
		flag = fmt.Sprintf("-%s", f.Name)
	}

	// Return flag as -flag if it's value is empty, otherwise -flag=value
	if f.Value != "" {
		flag = fmt.Sprintf("%s=%s", flag, f.Value)
	}
	return flag
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
