# CloudManager MVP Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build CloudManager MVP with edge asset discovery, liquid cooling monitoring via Modbus, and web dashboards for 1000+ server data centers.

**Architecture:** Hybrid deployment with Go-based edge agents collecting Redfish/Modbus data, reporting to cloud backend via gRPC/mTLS. Cloud uses microservices (API Gateway, Asset Service, Telemetry Service) with PostgreSQL/TimescaleDB storage and React frontend.

**Tech Stack:** Go, gRPC, PostgreSQL, TimescaleDB, Redis, NATS, React/TypeScript, Docker

---

## File Structure

```
cloudmanager/
├── edge-agent/                 # Edge data collection agent
│   ├── cmd/agent/
│   │   └── main.go
│   ├── internal/
│   │   ├── config/
│   │   │   └── config.go
│   │   ├── discovery/
│   │   │   ├── redfish.go
│   │   │   ├── snmp.go
│   │   │   └── scanner.go
│   │   ├── collector/
│   │   │   ├── modbus.go
│   │   │   ├── redfish_metrics.go
│   │   │   └── scheduler.go
│   │   ├── uploader/
│   │   │   ├── grpc_client.go
│   │   │   └── buffer.go
│   │   └── cert/
│   │       └── manager.go
│   ├── proto/
│   │   └── edge.proto
│   └── go.mod
├── cloud-backend/              # Cloud control plane
│   ├── api-gateway/
│   ├── asset-service/
│   ├── telemetry-service/
│   ├── discovery-service/
│   ├── shared/
│   └── go.mod
├── web-console/                # React frontend
├── migrations/                 # Database migrations
└── docker-compose.yml
```

---

## Phase 1: Project Setup & Infrastructure

### Task 1: Initialize Project Structure

**Files:**
- Create: `go.work`
- Create: `edge-agent/go.mod`
- Create: `cloud-backend/go.mod`
- Create: `.gitignore`

- [ ] **Step 1: Create workspace and module files**

```go
// go.work
go 1.22

use (
    ./edge-agent
    ./cloud-backend
)
```

```go
// edge-agent/go.mod
module github.com/cloudmanager/edge-agent

go 1.22

require (
    github.com/go-resty/resty/v2 v2.12.0
    github.com/goburrow/modbus v0.1.0
    github.com/google/uuid v1.6.0
    github.com/nats-io/nats.go v1.35.0
    google.golang.org/grpc v1.64.0
    google.golang.org/protobuf v1.34.1
    gopkg.in/yaml.v3 v3.0.1
)
```

```go
// cloud-backend/go.mod
module github.com/cloudmanager/cloud-backend

go 1.22

require (
    github.com/gin-gonic/gin v1.9.1
    github.com/go-redis/redis/v8 v8.11.5
    github.com/google/uuid v1.6.0
    github.com/jackc/pgx/v5 v5.6.0
    github.com/nats-io/nats.go v1.35.0
    google.golang.org/grpc v1.64.0
    google.golang.org/protobuf v1.34.1
)
```

```gitignore
# .gitignore
*.exe
*.dll
*.so
*.dylib
*.test
*.out
vendor/
.env
.env.local
.idea/
.vscode/
*.swp
*.swo
*~
.DS_Store
node_modules/
dist/
build/
*.log
```

- [ ] **Step 2: Commit project structure**

```bash
git add go.work edge-agent/go.mod cloud-backend/go.mod .gitignore
git commit -m "chore: initialize project structure with Go modules"
```

---

### Task 2: Setup Docker Infrastructure

**Files:**
- Create: `docker-compose.yml`
- Create: `.env.example`

- [ ] **Step 1: Create Docker Compose file**

