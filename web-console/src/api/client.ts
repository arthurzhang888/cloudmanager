import axios from 'axios';
import type {
  DataCenter,
  EdgeAgent,
  Server,
  CoolingDevice,
  SensorReading,
  Alert,
  PueCalculation,
  DashboardStats,
} from '../types';

const API_BASE_URL = import.meta.env.VITE_API_URL || '/api';

const client = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Dashboard
export const getDashboardStats = () =>
  client.get<DashboardStats>('/dashboard/stats').then((r) => r.data);

// Data Centers
export const getDataCenters = () =>
  client.get<DataCenter[]>('/datacenters').then((r) => r.data);

export const getDataCenter = (id: string) =>
  client.get<DataCenter>(`/datacenters/${id}`).then((r) => r.data);

// Edge Agents
export const getEdgeAgents = () =>
  client.get<EdgeAgent[]>('/agents').then((r) => r.data);

export const getEdgeAgent = (id: string) =>
  client.get<EdgeAgent>(`/agents/${id}`).then((r) => r.data);

// Servers
export const getServers = (params?: { agent_id?: string; status?: string }) =>
  client.get<Server[]>('/servers', { params }).then((r) => r.data);

export const getServer = (id: string) =>
  client.get<Server>(`/servers/${id}`).then((r) => r.data);

// Cooling Devices
export const getCoolingDevices = () =>
  client.get<CoolingDevice[]>('/cooling-devices').then((r) => r.data);

export const getCoolingDevice = (id: string) =>
  client.get<CoolingDevice>(`/cooling-devices/${id}`).then((r) => r.data);

// Sensor Readings
export const getSensorReadings = (params: {
  device_id?: string;
  metric_name?: string;
  start_time?: string;
  end_time?: string;
}) =>
  client.get<SensorReading[]>('/readings', { params }).then((r) => r.data);

// Alerts
export const getAlerts = (params?: { status?: string; severity?: string }) =>
  client.get<Alert[]>('/alerts', { params }).then((r) => r.data);

export const acknowledgeAlert = (id: string) =>
  client.post(`/alerts/${id}/acknowledge`).then((r) => r.data);

// PUE
export const getPueHistory = (dc_id: string, hours: number = 24) =>
  client
    .get<PueCalculation[]>(`/pue`, { params: { dc_id, hours } })
    .then((r) => r.data);

export default client;
