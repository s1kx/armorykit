package main

import (
	"github.com/urfave/cli"

	"github.com/s1kx/armorykit/cmd"
	"github.com/s1kx/armorykit/launcher"
)

var launchCmd = &cli.Command{
	Name:    "launch",
	Aliases: []string{"l"},
	Usage:   "launch a profile",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "profile",
			Aliases: []string{"k"},
			EnvVars: []string{"ARMORYKIT_PROFILE"},
			Usage:   "profile `key` from config",
		},
	},
	Action: runLaunchCmd,
}

func runLaunchCmd(ctx *cli.Context) error {
	// Check required flags
	if err := cmd.RequireFlags(ctx, "profile"); err != nil {
		return err
	}

	// Get profile.
	profileKey := ctx.String("profile")
	profile, err := conf.GetProfile(profileKey)
	if err != nil {
		return err
	}

	// Run launcher.
	l, err := launcher.New(profile)
	if err != nil {
		return err
	}

	err = l.Run()
	if err != nil {
		return err
	}

	return nil
}
