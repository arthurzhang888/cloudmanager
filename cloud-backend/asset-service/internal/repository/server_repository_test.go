package repository

import (
	"testing"

	"github.com/cloudmanager/cloud-backend/asset-service/internal/models"
	"github.com/google/uuid"
)

func TestServerRepository_ScanServers(t *testing.T) {
	// This is a basic test structure - in real tests, you'd use a test database
	repo := &ServerRepository{}
	_ = repo
}

func TestServerFilters(t *testing.T) {
	filters := map[string]interface{}{
		"agent_id": uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
		"status":   "online",
	}

	if _, ok := filters["agent_id"].(uuid.UUID); !ok {
		t.Error("agent_id should be a UUID")
	}

	if _, ok := filters["status"].(string); !ok {
		t.Error("status should be a string")
	}
}

func TestServerModel(t *testing.T) {
	server := &models.Server{
		ID:           uuid.New(),
		AgentID:      uuid.New(),
		Manufacturer: "Supermicro",
		Model:        "SuperServer 2029",
		SerialNumber: "SM2029A12345",
		CPUCount:     2,
		CPUModel:     "Intel Xeon Gold 6248R",
		MemoryGB:     256,
		Status:       "online",
		PowerState:   "On",
		Health:       "OK",
	}

	if server.Manufacturer != "Supermicro" {
		t.Errorf("expected manufacturer Supermicro, got %s", server.Manufacturer)
	}

	if server.CPUCount != 2 {
		t.Errorf("expected CPU count 2, got %d", server.CPUCount)
	}
}

func TestDataCenterModel(t *testing.T) {
	dc := &models.DataCenter{
		ID:          uuid.New(),
		Name:        "Beijing DC",
		Location:    "Beijing, China",
		NetworkType: "vpn",
	}

	if dc.Name != "Beijing DC" {
		t.Errorf("expected name Beijing DC, got %s", dc.Name)
	}

	if dc.NetworkType != "vpn" {
		t.Errorf("expected network type vpn, got %s", dc.NetworkType)
	}
}

func TestEdgeAgentModel(t *testing.T) {
	agent := &models.EdgeAgent{
		ID:       uuid.New(),
		DCID:     uuid.New(),
		Hostname: "edge-dc1-01",
		Version:  "0.1.0",
		Status:   "online",
	}

	if agent.Hostname != "edge-dc1-01" {
		t.Errorf("expected hostname edge-dc1-01, got %s", agent.Hostname)
	}

	if agent.Status != "online" {
		t.Errorf("expected status online, got %s", agent.Status)
	}
}

func TestCoolingDeviceModel(t *testing.T) {
	device := &models.CoolingDevice{
		ID:            uuid.New(),
		AgentID:       uuid.New(),
		DeviceType:    "cdu",
		Name:          "CDU-Rack-A01",
		ModbusAddress: "tcp://192.168.10.10:502",
		SlaveID:       1,
		Location:      "Rack A01-A10",
		Status:        "online",
	}

	if device.DeviceType != "cdu" {
		t.Errorf("expected device type cdu, got %s", device.DeviceType)
	}

	if device.SlaveID != 1 {
		t.Errorf("expected slave ID 1, got %d", device.SlaveID)
	}
}
