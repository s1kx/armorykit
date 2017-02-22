package launcher

import (
	"errors"
	"fmt"
	"net"

	"github.com/Sirupsen/logrus"

	"github.com/s1kx/armorykit/config"
	"github.com/s1kx/armorykit/process"
)

type bitcoindLauncher struct {
	profile *config.Profile

	proc *process.BitcoindProcess
}

func newBitcoindLauncher(profile *config.Profile) (*bitcoindLauncher, error) {
	// Create armory process/command
	proc, err := process.NewBitcoindProcess(profile)
	if err != nil {
		return nil, fmt.Errorf("bitcoind: %s", err)
	}

	launcher := bitcoindLauncher{
		profile: profile,
		proc:    proc,
	}
	return &launcher, nil
}

func (l *bitcoindLauncher) IsRunnable() error {
	// Check if port is already in use
	host := l.profile.BitcoindSettings.Host
	port := l.profile.BitcoindSettings.Port
	addr := fmt.Sprintf("%s:%d", host, port)

	logrus.WithField("addr", addr).Debug("launcher: checking if bitcoin port is in use")
	if err := CheckOpenPort(addr); err != nil {
		return fmt.Errorf("bitcoind: %s", err)
	}

	return nil
}

func (l *bitcoindLauncher) Start() error {
	// Start process
	if err := l.proc.Start(); err != nil {
		return fmt.Errorf("start bitcoind: %s", err)
	}

	pid := l.proc.Process().Pid

	logrus.WithField("pid", pid).Info("bitcoind: started")

	return nil
}

func (l *bitcoindLauncher) Stop() error {
	return l.proc.Stop()
}

func (l *bitcoindLauncher) ProcessAndWait() error {
	// Show output from bitcoind, blocks until execution is over
	l.proc.PrintOutput()

	// Wait until armory has shut down
	if err := l.proc.WaitForShutdown(); err != nil {
		logrus.Warnf("bitcoind: exited with error: %s", err)
		return nil
	}

	logrus.Info("bitcoind: exited")

	return nil
}

// CheckOpenPort tests if a port is not already in use by another application.
func CheckOpenPort(addr string) error {
	// Attempt to listen on port to determine if it is already in use
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	err = ln.Close()
	if err != nil {
		return errors.New("could not stop listening when checking open port")
	}

	return nil
}