```yaml
# docker-compose.yml
version: '3.8'

services:
  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: ${DB_USER:-cloudmanager}
      POSTGRES_PASSWORD: ${DB_PASSWORD:-changeme}
      POSTGRES_DB: ${DB_NAME:-cloudmanager}
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U cloudmanager"]
      interval: 5s
      timeout: 5s
      retries: 5

  timescaledb:
    image: timescale/timescaledb:latest-pg16
    environment:
      POSTGRES_USER: ${TSDB_USER:-tsdb}
      POSTGRES_PASSWORD: ${TSDB_PASSWORD:-changeme}
      POSTGRES_DB: ${TSDB_NAME:-telemetry}
    ports:
      - "5433:5432"
    volumes:
      - timescale_data:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data

  nats:
    image: nats:2-alpine
    ports:
      - "4222:4222"
      - "8222:8222"
    command: "-js -m 8222"

volumes:
  postgres_data:
  timescale_data:
  redis_data:
```

```bash
# .env.example
DB_USER=cloudmanager
DB_PASSWORD=changeme
DB_NAME=cloudmanager
DB_HOST=localhost
DB_PORT=5432

TSDB_USER=tsdb
TSDB_PASSWORD=changeme
TSDB_NAME=telemetry
TSDB_HOST=localhost
TSDB_PORT=5433

REDIS_HOST=localhost
REDIS_PORT=6379

NATS_URL=nats://localhost:4222

JWT_SECRET=change-this-in-production
```

- [ ] **Step 2: Commit Docker infrastructure**

```bash
git add docker-compose.yml .env.example
git commit -m "chore: add Docker Compose for local development infrastructure"
```

---

## Phase 2: Database Schema

### Task 3: Create Database Migrations

**Files:**
- Create: `migrations/001_initial.sql`

- [ ] **Step 1: Write initial migration**

```sql
-- migrations/001_initial.sql

CREATE EXTENSION IF NOT EXISTS timescaledb;

-- Data centers
CREATE TABLE data_centers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    location VARCHAR(500),
    network_type VARCHAR(50) NOT NULL DEFAULT 'vpn',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Edge agents
CREATE TABLE edge_agents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    dc_id UUID NOT NULL REFERENCES data_centers(id) ON DELETE CASCADE,
    hostname VARCHAR(255) NOT NULL,
    version VARCHAR(50) NOT NULL DEFAULT '0.1.0',
    last_seen TIMESTAMPTZ,
    cert_sn VARCHAR(255),
    status VARCHAR(50) DEFAULT 'offline',
    config JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_edge_agents_dc ON edge_agents(dc_id);
CREATE INDEX idx_edge_agents_status ON edge_agents(status);

-- Servers
CREATE TABLE servers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    agent_id UUID NOT NULL REFERENCES edge_agents(id) ON DELETE CASCADE,
    redfish_endpoint VARCHAR(500),
    manufacturer VARCHAR(255),
    model VARCHAR(255),
    serial_number VARCHAR(255) UNIQUE,
    cpu_count INTEGER DEFAULT 0,
    memory_gb INTEGER DEFAULT 0,
    status VARCHAR(50) DEFAULT 'unknown',
    power_state VARCHAR(50),
    health VARCHAR(50),
    raw_data JSONB DEFAULT '{}',
    discovered_at TIMESTAMPTZ DEFAULT NOW(),
    last_updated TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_servers_agent ON servers(agent_id);
CREATE INDEX idx_servers_status ON servers(status);

-- Cooling devices
CREATE TABLE cooling_devices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    agent_id UUID NOT NULL REFERENCES edge_agents(id) ON DELETE CASCADE,
    device_type VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    modbus_address VARCHAR(255) NOT NULL,
    slave_id INTEGER DEFAULT 1,
    location VARCHAR(500),
    register_map JSONB DEFAULT '{}',
    status VARCHAR(50) DEFAULT 'unknown',
    discovered_at TIMESTAMPTZ DEFAULT NOW(),
    last_updated TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_cooling_devices_agent ON cooling_devices(agent_id);

-- Sensor readings (TimescaleDB hypertable)
CREATE TABLE sensor_readings (
    time TIMESTAMPTZ NOT NULL,
    device_id UUID NOT NULL,
    device_type VARCHAR(50) NOT NULL,
    agent_id UUID NOT NULL,
    metric_name VARCHAR(255) NOT NULL,
    value DOUBLE PRECISION NOT NULL,
    unit VARCHAR(50),
    labels JSONB DEFAULT '{}'
);

SELECT create_hypertable('sensor_readings', 'time', chunk_time_interval => INTERVAL '1 day');

CREATE INDEX idx_sensor_readings_device ON sensor_readings(device_id, metric_name, time DESC);

-- Alerts
CREATE TABLE alerts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    severity VARCHAR(50) NOT NULL,
    category VARCHAR(50) NOT NULL,
    source_type VARCHAR(50) NOT NULL,
    source_id UUID NOT NULL,
    title VARCHAR(500) NOT NULL,
    description TEXT,
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_alerts_status ON alerts(status);
CREATE INDEX idx_alerts_created ON alerts(created_at DESC);
```

