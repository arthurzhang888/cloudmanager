// Data Center
export interface DataCenter {
  id: string;
  name: string;
  location: string;
  network_type: 'vpn' | 'public' | 'airgap';
  created_at: string;
  updated_at: string;
}

// Edge Agent
export interface EdgeAgent {
  id: string;
  dc_id: string;
  hostname: string;
  version: string;
  last_seen: string | null;
  status: 'online' | 'offline' | 'error';
  config: Record<string, unknown>;
  created_at: string;
  updated_at: string;
}

// Server
export interface Server {
  id: string;
  agent_id: string;
  redfish_endpoint: string;
  manufacturer: string;
  model: string;
  serial_number: string;
  cpu_count: number;
  cpu_model: string;
  memory_gb: number;
  status: 'online' | 'offline' | 'error' | 'maintenance' | 'unknown';
  power_state: 'On' | 'Off' | string;
  health: 'OK' | 'Warning' | 'Critical' | string;
  discovered_at: string;
  last_updated: string;
}

// Cooling Device
export interface CoolingDevice {
  id: string;
  agent_id: string;
  device_type: 'cdu' | 'cooling_tower';
  name: string;
  modbus_address: string;
  slave_id: number;
  location: string;
  status: 'online' | 'offline' | 'error' | 'unknown';
  discovered_at: string;
  last_updated: string;
}

// Sensor Reading
export interface SensorReading {
  time: string;
  device_id: string;
  device_type: string;
  metric_name: string;
  value: number;
  unit: string;
}

// Alert
export interface Alert {
  id: string;
  severity: 'critical' | 'warning' | 'info';
  category: 'hardware' | 'cooling' | 'network' | 'threshold';
  source_type: string;
  source_id: string;
  title: string;
  description: string;
  status: 'active' | 'acknowledged' | 'resolved';
  created_at: string;
}

// PUE Calculation
export interface PueCalculation {
  time: string;
  dc_id: string;
  pue_value: number;
  it_power_kw: number;
  facility_power_kw: number;
  cooling_power_kw: number;
}

// Dashboard Stats
export interface DashboardStats {
  total_servers: number;
  online_servers: number;
  offline_servers: number;
  cooling_devices: number;
  active_alerts: number;
  avg_pue: number;
}
