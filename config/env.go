package config

import (
	"io/ioutil"
	"nav_sync/utils"

	"gopkg.in/yaml.v2"
)

// Model classes
type Configuration struct {
	Vendor        UrlConfig  `yaml:"vendor"`
	Invoice       UrlConfig  `yaml:"invoice"`
	LedgerEntries UrlConfig  `yaml:"ledger_entries"`
	Auth          AuthConfig `yaml:"auth"`
}

type UrlConfig struct {
	Fetch      FetchConfig `yaml:"fetch"`
	Sync       SyncConfig  `yaml:"sync"`
	Post       SyncConfig  `yaml:"post"`
	Save       SyncConfig  `yaml:"save"`
	FakeInsert bool        `yaml:"fake_insert"`
	Prefix     string      `yaml:"prefix"`
	EmptyLogs  bool        `yaml:"empty_logs"`
}

type AuthConfig struct {
	Ntlm NtlmConfig `yaml:"ntlm"`
}

type FetchConfig struct {
	URL    string `yaml:"url"`
	APIKey string `yaml:"apikey"`
}

type SyncConfig struct {
	URL string `yaml:"url"`
}

type NtlmConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// Load configuration yaml file
func LoadYamlFile() Configuration {
	//Reading config.yml file
	data, err := ioutil.ReadFile("config.yml")
	if err != nil {
		utils.Fatal("Error loading config.yml file")
	}

	//Unmarshal config data
	err = yaml.Unmarshal(data, &Config)
	if err != nil {
		utils.Fatal("Error reading config.yml file")
	}
	return Config
}

var Config Configuration
