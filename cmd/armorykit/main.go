package main

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/s1kx/armorykit/config"
)

const (
	defaultConfigPath = "~/.config/armorykit.toml"
)

var (
	Version   string
	BuildTime string
)

// conf is a local package variable for access to the config from all commands
var conf config.Config

func main() {
	app := &cli.App{
		Name:    "armorykit",
		Usage:   "utility for creating and managing profiles to use with armory",
		Version: Version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				EnvVars: []string{"ARMORYKIT_CONFIG"},
				Value:   defaultConfigPath,
				Usage:   "path to configuration `file`",
			},
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"d"},
				EnvVars: []string{"ARMORYKIT_DEBUG"},
				Usage:   "debug mode",
			},
		},
		Before: initApplication,
		Commands: []*cli.Command{
			launchCmd,
		},
	}
	app.Run(os.Args)
}

var logFormatter = logrus.TextFormatter{
	FullTimestamp:   true,
	TimestampFormat: "2006-01-02 15:04:05",
}

func initApplication(c *cli.Context) error {
	configPath := c.String("config")
	debug := c.Bool("debug")

	// Configure logger.
	logrus.SetFormatter(&logFormatter)
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	// Load configuration in to package variable conf.
	err := config.Load(configPath, &conf)
	if err != nil {
		return err
	}

	return nil
}