- [ ] **Step 2: Commit migration**

```bash
git add migrations/001_initial.sql
git commit -m "feat(db): add initial database schema"
```

---

## Phase 3: Protocol Buffers

### Task 4: Define gRPC Proto Files

**Files:**
- Create: `edge-agent/proto/edge.proto`

- [ ] **Step 1: Create edge.proto**

```protobuf
syntax = "proto3";

package edge;
option go_package = "github.com/cloudmanager/edge-agent/proto;edgepb";

import "google/protobuf/timestamp.proto";

service EdgeAgentService {
  rpc Bootstrap(BootstrapRequest) returns (BootstrapResponse);
  rpc Heartbeat(stream HeartbeatRequest) returns (stream HeartbeatResponse);
  rpc ReportDiscovery(ReportDiscoveryRequest) returns (ReportDiscoveryResponse);
  rpc ReportTelemetry(ReportTelemetryRequest) returns (ReportTelemetryResponse);
}

message BootstrapRequest {
  string bootstrap_token = 1;
  string hostname = 2;
  string version = 3;
  string public_key = 4;
}

message BootstrapResponse {
  string agent_id = 1;
  bytes certificate = 2;
  bytes ca_certificate = 3;
  AgentConfig config = 4;
}

message AgentConfig {
  string dc_id = 1;
  int64 heartbeat_interval_sec = 2;
  DiscoveryConfig discovery = 3;
  CollectionConfig collection = 4;
}

message DiscoveryConfig {
  repeated string redfish_ip_ranges = 1;
  repeated string snmp_ip_ranges = 2;
  int64 scan_interval_sec = 3;
}

message CollectionConfig {
  repeated CoolingDeviceConfig cooling_devices = 1;
  int64 collection_interval_sec = 2;
}

message CoolingDeviceConfig {
  string device_id = 1;
  string name = 2;
  string device_type = 3;
  string modbus_address = 4;
  int32 slave_id = 5;
}

message HeartbeatRequest {
  string agent_id = 1;
  google.protobuf.Timestamp timestamp = 2;
  string status = 3;
}

message HeartbeatResponse {
  google.protobuf.Timestamp timestamp = 1;
  AgentConfig config_update = 2;
}

message ReportDiscoveryRequest {
  string agent_id = 1;
  google.protobuf.Timestamp timestamp = 2;
  repeated ServerInfo servers = 3;
  repeated CoolingDeviceInfo cooling_devices = 4;
}

message ReportDiscoveryResponse {
  bool success = 1;
}

message ServerInfo {
  string redfish_endpoint = 1;
  string manufacturer = 2;
  string model = 3;
  string serial_number = 4;
  int32 cpu_count = 5;
  int64 memory_gb = 6;
  string power_state = 7;
  string health = 8;
}

message CoolingDeviceInfo {
  string name = 1;
  string device_type = 2;
  string modbus_address = 3;
  int32 slave_id = 4;
}

message ReportTelemetryRequest {
  string agent_id = 1;
  google.protobuf.Timestamp timestamp = 2;
  repeated SensorReading readings = 3;
}

message ReportTelemetryResponse {
  bool success = 1;
  int32 accepted_count = 2;
}

message SensorReading {
  string device_id = 1;
  string device_type = 2;
  google.protobuf.Timestamp timestamp = 3;
  string metric_name = 4;
  double value = 5;
  string unit = 6;
}
```

