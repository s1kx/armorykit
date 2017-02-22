package launcher

import (
	"fmt"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/s1kx/armorykit/config"
	"github.com/s1kx/armorykit/process"
)

type armoryLauncher struct {
	sync.Mutex

	profile *config.Profile

	proc *process.ArmoryProcess
}

func newArmoryLauncher(profile *config.Profile) (*armoryLauncher, error) {
	// Create armory process/command
	proc, err := process.NewArmoryProcess(profile)
	if err != nil {
		return nil, fmt.Errorf("armory: %s", err)
	}

	launcher := armoryLauncher{
		profile: profile,
		proc:    proc,
	}
	return &launcher, nil
}

func (l *armoryLauncher) Start() error {
	// Start armory
	if err := l.proc.Start(); err != nil {
		return fmt.Errorf("start armory: %s", err)
	}

	pid := l.proc.Process().Pid

	logrus.WithField("pid", pid).Info("armory: started")

	return nil
}

func (l *armoryLauncher) Stop() error {
	return l.proc.Stop()
}

func (l *armoryLauncher) ProcessAndWait() error {
	// Show output from armory, blocks until execution is over
	l.proc.PrintOutput()

	// Wait until armory has shut down
	if err := l.proc.WaitForShutdown(); err != nil {
		logrus.Warnf("armory exited with error: %s", err)
		return nil
	}

	logrus.Info("armory: exited")

	return nil
}
