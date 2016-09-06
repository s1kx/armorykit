package launcher

import (
	"github.com/Sirupsen/logrus"

	"github.com/s1kx/armorykit/config"
	"github.com/s1kx/armorykit/process"
)

type Launcher struct {
	profile *config.Profile

	bitcoind *process.ArmoryProcess
	armory   *process.BitcoindProcess
}

func New(profile *config.Profile) (*Launcher, error) {
	l := Launcher{
		profile: profile,
	}
	return &l, nil
}

func (l *Launcher) Run() error {
	// Create armory instance.
	c, err := process.NewArmoryProcess(l.profile)
	if err != nil {
		logrus.Fatalf("Error creating Armory instance: %s", err)
		return err
	}

	// Start armory.
	if err = c.Start(); err != nil {
		logrus.Fatalf("Error starting Armory: %s", err)
		return err
	}

	// Show output from armory, blocks until execution is over.
	c.PrintOutput()

	// Wait until armory has shut down.
	if err = c.WaitForShutdown(); err != nil {
		logrus.Errorf("Armory exited with error: %s", err)
		return nil
	}

	return nil
}
