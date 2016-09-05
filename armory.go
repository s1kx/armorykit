package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"sort"
	"strings"
	"sync"

	"github.com/Sirupsen/logrus"
)

type armorySettings struct {
	Bitcoind bitcoindSettings
	Profile  profileSettings
}

type ArmoryInstance struct {
	settings armorySettings

	cmd    *exec.Cmd
	stdout io.ReadCloser
	stderr io.ReadCloser
}

func NewArmoryInstance(settings armorySettings) (*ArmoryInstance, error) {
	// Build Command
	armoryFlags := armoryProfileFlags(&settings)
	armoryCmd := []string{"armory"}
	armoryCmd = append(armoryCmd, armoryFlags...)

	args := []string{"-xc", strings.Join(armoryCmd, " ")}
	args = append(args, armoryFlags...)
	fmt.Printf("Shell args: %v\n", args)
	cmd := exec.Command("/bin/sh", args...)

	// Read stdout and stderr from command
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	c := &ArmoryInstance{
		settings: settings,
		cmd:      cmd,
		stdout:   stdout,
		stderr:   stderr,
	}
	return c, nil
}

// Start launches the Armory instance.
func (c *ArmoryInstance) Start() error {
	return c.cmd.Start()
}

// Stop sends a TERM signal to the Armory instance.
func (c *ArmoryInstance) Stop() error {
	// Get Armory process
	process := c.cmd.Process
	if process == nil {
		return errors.New("armory is not running")
	}

	// TO-DO: Send TERM signal to initiate shutdown

	return nil
}

// PrintOutput prints the output of stdout/stderr from Armory to the logger.
func (c *ArmoryInstance) PrintOutput() {
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
func (c *ArmoryInstance) WaitForShutdown() error {
	return c.cmd.Wait()
}

func (c *ArmoryInstance) printPipeOutput(label string, pipe io.Reader) {
	rd := bufio.NewScanner(pipe)

	for rd.Scan() {
		line := rd.Text()

		logrus.Infof("Armory[%s]: %s", label, line)
	}
}

func armoryProfileFlags(s *armorySettings) []string {
	// Map of flags to be passed to armory
	flagMap := map[string]string{
		"--datadir":         s.Profile.ArmoryDataDir,
		"--satoshi-datadir": s.Bitcoind.DataDir,
	}

	// Create array of flags from map
	flags := armoryFlagsFromMap(flagMap)

	return flags
}

func armoryFlagsFromMap(flagMap map[string]string) []string {
	// Sort keys alphabetically to avoid random argument order
	keys := make([]string, 0, len(flagMap))
	for k, _ := range flagMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Create flags from map
	flags := make([]string, 0, len(flagMap))

	for _, k := range keys {
		arg := flagMap[k]

		flag := fmt.Sprintf("%s=%s", k, arg)
		flags = append(flags, flag)
	}

	return flags
}