- [ ] **Step 2: Generate Go code from proto**

```bash
cd /Users/arthurzhang/dev/llm/cloudManager
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
mkdir -p cloud-backend/shared/proto
cp edge-agent/proto/edge.proto cloud-backend/shared/proto/

cd edge-agent
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       proto/edge.proto
```

- [ ] **Step 3: Commit proto files**

```bash
git add edge-agent/proto/ cloud-backend/shared/proto/
git commit -m "feat(proto): define gRPC protocol"
```

---

## Phase 4: Edge Agent Implementation

### Task 5: Edge Agent Configuration Module

**Files:**
- Create: `edge-agent/internal/config/config.go`
- Create: `edge-agent/internal/config/config_test.go`

- [ ] **Step 1: Write config module with tests**

```go
// edge-agent/internal/config/config.go
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

func (c *Config) Validate() error {
	if c.CloudEndpoint == "" {
		return fmt.Errorf("cloud_endpoint is required")
	}
	if c.BootstrapToken == "" && c.AgentID == "" {
		return fmt.Errorf("either bootstrap_token or agent_id is required")
	}
	return nil
}
```

```go
// edge-agent/internal/config/config_test.go
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

func TestValidate_MissingCloudEndpoint(t *testing.T) {
	cfg := &Config{BootstrapToken: "test-token"}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for missing cloud_endpoint")
	}
}
```

- [ ] **Step 2: Run tests**

```bash
cd /Users/arthurzhang/dev/llm/cloudManager/edge-agent
go test ./internal/config/... -v
```

Expected: All tests pass

- [ ] **Step 3: Commit**

```bash
git add edge-agent/internal/config/
git commit -m "feat(edge): add configuration module with validation"
```

---

### Task 6: Modbus Collector Module

**Files:**
- Create: `edge-agent/internal/collector/modbus.go`
- Create: `edge-agent/internal/collector/modbus_test.go`

- [ ] **Step 1: Write Modbus collector**

```go
// edge-agent/internal/collector/modbus.go
package collector

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/goburrow/modbus"
)

type ModbusClient interface {
	Connect() error
	Close() error
	ReadHoldingRegisters(address, quantity uint16) ([]byte, error)
	ReadInputRegisters(address, quantity uint16) ([]byte, error)
}

type ModbusCollector struct {
	client ModbusClient
	slaveID byte
}

func NewModbusCollector(address string, slaveID byte) (*ModbusCollector, error) {
	handler := modbus.NewTCPClientHandler(address)
	handler.SlaveId = slaveID
	handler.Timeout = 10 * time.Second

	return &ModbusCollector{
		client:  handler,
		slaveID: slaveID,
	}, nil
}

func (c *ModbusCollector) Connect() error {
	return c.client.Connect()
}

func (c *ModbusCollector) Close() error {
	return c.client.Close()
}

func (c *ModbusCollector) ReadTemperature(address uint16, scale float64) (float64, error) {
	data, err := c.client.ReadInputRegisters(address, 1)
	if err != nil {
		return 0, fmt.Errorf("failed to read temperature: %w", err)
	}
	value := binary.BigEndian.Uint16(data)
	return float64(value) * scale, nil
}

func (c *ModbusCollector) ReadFloat32(address uint16) (float64, error) {
	data, err := c.client.ReadInputRegisters(address, 2)
	if err != nil {
		return 0, fmt.Errorf("failed to read float32: %w", err)
	}
	bits := binary.BigEndian.Uint32(data)
	return float64(bits), nil
}

type CoolingDeviceReader struct {
	collector *ModbusCollector
	config    DeviceConfig
}

type DeviceConfig struct {
	Name           string
	DeviceType     string
	RegisterMap    map[string]RegisterDef
}

type RegisterDef struct {
	Address  uint16
	DataType string
	Scale    float64
	Unit     string
}

func NewCoolingDeviceReader(collector *ModbusCollector, config DeviceConfig) *CoolingDeviceReader {
	return &CoolingDeviceReader{
		collector: collector,
		config:    config,
	}
}

func (r *CoolingDeviceReader) ReadAll() (map[string]float64, error) {
	results := make(map[string]float64)
	for metric, regDef := range r.config.RegisterMap {
		value, err := r.readRegister(regDef)
		if err != nil {
			return nil, fmt.Errorf("failed to read %s: %w", metric, err)
		}
		results[metric] = value
	}
	return results, nil
}

func (r *CoolingDeviceReader) readRegister(def RegisterDef) (float64, error) {
	switch def.DataType {
	case "uint16":
		data, err := r.collector.client.ReadInputRegisters(def.Address, 1)
		if err != nil {
			return 0, err
		}
		return float64(binary.BigEndian.Uint16(data)) * def.Scale, nil
	case "float32":
		data, err := r.collector.client.ReadInputRegisters(def.Address, 2)
		if err != nil {
			return 0, err
		}
		bits := binary.BigEndian.Uint32(data)
		return float64(bits) * def.Scale, nil
	default:
		return 0, fmt.Errorf("unsupported data type: %s", def.DataType)
	}
}
```

