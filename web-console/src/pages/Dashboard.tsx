import { useQuery } from '@tanstack/react-query';
import {
  Server,
  Thermometer,
  AlertCircle,
  Activity,
  CheckCircle,
  XCircle,
  Droplets,
} from 'lucide-react';
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  AreaChart,
  Area,
  BarChart,
  Bar,
  PieChart,
  Pie,
  Cell,
  Legend,
} from 'recharts';
import {
  mockDashboardStats,
  mockAlerts,
  generateMockPueHistory,
  generateMockTemperatureHistory,
  generateMockAlertHistory,
  generateMockServerUtilization,
  generateMockCoolingEfficiency,
  generateMockServerHealthDistribution,
} from '../api/mock';
import type { Alert } from '../types';

function Dashboard() {
  const { data: stats } = useQuery({
    queryKey: ['dashboardStats'],
    queryFn: () => mockDashboardStats,
  });

  const pueHistory = generateMockPueHistory(24);
  const temperatureHistory = generateMockTemperatureHistory(24);
  const alertHistory = generateMockAlertHistory(7);
  const serverUtilization = generateMockServerUtilization(24);
  const coolingEfficiency = generateMockCoolingEfficiency();
  const serverHealthDistribution = generateMockServerHealthDistribution();

  const statCards = [
    {
      label: 'Total Servers',
      value: stats?.total_servers || 0,
      icon: Server,
      color: 'blue',
    },
    {
      label: 'Online',
      value: stats?.online_servers || 0,
      icon: CheckCircle,
      color: 'green',
    },
    {
      label: 'Offline',
      value: stats?.offline_servers || 0,
      icon: XCircle,
      color: 'red',
    },
    {
      label: 'Cooling Devices',
      value: stats?.cooling_devices || 0,
      icon: Thermometer,
      color: 'cyan',
    },
    {
      label: 'Active Alerts',
      value: stats?.active_alerts || 0,
      icon: AlertCircle,
      color: 'orange',
    },
    {
      label: 'Avg PUE',
      value: stats?.avg_pue?.toFixed(2) || '0.00',
      icon: Activity,
      color: 'purple',
    },
  ];

  const getSeverityColor = (severity: Alert['severity']) => {
    switch (severity) {
      case 'critical':
        return 'text-red-600 bg-red-50';
      case 'warning':
        return 'text-yellow-600 bg-yellow-50';
      case 'info':
        return 'text-blue-600 bg-blue-50';
      default:
        return 'text-gray-600 bg-gray-50';
    }
  };

  const getStatusColor = (status: Alert['status']) => {
    switch (status) {
      case 'active':
        return 'text-red-600';
      case 'acknowledged':
        return 'text-yellow-600';
      case 'resolved':
        return 'text-green-600';
      default:
        return 'text-gray-600';
    }
  };

  return (
    <div className="space-y-6">
      <h2 className="text-2xl font-bold text-gray-800">Dashboard</h2>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-6 gap-4">
        {statCards.map((stat) => {
          const Icon = stat.icon;
          return (
            <div key={stat.label} className="bg-white rounded-lg shadow p-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-gray-500">{stat.label}</p>
                  <p className="text-2xl font-bold text-gray-800">{stat.value}</p>
                </div>
                <div className={`p-2 rounded-full bg-${stat.color}-100`}>
                  <Icon className={`text-${stat.color}-600`} size={20} />
                </div>
              </div>
            </div>
          );
        })}
      </div>

      {/* Charts Row */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* PUE Chart */}
        <div className="bg-white rounded-lg shadow p-6">
          <h3 className="text-lg font-semibold text-gray-800 mb-4">PUE Trend (24h)</h3>
          <div className="h-64">
            <ResponsiveContainer width="100%" height="100%">
              <AreaChart data={pueHistory}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis
                  dataKey="time"
                  tickFormatter={(time) => new Date(time).getHours() + ':00'}
                />
                <YAxis domain={[1.2, 1.6]} />
                <Tooltip
                  formatter={(value: number) => value.toFixed(3)}
                  labelFormatter={(label) => new Date(label).toLocaleString()}
                />
                <Area
                  type="monotone"
                  dataKey="pue_value"
                  stroke="#8b5cf6"
                  fill="#8b5cf6"
                  fillOpacity={0.2}
                />
              </AreaChart>
            </ResponsiveContainer>
          </div>
        </div>

        {/* Power Chart */}
        <div className="bg-white rounded-lg shadow p-6">
          <h3 className="text-lg font-semibold text-gray-800 mb-4">Power Consumption (24h)</h3>
          <div className="h-64">
            <ResponsiveContainer width="100%" height="100%">
              <LineChart data={pueHistory}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis
                  dataKey="time"
                  tickFormatter={(time) => new Date(time).getHours() + ':00'}
                />
                <YAxis />
                <Tooltip
                  formatter={(value: number) => value.toFixed(2) + ' kW'}
                  labelFormatter={(label) => new Date(label).toLocaleString()}
                />
                <Line
                  type="monotone"
                  dataKey="it_power_kw"
                  stroke="#3b82f6"
                  name="IT Power"
                />
                <Line
                  type="monotone"
                  dataKey="cooling_power_kw"
                  stroke="#06b6d4"
                  name="Cooling Power"
                />
              </LineChart>
            </ResponsiveContainer>
          </div>
        </div>
      </div>

      {/* Temperature & Utilization Row */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Temperature Chart */}
        <div className="bg-white rounded-lg shadow p-6">
          <div className="flex items-center gap-2 mb-4">
            <Thermometer className="text-red-500" size={20} />
            <h3 className="text-lg font-semibold text-gray-800">Cooling Temperature (24h)</h3>
          </div>
          <div className="h-64">
            <ResponsiveContainer width="100%" height="100%">
              <LineChart data={temperatureHistory}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis
                  dataKey="time"
                  tickFormatter={(time) => new Date(time).getHours() + ':00'}
                />
                <YAxis domain={[15, 35]} />
                <Tooltip
                  formatter={(value: number) => value.toFixed(1) + '°C'}
                  labelFormatter={(label) => new Date(label).toLocaleString()}
                />
                <Legend />
                <Line
                  type="monotone"
                  dataKey="supply_temp"
                  stroke="#3b82f6"
                  name="Supply Temp"
                  strokeWidth={2}
                />
                <Line
                  type="monotone"
                  dataKey="return_temp"
                  stroke="#ef4444"
                  name="Return Temp"
                  strokeWidth={2}
                />
                <Line
                  type="monotone"
                  dataKey="ambient_temp"
                  stroke="#10b981"
                  name="Ambient Temp"
                  strokeWidth={2}
                  strokeDasharray="5 5"
                />
              </LineChart>
            </ResponsiveContainer>
          </div>
        </div>

        {/* Server Utilization Chart */}
        <div className="bg-white rounded-lg shadow p-6">
          <div className="flex items-center gap-2 mb-4">
            <Activity className="text-blue-500" size={20} />
            <h3 className="text-lg font-semibold text-gray-800">Server Utilization (24h)</h3>
          </div>
          <div className="h-64">
            <ResponsiveContainer width="100%" height="100%">
              <AreaChart data={serverUtilization}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis
                  dataKey="time"
                  tickFormatter={(time) => new Date(time).getHours() + ':00'}
                />
                <YAxis domain={[0, 100]} />
                <Tooltip
                  formatter={(value: number, name: string) => [
                    value.toFixed(1) + (name === 'Power' ? ' kW' : '%'),
                    name,
                  ]}
                  labelFormatter={(label) => new Date(label).toLocaleString()}
                />
                <Legend />
                <Area
                  type="monotone"
                  dataKey="cpu_percent"
                  stroke="#3b82f6"
                  fill="#3b82f6"
                  fillOpacity={0.2}
                  name="CPU"
                />
                <Area
                  type="monotone"
                  dataKey="memory_percent"
                  stroke="#8b5cf6"
                  fill="#8b5cf6"
                  fillOpacity={0.2}
                  name="Memory"
                />
              </AreaChart>
            </ResponsiveContainer>
          </div>
        </div>
      </div>

      {/* Analytics Row */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Alert History */}
        <div className="bg-white rounded-lg shadow p-6">
          <div className="flex items-center gap-2 mb-4">
            <AlertCircle className="text-orange-500" size={20} />
            <h3 className="text-lg font-semibold text-gray-800">Alert History (7d)</h3>
          </div>
          <div className="h-64">
            <ResponsiveContainer width="100%" height="100%">
              <BarChart data={alertHistory}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis
                  dataKey="time"
                  tickFormatter={(time) =>
                    new Date(time).toLocaleDateString('en-US', { weekday: 'short' })
                  }
                />
                <YAxis />
                <Tooltip
                  labelFormatter={(label) => new Date(label).toLocaleDateString()}
                />
                <Legend />
                <Bar dataKey="critical" fill="#ef4444" name="Critical" />
                <Bar dataKey="warning" fill="#f59e0b" name="Warning" />
                <Bar dataKey="info" fill="#3b82f6" name="Info" />
              </BarChart>
            </ResponsiveContainer>
          </div>
        </div>

        {/* Server Health Distribution */}
        <div className="bg-white rounded-lg shadow p-6">
          <div className="flex items-center gap-2 mb-4">
            <Server className="text-green-500" size={20} />
            <h3 className="text-lg font-semibold text-gray-800">Server Health</h3>
          </div>
          <div className="h-64">
            <ResponsiveContainer width="100%" height="100%">
              <PieChart>
                <Pie
                  data={serverHealthDistribution}
                  cx="50%"
                  cy="50%"
                  innerRadius={60}
                  outerRadius={80}
                  paddingAngle={5}
                  dataKey="value"
                >
                  {serverHealthDistribution.map((entry: { color: string }, index: number) => (
                    <Cell key={`cell-${index}`} fill={entry.color} />
                  ))}
                </Pie>
                <Tooltip formatter={(value: number) => value.toString() + ' servers'} />
                <Legend />
              </PieChart>
            </ResponsiveContainer>
          </div>
        </div>

        {/* Cooling Efficiency */}
        <div className="bg-white rounded-lg shadow p-6">
          <div className="flex items-center gap-2 mb-4">
            <Droplets className="text-cyan-500" size={20} />
            <h3 className="text-lg font-semibold text-gray-800">Cooling Efficiency</h3>
          </div>
          <div className="h-64">
            <ResponsiveContainer width="100%" height="100%">
              <BarChart data={coolingEfficiency} layout="vertical">
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis type="number" domain={[0, 100]} />
                <YAxis type="category" dataKey="name" width={100} />
                <Tooltip formatter={(value: number) => value.toFixed(1) + '%'} />
                <Bar dataKey="value" fill="#06b6d4" radius={[0, 4, 4, 0]} />
              </BarChart>
            </ResponsiveContainer>
          </div>
        </div>
      </div>

      {/* Recent Alerts */}
      <div className="bg-white rounded-lg shadow">
        <div className="p-4 border-b">
          <h3 className="text-lg font-semibold text-gray-800">Recent Alerts</h3>
        </div>
        <div className="divide-y">
          {mockAlerts.length > 0 ? (
            mockAlerts.map((alert) => (
              <div key={alert.id} className="p-4 hover:bg-gray-50">
                <div className="flex items-start justify-between">
                  <div className="flex items-start gap-3">
                    <span
                      className={`px-2 py-1 rounded text-xs font-medium ${getSeverityColor(
                        alert.severity
                      )}`}
                    >
                      {alert.severity.toUpperCase()}
                    </span>
                    <div>
                      <p className="font-medium text-gray-800">{alert.title}</p>
                      <p className="text-sm text-gray-500">{alert.description}</p>
                      <p className="text-xs text-gray-400 mt-1">
                        {new Date(alert.created_at).toLocaleString()}
                      </p>
                    </div>
                  </div>
                  <span className={`text-sm font-medium ${getStatusColor(alert.status)}`}>
                    {alert.status}
                  </span>
                </div>
              </div>
            ))
          ) : (
            <div className="p-4 text-gray-500">No recent alerts</div>
          )}
        </div>
      </div>
    </div>
  );
}

export default Dashboard;
