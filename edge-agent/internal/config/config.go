package config

import (
	"fmt"
	"os"
	"gopkg.in/yaml.v3"
)

type Config struct {
	AgentID        string           `yaml:"agent_id"`
	Hostname       string           `yaml:"hostname"`
	Version        string           `yaml:"version"`
	BootstrapToken string           `yaml:"bootstrap_token"`
	CloudEndpoint  string           `yaml:"cloud_endpoint"`
	DataCenterID   string           `yaml:"datacenter_id"`
	Discovery      DiscoveryConfig  `yaml:"discovery"`
	Collection     CollectionConfig `yaml:"collection"`
}

type DiscoveryConfig struct {
	RedfishRanges []string `yaml:"redfish_ranges"`
	SNMPRanges    []string `yaml:"snmp_ranges"`
	IntervalSec   int      `yaml:"interval_sec"`
	TimeoutSec    int      `yaml:"timeout_sec"`
}

type CollectionConfig struct {
	IntervalSec int `yaml:"interval_sec"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	if cfg.Version == "" {
		cfg.Version = "0.1.0"
	}
	if cfg.Discovery.IntervalSec == 0 {
		cfg.Discovery.IntervalSec = 3600
	}
	if cfg.Discovery.TimeoutSec == 0 {
		cfg.Discovery.TimeoutSec = 30
	}
	if cfg.Collection.IntervalSec == 0 {
		cfg.Collection.IntervalSec = 30
	}
	return &cfg, nil
}

func (c *Config) Save(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	return nil
}

func (c *Config) Validate() error {
	if c.CloudEndpoint == "" {
		return fmt.Errorf("cloud_endpoint is required")
	}
	if c.BootstrapToken == "" && c.AgentID == "" {
		return fmt.Errorf("either bootstrap_token or agent_id is required")
	}
	return nil
}