- [ ] **Step 2: Commit**

```bash
git add edge-agent/internal/collector/
git commit -m "feat(edge): add Modbus collector for cooling devices"
```

---

## Phase 5: Cloud Backend Services

### Task 7: Asset Service - Database Connection

**Files:**
- Create: `cloud-backend/shared/db/postgres.go`
- Create: `cloud-backend/asset-service/internal/models/server.go`

- [ ] **Step 1: Create database connection pool**

```go
// cloud-backend/shared/db/postgres.go
package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresDB struct {
	Pool *pgxpool.Pool
}

func NewPostgresDB(connString string) (*PostgresDB, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	config.MaxConns = 20
	config.MinConns = 5
	config.MaxConnLifetime = time.Hour

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresDB{Pool: pool}, nil
}

func (db *PostgresDB) Close() {
	db.Pool.Close()
}
```

```go
// cloud-backend/asset-service/internal/models/server.go
package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Server struct {
	ID               uuid.UUID `json:"id"`
	AgentID          uuid.UUID `json:"agent_id"`
	RedfishEndpoint  string    `json:"redfish_endpoint"`
	Manufacturer     string    `json:"manufacturer"`
	Model            string    `json:"model"`
	SerialNumber     string    `json:"serial_number"`
	CPUCount         int       `json:"cpu_count"`
	CPUModel         string    `json:"cpu_model"`
	MemoryGB         int       `json:"memory_gb"`
	Status           string    `json:"status"`
	PowerState       string    `json:"power_state"`
	Health           string    `json:"health"`
	DiscoveredAt     time.Time `json:"discovered_at"`
	LastUpdated      time.Time `json:"last_updated"`
}

type ServerRepository struct {
	db *pgxpool.Pool
}

func NewServerRepository(db *pgxpool.Pool) *ServerRepository {
	return &ServerRepository{db: db}
}

func (r *ServerRepository) CreateOrUpdate(ctx context.Context, server *Server) error {
	query := `
		INSERT INTO servers (agent_id, redfish_endpoint, manufacturer, model, serial_number,
			cpu_count, cpu_model, memory_gb, status, power_state, health, raw_data)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		ON CONFLICT (serial_number) DO UPDATE SET
			agent_id = EXCLUDED.agent_id,
			redfish_endpoint = EXCLUDED.redfish_endpoint,
			status = EXCLUDED.status,
			power_state = EXCLUDED.power_state,
			health = EXCLUDED.health,
			last_updated = NOW()
		RETURNING id
	`

	rawData := pgtype.JSONB{Valid: true}

	err := r.db.QueryRow(ctx, query,
		server.AgentID,
		server.RedfishEndpoint,
		server.Manufacturer,
		server.Model,
		server.SerialNumber,
		server.CPUCount,
		server.CPUModel,
		server.MemoryGB,
		server.Status,
		server.PowerState,
		server.Health,
		rawData,
	).Scan(&server.ID)

	return err
}

func (r *ServerRepository) ListByAgent(ctx context.Context, agentID uuid.UUID) ([]Server, error) {
	query := `
		SELECT id, agent_id, redfish_endpoint, manufacturer, model, serial_number,
			cpu_count, cpu_model, memory_gb, status, power_state, health, discovered_at, last_updated
		FROM servers
		WHERE agent_id = $1
		ORDER BY discovered_at DESC
	`

	rows, err := r.db.Query(ctx, query, agentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var servers []Server
	for rows.Next() {
		var s Server
		if err := rows.Scan(
			&s.ID, &s.AgentID, &s.RedfishEndpoint, &s.Manufacturer, &s.Model,
			&s.SerialNumber, &s.CPUCount, &s.CPUModel, &s.MemoryGB,
			&s.Status, &s.PowerState, &s.Health, &s.DiscoveredAt, &s.LastUpdated,
		); err != nil {
			return nil, err
		}
		servers = append(servers, s)
	}

	return servers, rows.Err()
}
```

