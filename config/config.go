package config

import "errors"

type Config struct {
	BitcoindSettings BitcoindSettings           `toml:"bitcoind"`
	ProfileSettings  map[string]ProfileSettings `toml:"profiles"`
}

type BitcoindSettings struct {
	Host    string
	Port    int
	RPCPort int
	DataDir string
	Config  string
	Flags   []string
}

type ProfileSettings struct {
	ArmoryDataDir string   `toml:"armory_datadir"`
	BitcoindFlags []string `toml:"bitcoind_flags"`
}

type Profile struct {
	BitcoindSettings
	ProfileSettings
}

func (c Config) GetProfile(profileKey string) (*Profile, error) {
	profileSettings, ok := c.ProfileSettings[profileKey]
	if !ok {
		return nil, errors.New("invalid profile name, check your config")
	}

	profile := Profile{
		BitcoindSettings: c.BitcoindSettings,
		ProfileSettings:  profileSettings,
	}
	return &profile, nil
}
