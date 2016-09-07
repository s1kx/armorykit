package process

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"
	"syscall"

	"github.com/Sirupsen/logrus"

	"github.com/s1kx/armorykit/config"
)

const armoryBinary = "armory"

type ArmoryProcess struct {
	profile *config.Profile

	cmd    *exec.Cmd
	stdout io.ReadCloser
	stderr io.ReadCloser
}

func NewArmoryProcess(profile *config.Profile) (*ArmoryProcess, error) {
	// Build Command
	flags := armoryFlags(profile)
	shellCmd := fmt.Sprintf("%s %s", armoryBinary, strings.Join(flags, " "))
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

	c := &ArmoryProcess{
		profile: profile,

		cmd:    cmd,
		stdout: stdout,
		stderr: stderr,
	}
	return c, nil
}

// Start launches the Armory instance.
func (c *ArmoryProcess) Start() error {
	return c.cmd.Start()
}

// Stop sends a TERM signal to the Armory instance.
func (c *ArmoryProcess) Stop() error {
	// Get Armory process
	process := c.cmd.Process
	if process == nil {
		return errors.New("armory is not running")
	}

	// Send SIGINT signal to initiate graceful shutdown
	process.Signal(syscall.SIGINT)

	return nil
}

// PrintOutput prints the output of stdout/stderr from Armory to the logger.
func (c *ArmoryProcess) PrintOutput() {
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
func (c *ArmoryProcess) WaitForShutdown() error {
	return c.cmd.Wait()
}

func (c *ArmoryProcess) printPipeOutput(label string, pipe io.Reader) {
	r := bufio.NewScanner(pipe)
	for r.Scan() {
		line := r.Text()

		logrus.Debugf("Armory[%s]: %s", label, line)
	}
}

func armoryFlags(p *config.Profile) []string {
	flags := FlagList{
		{"--datadir", p.ProfileSettings.ArmoryDataDir},
		{"--satoshi-datadir", p.BitcoindSettings.DataDir},
	}

	return flags.StringList()
}
