import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import {
  Thermometer,
  Droplets,
  Gauge,
  Activity,
  AlertTriangle,
} from 'lucide-react';
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  Legend,
} from 'recharts';
import { mockCoolingDevices, generateMockReadings } from '../api/mock';
import type { CoolingDevice } from '../types';

function Cooling() {
  const [selectedDevice, setSelectedDevice] = useState<string | null>(null);

  const { data: devices } = useQuery({
    queryKey: ['coolingDevices'],
    queryFn: () => mockCoolingDevices,
  });

  const supplyTempData = generateMockReadings('cdu-001', 'supply_temp', 4);
  const returnTempData = generateMockReadings('cdu-001', 'return_temp', 4);
  const flowRateData = generateMockReadings('cdu-001', 'flow_rate', 4);

  const combinedTempData = supplyTempData.map((item, index) => ({
    time: item.time,
    supply: item.value,
    return: returnTempData[index]?.value || 0,
  }));

  const getStatusColor = (status: CoolingDevice['status']) => {
    switch (status) {
      case 'online':
        return 'bg-green-100 text-green-800';
      case 'offline':
        return 'bg-red-100 text-red-800';
      case 'error':
        return 'bg-red-100 text-red-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  const getDeviceIcon = (type: CoolingDevice['device_type']) => {
    if (type === 'cdu') {
      return <Droplets className="text-blue-500" size={24} />;
    }
    return <Thermometer className="text-cyan-500" size={24} />;
  };

  return (
    <div className="space-y-6">
      <h2 className="text-2xl font-bold text-gray-800">Liquid Cooling</h2>

      {/* Devices Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {devices?.map((device) => (
          <div
            key={device.id}
            className={`bg-white rounded-lg shadow p-4 cursor-pointer transition-all ${
              selectedDevice === device.id ? 'ring-2 ring-blue-500' : 'hover:shadow-md'
            }`}
            onClick={() => setSelectedDevice(device.id)}
          >
            <div className="flex items-start justify-between">
              <div className="flex items-center gap-3">
                {getDeviceIcon(device.device_type)}
                <div>
                  <h3 className="font-semibold text-gray-800">{device.name}</h3>
                  <p className="text-sm text-gray-500">{device.location}</p>
                </div>
              </div>
              <span
                className={`px-2 py-1 text-xs font-medium rounded-full ${getStatusColor(
                  device.status
                )}`}
              >
                {device.status}
              </span>
            </div>
            <div className="mt-4 pt-4 border-t">
              <div className="flex items-center justify-between text-sm">
                <span className="text-gray-500">Type</span>
                <span className="font-medium uppercase">{device.device_type}</span>
              </div>
              <div className="flex items-center justify-between text-sm mt-1">
                <span className="text-gray-500">Modbus</span>
                <span className="font-mono text-xs">{device.modbus_address}</span>
              </div>
            </div>
          </div>
        ))}
      </div>

      {/* Metrics */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div className="bg-white rounded-lg shadow p-4">
          <div className="flex items-center gap-2 mb-2">
            <Thermometer className="text-blue-500" size={20} />
            <span className="text-sm text-gray-500">Supply Temp</span>
          </div>
          <p className="text-2xl font-bold text-gray-800">
            {supplyTempData[supplyTempData.length - 1]?.value.toFixed(1)}°C
          </p>
        </div>
        <div className="bg-white rounded-lg shadow p-4">
          <div className="flex items-center gap-2 mb-2">
            <Thermometer className="text-red-500" size={20} />
            <span className="text-sm text-gray-500">Return Temp</span>
          </div>
          <p className="text-2xl font-bold text-gray-800">
            {returnTempData[returnTempData.length - 1]?.value.toFixed(1)}°C
          </p>
        </div>
        <div className="bg-white rounded-lg shadow p-4">
          <div className="flex items-center gap-2 mb-2">
            <Activity className="text-green-500" size={20} />
            <span className="text-sm text-gray-500">Flow Rate</span>
          </div>
          <p className="text-2xl font-bold text-gray-800">
            {flowRateData[flowRateData.length - 1]?.value.toFixed(0)} L/min
          </p>
        </div>
        <div className="bg-white rounded-lg shadow p-4">
          <div className="flex items-center gap-2 mb-2">
            <Gauge className="text-purple-500" size={20} />
            <span className="text-sm text-gray-500">Pressure</span>
          </div>
          <p className="text-2xl font-bold text-gray-800">2.3 bar</p>
        </div>
      </div>

      {/* Temperature Chart */}
      <div className="bg-white rounded-lg shadow p-6">
        <h3 className="text-lg font-semibold text-gray-800 mb-4">Temperature Trend</h3>
        <div className="h-64">
          <ResponsiveContainer width="100%" height="100%">
            <LineChart data={combinedTempData}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis
                dataKey="time"
                tickFormatter={(time) =>
                  new Date(time).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
                }
              />
              <YAxis domain={[15, 35]} />
              <Tooltip
                formatter={(value: number) => value.toFixed(2) + '°C'}
                labelFormatter={(label) => new Date(label).toLocaleString()}
              />
              <Legend />
              <Line
                type="monotone"
                dataKey="supply"
                stroke="#3b82f6"
                name="Supply Temp"
                strokeWidth={2}
              />
              <Line
                type="monotone"
                dataKey="return"
                stroke="#ef4444"
                name="Return Temp"
                strokeWidth={2}
              />
            </LineChart>
          </ResponsiveContainer>
        </div>
      </div>

      {/* Alerts */}
      <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-4">
        <div className="flex items-start gap-3">
          <AlertTriangle className="text-yellow-600 flex-shrink-0" size={20} />
          <div>
            <h4 className="font-medium text-yellow-800">Temperature Threshold Alert</h4>
            <p className="text-sm text-yellow-700">
              Supply temperature (25.2°C) is approaching the upper limit (26°C) for CDU-Rack-A01.
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}

export default Cooling;
