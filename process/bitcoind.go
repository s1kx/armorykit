package process

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"

	"github.com/Sirupsen/logrus"

	"github.com/s1kx/armorykit/config"
)

const BitcoindBinary = "bitcoind"

type BitcoindProcess struct {
	profile *config.Profile

	cmd    *exec.Cmd
	stdout io.ReadCloser
	stderr io.ReadCloser
}

func NewBitcoindProcess(profile *config.Profile) (*BitcoindProcess, error) {
	// Build Command
	flags := bitcoindFlags(profile)
	shellCmd := fmt.Sprintf("%s %s", BitcoindBinary, strings.Join(flags, " "))
	cmd := exec.Command("/bin/sh", "-xc", shellCmd)

	// Read stdout and stderr from command
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	c := &BitcoindProcess{
		profile: profile,

		cmd:    cmd,
		stdout: stdout,
		stderr: stderr,
	}
	return c, nil
}

// Start launches the Armory instance.
func (c *BitcoindProcess) Start() error {
	return c.cmd.Start()
}

// Stop sends a TERM signal to the Armory instance.
func (c *BitcoindProcess) Stop() error {
	// Get Armory process
	process := c.cmd.Process
	if process == nil {
		return errors.New("armory is not running")
	}

	// TO-DO: Send TERM signal to initiate shutdown

	return nil
}

// PrintOutput prints the output of stdout/stderr from Armory to the logger.
func (c *BitcoindProcess) PrintOutput() {
	wg := &sync.WaitGroup{}

	pipes := []struct {
		label string
		pipe  io.ReadCloser
	}{
		{"stdout", c.stdout},
		{"stderr", c.stderr},
	}
	for _, p := range pipes {
		wg.Add(1)
		go func() {
			c.printPipeOutput(p.label, p.pipe)
			wg.Done()
		}()
	}

	wg.Wait()
}

// WaitForShutdown waits until the process finishes.
// This function should be called after PrintOutput has completed.
func (c *BitcoindProcess) WaitForShutdown() error {
	return c.cmd.Wait()
}

func (c *BitcoindProcess) printPipeOutput(label string, pipe io.Reader) {
	r := bufio.NewScanner(pipe)
	for r.Scan() {
		line := r.Text()

		logrus.Debugf("bitcoind[%s]: %s", label, line)
	}
}

func bitcoindFlags(p *config.Profile) []string {
	flagList := FlagList{
		{"--conf", p.BitcoindSettings.ConfigFile},
		{"--datadir", p.ProfileSettings.ArmoryDataDir},
	}

	flags := flagList.StringList()

	// Add flags from bitcoind config
	flags = append(flags, p.BitcoindSettings.Flags...)

	// Add flags from profile config
	flags = append(flags, p.ProfileSettings.BitcoindFlags...)

	return flags
}