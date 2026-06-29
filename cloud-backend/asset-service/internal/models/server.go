package models

import (
	"time"

	"github.com/google/uuid"
)

// Server represents a physical server in the data center
type Server struct {
	ID              uuid.UUID `json:"id"`
	AgentID         uuid.UUID `json:"agent_id"`
	RedfishEndpoint string    `json:"redfish_endpoint"`
	Manufacturer    string    `json:"manufacturer"`
	Model           string    `json:"model"`
	SerialNumber    string    `json:"serial_number"`
	SKU             string    `json:"sku"`
	CPUCount        int       `json:"cpu_count"`
	CPUModel        string    `json:"cpu_model"`
	MemoryGB        int       `json:"memory_gb"`
	DiskCount       int       `json:"disk_count"`
	TotalDiskGB     int64     `json:"total_disk_gb"`
	Status          string    `json:"status"` // 'online', 'offline', 'error', 'maintenance', 'unknown'
	PowerState      string    `json:"power_state"`
	Health          string    `json:"health"`
	RawData         []byte    `json:"raw_data"`
	DiscoveredAt    time.Time `json:"discovered_at"`
	LastUpdated     time.Time `json:"last_updated"`
	CreatedAt       time.Time `json:"created_at"`
}

// EdgeAgent represents an edge agent that collects data from devices
type EdgeAgent struct {
	ID         uuid.UUID              `json:"id"`
	DCID       uuid.UUID              `json:"dc_id"`
	Hostname   string                 `json:"hostname"`
	Version    string                 `json:"version"`
	LastSeen   *time.Time             `json:"last_seen"`
	CertSN     string                 `json:"cert_sn"`
	Status     string                 `json:"status"` // 'online', 'offline', 'error'
	Config     map[string]interface{} `json:"config"`
	CreatedAt  time.Time              `json:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at"`
}

// DataCenter represents a data center location
type DataCenter struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Location     string    `json:"location"`
	NetworkType  string    `json:"network_type"` // 'vpn', 'public', 'airgap'
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CoolingDevice represents a liquid cooling device (CDU or cooling tower)
type CoolingDevice struct {
	ID            uuid.UUID              `json:"id"`
	AgentID       uuid.UUID              `json:"agent_id"`
	DeviceType    string                 `json:"device_type"` // 'cdu', 'cooling_tower'
	Name          string                 `json:"name"`
	ModbusAddress string                 `json:"modbus_address"`
	SlaveID       int                    `json:"slave_id"`
	Location      string                 `json:"location"`
	RegisterMap   map[string]interface{} `json:"register_map"`
	Status        string                 `json:"status"` // 'online', 'offline', 'error', 'unknown'
	DiscoveredAt  time.Time              `json:"discovered_at"`
	LastUpdated   time.Time              `json:"last_updated"`
}

// NetworkDevice represents a network device discovered via SNMP
type NetworkDevice struct {
	ID           uuid.UUID `json:"id"`
	AgentID      uuid.UUID `json:"agent_id"`
	IPAddress    string    `json:"ip_address"`
	SNMPCommunity string   `json:"snmp_community"`
	Manufacturer string    `json:"manufacturer"`
	Model        string    `json:"model"`
	SerialNumber string    `json:"serial_number"`
	DeviceType   string    `json:"device_type"` // 'switch', 'router', 'firewall'
	PortCount    int       `json:"port_count"`
	Status       string    `json:"status"`
	DiscoveredAt time.Time `json:"discovered_at"`
	LastUpdated  time.Time `json:"last_updated"`
}
