package main

import (
	"errors"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

const (
	defaultConfigPath = "~/.config/armory-kit.toml"
)

var (
	Version   string
	BuildTime string
)

// // appContextFunc is a function that receives the application config in addition to the cli context.
// type appContextFunc func(c *cli.Context, conf *config) error
//
// // Helper function to pass the application config to a cli function.
// func (f appContextFunc) WithConfig(conf *config) func(c *cli.Context) error {
// 	return func(c *cli.Context) error { return f(c, conf) }
// }

func main() {
	app := &cli.App{
		Name:    "armory-kit",
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

			&cli.StringFlag{
				Name:    "profile-key",
				Aliases: []string{"k"},
				EnvVars: []string{"ARMORYKIT_PROFILE"},
				Usage:   "profile `name` from config",
			},
		},
		Before: initApplication,
		Action: launchCmd,
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

	// Configure logger
	logrus.SetFormatter(&logFormatter)
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	// Load configuration
	var conf config
	err := loadConfig(configPath, &conf)
	if err != nil {
		return err
	}

	// Save config in app metadata
	c.App.Metadata["config"] = conf

	return nil
}

func launchCmd(c *cli.Context) error {
	conf := c.App.Metadata["config"].(config)

	// Get profile
	if !c.IsSet("profile-key") {
		return errors.New("Missing --profile-key flag")
	}
	profileKey := c.String("profile-key")
	profile, ok := conf.Profiles[profileKey]
	if !ok {
		return errors.New("invalid profile name, check your config")
	}

	// Create armory instance
	settings := armorySettings{
		Bitcoind: conf.Bitcoind,
		Profile:  profile,
	}
	armory, err := NewArmoryInstance(settings)
	if err != nil {
		logrus.Fatalf("Error creating Armory instance: %s", err)
		return nil
	}

	// Start armory
	if err = armory.Start(); err != nil {
		logrus.Fatalf("Error starting Armory: %s", err)
	}

	// Show output from armory, blocks until execution is over
	armory.PrintOutput()

	if err = armory.WaitForShutdown(); err != nil {
		logrus.Errorf("Armory exited with error: %s", err)
		return nil
	}

	return nil
}
