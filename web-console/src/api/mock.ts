import type {
  DashboardStats,
  DataCenter,
  EdgeAgent,
  Server,
  CoolingDevice,
  Alert,
  SensorReading,
  PueCalculation,
} from '../types';

// Mock Dashboard Stats
export const mockDashboardStats: DashboardStats = {
  total_servers: 1247,
  online_servers: 1189,
  offline_servers: 58,
  cooling_devices: 48,
  active_alerts: 12,
  avg_pue: 1.35,
};

// Mock Data Centers
export const mockDataCenters: DataCenter[] = [
  {
    id: 'dc-1',
    name: 'Beijing DC',
    location: 'Beijing, China',
    network_type: 'vpn',
    created_at: '2024-01-15T08:00:00Z',
    updated_at: '2024-06-29T10:00:00Z',
  },
  {
    id: 'dc-2',
    name: 'Shanghai DC',
    location: 'Shanghai, China',
    network_type: 'public',
    created_at: '2024-02-20T08:00:00Z',
    updated_at: '2024-06-28T15:30:00Z',
  },
  {
    id: 'dc-3',
    name: 'Shenzhen DC',
    location: 'Shenzhen, China',
    network_type: 'vpn',
    created_at: '2024-03-10T08:00:00Z',
    updated_at: '2024-06-29T08:15:00Z',
  },
];

// Mock Edge Agents
export const mockEdgeAgents: EdgeAgent[] = [
  {
    id: 'agent-001',
    dc_id: 'dc-1',
    hostname: 'edge-beijing-01',
    version: '0.1.0',
    last_seen: new Date().toISOString(),
    status: 'online',
    config: {},
    created_at: '2024-01-15T08:00:00Z',
    updated_at: '2024-06-29T10:00:00Z',
  },
  {
    id: 'agent-002',
    dc_id: 'dc-2',
    hostname: 'edge-shanghai-01',
    version: '0.1.0',
    last_seen: new Date(Date.now() - 5 * 60 * 1000).toISOString(),
    status: 'online',
    config: {},
    created_at: '2024-02-20T08:00:00Z',
    updated_at: '2024-06-28T15:30:00Z',
  },
  {
    id: 'agent-003',
    dc_id: 'dc-3',
    hostname: 'edge-shenzhen-01',
    version: '0.1.0',
    last_seen: new Date(Date.now() - 2 * 60 * 60 * 1000).toISOString(),
    status: 'offline',
    config: {},
    created_at: '2024-03-10T08:00:00Z',
    updated_at: '2024-06-29T08:15:00Z',
  },
];

// Mock Servers
export const mockServers: Server[] = [
  {
    id: 'srv-001',
    agent_id: 'agent-001',
    redfish_endpoint: 'https://192.168.1.10',
    manufacturer: 'Supermicro',
    model: 'SuperServer 2029TP-HTR',
    serial_number: 'SM2029A12345',
    cpu_count: 2,
    cpu_model: 'Intel Xeon Gold 6248R',
    memory_gb: 256,
    status: 'online',
    power_state: 'On',
    health: 'OK',
    discovered_at: '2024-01-15T08:30:00Z',
    last_updated: new Date().toISOString(),
  },
  {
    id: 'srv-002',
    agent_id: 'agent-001',
    redfish_endpoint: 'https://192.168.1.11',
    manufacturer: 'Supermicro',
    model: 'SuperServer 2029TP-HTR',
    serial_number: 'SM2029A12346',
    cpu_count: 2,
    cpu_model: 'Intel Xeon Gold 6248R',
    memory_gb: 256,
    status: 'online',
    power_state: 'On',
    health: 'OK',
    discovered_at: '2024-01-15T08:30:00Z',
    last_updated: new Date().toISOString(),
  },
  {
    id: 'srv-003',
    agent_id: 'agent-001',
    redfish_endpoint: 'https://192.168.1.12',
    manufacturer: 'Dell',
    model: 'PowerEdge R750',
    serial_number: 'DELL750B78901',
    cpu_count: 2,
    cpu_model: 'Intel Xeon Gold 6330',
    memory_gb: 512,
    status: 'offline',
    power_state: 'Off',
    health: 'Warning',
    discovered_at: '2024-01-15T08:30:00Z',
    last_updated: new Date(Date.now() - 24 * 60 * 60 * 1000).toISOString(),
  },
];

