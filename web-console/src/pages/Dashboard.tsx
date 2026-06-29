import { Server, Thermometer, AlertCircle, Activity } from 'lucide-react'

function Dashboard() {
  const stats = [
    { label: 'Total Servers', value: '1,234', icon: Server, color: 'blue' },
    { label: 'Cooling Devices', value: '48', icon: Thermometer, color: 'green' },
    { label: 'Active Alerts', value: '12', icon: AlertCircle, color: 'red' },
    { label: 'Avg PUE', value: '1.35', icon: Activity, color: 'purple' },
  ]

  return (
    <div>
      <h2 className="text-2xl font-bold text-gray-800 mb-6">Dashboard</h2>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        {stats.map((stat) => {
          const Icon = stat.icon
          return (
            <div key={stat.label} className="bg-white rounded-lg shadow p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-gray-500">{stat.label}</p>
                  <p className="text-2xl font-bold text-gray-800">{stat.value}</p>
                </div>
                <div className={`p-3 rounded-full bg-${stat.color}-100`}>
                  <Icon className={`text-${stat.color}-600`} size={24} />
                </div>
              </div>
            </div>
          )
        })}
      </div>

      {/* Recent Alerts */}
      <div className="bg-white rounded-lg shadow">
        <div className="p-4 border-b">
          <h3 className="text-lg font-semibold text-gray-800">Recent Alerts</h3>
        </div>
        <div className="p-4">
          <p className="text-gray-500">No recent alerts</p>
        </div>
      </div>
    </div>
  )
}

export default Dashboard
