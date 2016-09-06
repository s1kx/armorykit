package config

import "errors"

var (
	ErrInvalidProfile = errors.New("config: invalid profile name")
)

type Config struct {
	BitcoindSettings BitcoindSettings           `toml:"bitcoind"`
	ProfileSettings  map[string]ProfileSettings `toml:"profiles"`
}

type BitcoindSettings struct {
	Host       string
	Port       int
	RPCPort    int
	DataDir    string
	ConfigFile string `toml:"config"`
	Flags      []string
}

type ProfileSettings struct {
	ArmoryDataDir string   `toml:"armory_datadir"`
	BitcoindFlags []string `toml:"bitcoind_flags"`
}

// Profile is a composite of the global bitcoind settings and a specific profile.
type Profile struct {
	BitcoindSettings
	ProfileSettings
}

func (conf *Config) GetProfile(profileKey string) (*Profile, error) {
	profileSettings, ok := conf.ProfileSettings[profileKey]
	if !ok {
		return nil, ErrInvalidProfile
	}

	profile := Profile{
		BitcoindSettings: conf.BitcoindSettings,
		ProfileSettings:  profileSettings,
	}
	return &profile, nil
}
