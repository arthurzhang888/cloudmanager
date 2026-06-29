-- migrations/001_initial.sql

-- Enable TimescaleDB extension
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
CREATE INDEX idx_servers_serial ON servers(serial_number);

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
CREATE INDEX idx_cooling_devices_type ON cooling_devices(device_type);

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
CREATE INDEX idx_sensor_readings_agent ON sensor_readings(agent_id, time DESC);

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
CREATE INDEX idx_alerts_severity ON alerts(severity);
CREATE INDEX idx_alerts_created ON alerts(created_at DESC);

-- Update trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_data_centers_updated_at BEFORE UPDATE ON data_centers
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_edge_agents_updated_at BEFORE UPDATE ON edge_agents
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
