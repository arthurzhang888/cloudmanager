package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_ValidConfig(t *testing.T) {
	content := `
agent_id: "agent-123"
hostname: "edge-dc1-01"
cloud_endpoint: "cloud.example.com:443"
bootstrap_token: "test-token"
discovery:
  redfish_ranges:
    - "192.168.1.0/24"
  interval_sec: 3600
`
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	os.WriteFile(configPath, []byte(content), 0644)

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if cfg.AgentID != "agent-123" {
		t.Errorf("expected AgentID=agent-123, got %s", cfg.AgentID)
	}
	if cfg.CloudEndpoint != "cloud.example.com:443" {
		t.Errorf("expected CloudEndpoint=cloud.example.com:443, got %s", cfg.CloudEndpoint)
	}
}

func TestLoad_Defaults(t *testing.T) {
	content := `
cloud_endpoint: "cloud.example.com:443"
bootstrap_token: "test-token"
`
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	os.WriteFile(configPath, []byte(content), 0644)

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if cfg.Version != "0.1.0" {
		t.Errorf("expected default Version=0.1.0, got %s", cfg.Version)
	}
	if cfg.Discovery.IntervalSec != 3600 {
		t.Errorf("expected default IntervalSec=3600, got %d", cfg.Discovery.IntervalSec)
	}
	if cfg.Collection.IntervalSec != 30 {
		t.Errorf("expected default CollectionInterval=30, got %d", cfg.Collection.IntervalSec)
	}
}

func TestValidate_MissingCloudEndpoint(t *testing.T) {
	cfg := &Config{BootstrapToken: "test-token"}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for missing cloud_endpoint")
	}
}

func TestValidate_MissingAuth(t *testing.T) {
	cfg := &Config{CloudEndpoint: "cloud.example.com:443"}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for missing auth")
	}
}

func TestSave(t *testing.T) {
	cfg := &Config{
		AgentID:        "agent-456",
		Hostname:       "test-host",
		CloudEndpoint:  "cloud.example.com:443",
		BootstrapToken: "test-token",
	}

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	if err := cfg.Save(configPath); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load after save failed: %v", err)
	}

	if loaded.AgentID != cfg.AgentID {
		t.Errorf("expected AgentID=%s, got %s", cfg.AgentID, loaded.AgentID)
	}
}