- [ ] **Step 2: Commit**

```bash
git add cloud-backend/shared/db/ cloud-backend/asset-service/internal/models/
git commit -m "feat(cloud): add database connection and server models"
```

---

## Phase 6: Web Console

### Task 8: React Frontend Setup

**Files:**
- Create: `web-console/package.json`
- Create: `web-console/vite.config.ts`
- Create: `web-console/src/main.tsx`

- [ ] **Step 1: Initialize React project with Vite**

```json
{
  "name": "cloudmanager-web-console",
  "private": true,
  "version": "0.1.0",
  "type": "module",
  "scripts": {
    "dev": "vite",
    "build": "tsc && vite build",
    "preview": "vite preview",
    "lint": "eslint . --ext ts,tsx --report-unused-disable-directives --max-warnings 0"
  },
  "dependencies": {
    "react": "^18.2.0",
    "react-dom": "^18.2.0",
    "react-router-dom": "^6.23.0",
    "@tanstack/react-query": "^5.32.0",
    "axios": "^1.6.8",
    "recharts": "^2.12.6",
    "lucide-react": "^0.378.0"
  },
  "devDependencies": {
    "@types/react": "^18.2.66",
    "@types/react-dom": "^18.2.22",
    "@typescript-eslint/eslint-plugin": "^7.2.0",
    "@typescript-eslint/parser": "^7.2.0",
    "@vitejs/plugin-react": "^4.2.1",
    "autoprefixer": "^10.4.19",
    "eslint": "^8.57.0",
    "eslint-plugin-react-hooks": "^4.6.0",
    "eslint-plugin-react-refresh": "^0.4.6",
    "postcss": "^8.4.38",
    "tailwindcss": "^3.4.3",
    "typescript": "^5.2.2",
    "vite": "^5.2.0"
  }
}
```

```typescript
// web-console/vite.config.ts
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    port: 3000,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
})
```

```typescript
// web-console/src/main.tsx
import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App.tsx'
import './index.css'

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
)
```

- [ ] **Step 2: Install dependencies**

```bash
cd /Users/arthurzhang/dev/llm/cloudManager/web-console
npm install
```

- [ ] **Step 3: Commit**

```bash
git add web-console/
git commit -m "feat(web): initialize React frontend with Vite and Tailwind"
```

---

## Summary

This plan covers:
- **Phase 1:** Project setup, Docker infrastructure
- **Phase 2:** Database schema with TimescaleDB
- **Phase 3:** gRPC protocol definitions
- **Phase 4:** Edge Agent (config, Modbus collector)
- **Phase 5:** Cloud Backend (database, models)
- **Phase 6:** Web Console (React setup)

### Execution Options:

**1. Subagent-Driven (recommended)** - Dispatch a fresh subagent per task, review between tasks, fast iteration

**2. Inline Execution** - Execute tasks in this session using executing-plans, batch execution with checkpoints

Which approach do you prefer? Or would you like me to expand the plan with more detailed tasks before starting implementation?
