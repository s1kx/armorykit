package launcher

import (
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/s1kx/armorykit/config"
)

type Launcher struct {
	profile *config.Profile

	bitcoind *bitcoindLauncher
	armory   *armoryLauncher
}

func New(profile *config.Profile) (*Launcher, error) {
	bl, err := newBitcoindLauncher(profile)
	if err != nil {
		return nil, err
	}

	al, err := newArmoryLauncher(profile)
	if err != nil {
		return nil, err
	}

	l := Launcher{
		profile: profile,

		bitcoind: bl,
		armory:   al,
	}
	return &l, nil
}

func (l *Launcher) Run() error {
	// Check if bitcoin is already running
	if err := l.bitcoind.IsRunnable(); err != nil {
		return err
	}

	// Start bitcoind
	if err := l.bitcoind.Start(); err != nil {
		return err
	}
	defer l.bitcoind.Stop()

	// Start armory
	if err := l.armory.Start(); err != nil {
		return err
	}
	defer l.armory.Stop()

	// Run processes, stop both if one stops
	wg := &sync.WaitGroup{}
	quit := make(chan struct{}, 2)

	// Redirect bitcoind output and wait for shutdown
	wg.Add(1)
	go func() {
		if err := l.bitcoind.ProcessAndWait(); err != nil {
			logrus.Errorf("launcher: %s", err)
		}
		logrus.Debug("bitcoind has exited")

		quit <- struct{}{}
		wg.Done()
	}()

	// Redirect armory output and wait for shutdown
	wg.Add(1)
	go func() {
		if err := l.armory.ProcessAndWait(); err != nil {
			logrus.Errorf("launcher: %s", err)
		}
		logrus.Debug("armory has exited")

		quit <- struct{}{}
		wg.Done()
	}()

	// Wait for first process to exit
	<-quit

	// Stop both processes
	l.bitcoind.Stop()
	l.armory.Stop()

	// Wait for processes to finish
	wg.Wait()

	return nil
}