// Mock Cooling Devices
export const mockCoolingDevices: CoolingDevice[] = [
  {
    id: 'cdu-001',
    agent_id: 'agent-001',
    device_type: 'cdu',
    name: 'CDU-Rack-A01',
    modbus_address: 'tcp://192.168.10.10:502',
    slave_id: 1,
    location: 'Rack A01-A10',
    status: 'online',
    discovered_at: '2024-01-15T08:30:00Z',
    last_updated: new Date().toISOString(),
  },
  {
    id: 'cdu-002',
    agent_id: 'agent-001',
    device_type: 'cdu',
    name: 'CDU-Rack-A11',
    modbus_address: 'tcp://192.168.10.11:502',
    slave_id: 1,
    location: 'Rack A11-A20',
    status: 'online',
    discovered_at: '2024-01-15T08:30:00Z',
    last_updated: new Date().toISOString(),
  },
  {
    id: 'tower-001',
    agent_id: 'agent-001',
    device_type: 'cooling_tower',
    name: 'Cooling Tower 1',
    modbus_address: 'tcp://192.168.10.100:502',
    slave_id: 1,
    location: 'Roof',
    status: 'online',
    discovered_at: '2024-01-15T08:30:00Z',
    last_updated: new Date().toISOString(),
  },
];

// Mock Alerts
export const mockAlerts: Alert[] = [
  {
    id: 'alert-001',
    severity: 'warning',
    category: 'cooling',
    source_type: 'cdu',
    source_id: 'cdu-001',
    title: 'CDU Supply Temperature High',
    description: 'Supply temperature exceeded 25°C threshold',
    status: 'active',
    created_at: new Date(Date.now() - 2 * 60 * 60 * 1000).toISOString(),
  },
  {
    id: 'alert-002',
    severity: 'critical',
    category: 'hardware',
    source_type: 'server',
    source_id: 'srv-003',
    title: 'Server Offline',
    description: 'Server has been offline for more than 24 hours',
    status: 'active',
    created_at: new Date(Date.now() - 24 * 60 * 60 * 1000).toISOString(),
  },
  {
    id: 'alert-003',
    severity: 'info',
    category: 'threshold',
    source_type: 'cdu',
    source_id: 'cdu-002',
    title: 'PUE Above Target',
    description: 'Current PUE (1.42) exceeds target (1.35)',
    status: 'acknowledged',
    created_at: new Date(Date.now() - 6 * 60 * 60 * 1000).toISOString(),
  },
];

// Generate mock sensor readings
export const generateMockReadings = (
  deviceId: string,
  metricName: string,
  hours: number = 24
): SensorReading[] => {
  const readings: SensorReading[] = [];
  const now = Date.now();
  const interval = 5 * 60 * 1000; // 5 minutes

  for (let i = 0; i < (hours * 60 * 60 * 1000) / interval; i++) {
    const time = new Date(now - i * interval).toISOString();
    let value: number;
    let unit: string;

    switch (metricName) {
      case 'supply_temp':
        value = 18 + Math.random() * 4;
        unit = '°C';
        break;
      case 'return_temp':
        value = 25 + Math.random() * 5;
        unit = '°C';
        break;
      case 'flow_rate':
        value = 100 + Math.random() * 20;
        unit = 'L/min';
        break;
      case 'pressure':
        value = 2 + Math.random() * 0.5;
        unit = 'bar';
        break;
      case 'power':
        value = 50 + Math.random() * 30;
        unit = 'kW';
        break;
      default:
        value = Math.random() * 100;
        unit = '';
    }

    readings.push({
      time,
      device_id: deviceId,
      device_type: 'cdu',
      metric_name: metricName,
      value: parseFloat(value.toFixed(2)),
      unit,
    });
  }

  return readings.reverse();
};

// Generate mock PUE history
export const generateMockPueHistory = (hours: number = 24): PueCalculation[] => {
  const history: PueCalculation[] = [];
  const now = Date.now();
  const interval = 60 * 60 * 1000; // 1 hour

  for (let i = 0; i < hours; i++) {
    const time = new Date(now - i * interval).toISOString();
    history.push({
      time,
      dc_id: 'dc-1',
      pue_value: parseFloat((1.3 + Math.random() * 0.15).toFixed(3)),
      it_power_kw: parseFloat((800 + Math.random() * 200).toFixed(2)),
      facility_power_kw: parseFloat((1100 + Math.random() * 250).toFixed(2)),
      cooling_power_kw: parseFloat((200 + Math.random() * 50).toFixed(2)),
    });
  }

  return history.reverse();
};
