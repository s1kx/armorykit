package main

import (
	"os/user"
	"strings"

	"github.com/BurntSushi/toml"
)

type config struct {
	Bitcoind bitcoindSettings
	Profiles map[string]profileSettings
}

type bitcoindSettings struct {
	Host    string
	Port    int
	RPCPort int
	DataDir string
	Config  string
	Flags   []string
}

type profileSettings struct {
	ArmoryDataDir string   `toml:"armory_datadir"`
	BitcoindFlags []string `toml:"bitcoind_flags"`
}

func loadConfig(path string, conf *config) error {
	path = expandUserHomePath(path)
	_, err := toml.DecodeFile(path, &conf)
	return err
}

func expandUserHomePath(path string) string {
	usr, _ := user.Current()
	dir := usr.HomeDir
	if !strings.HasSuffix(dir, "/") {
		dir += "/"
	}
	if path[:2] == "~/" {
		path = strings.Replace(path, "~/", dir, 1)
	}
	return path
}
