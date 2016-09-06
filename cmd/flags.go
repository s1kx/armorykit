package cmd

import (
	"fmt"
	"strings"

	"github.com/urfave/cli"
)

type MissingFlagsError struct {
	RequiredFlags []string
	MissingFlags  []string
}

func (e MissingFlagsError) Error() string {
	return fmt.Sprintf("missing required flags: %s",
		strings.Join(e.MissingFlags, ", "))
}

func RequireFlags(c *cli.Context, flags ...string) error {
	missing := []string{}
	for _, flag := range flags {
		if !c.IsSet(flag) {
			missing = append(missing, flag)
		}
	}
	if len(missing) > 0 {
		err := &MissingFlagsError{
			RequiredFlags: flags,
			MissingFlags:  missing,
		}
		return err
	}
	return nil
}
