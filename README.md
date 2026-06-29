# CloudManager

Enterprise data center infrastructure management platform supporting multi-site large-scale data centers (1000+ servers) with unified management, composable infrastructure capabilities, and deep liquid cooling system integration. Similar to Supermicro SuperCloud Composer.

## Features

### Implemented

- **API Gateway** - RESTful HTTP API with JWT authentication, role-based access control, and comprehensive middleware
- **Alert Rule Engine** - Configurable alert rules with condition evaluation, duration-based triggering, and action execution
- **User Authentication** - Secure authentication with bcrypt password hashing and JWT tokens
- **Dashboard Visualizations** - Real-time charts for PUE, temperature, server utilization, alerts, and cooling efficiency
- **Edge Agent** - Data collection agent supporting Modbus, Redfish, and SNMP protocols
- **Asset Management** - Server and cooling device discovery and inventory
- **Telemetry** - Sensor data collection with time-series storage

### Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    Cloud Control Platform                │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐   │
│  │ API      │ │ Asset    │ │Telemetry │ │ Web      │   │
│  │ Gateway  │ │ Service  │ │ Service  │ │ Console  │   │
│  └─────┬────┘ └─────┬────┘ └─────┬────┘ └─────┬────┘   │
│        └─────────────┴─────────────┴─────────────┘       │
│                      │                                   │
│        ┌─────────────┴─────────────┐                    │
│        │      Message Queue        │                    │
│        │       (NATS/Kafka)        │                    │
│        └─────────────┬─────────────┘                    │
└──────────────────────┼──────────────────────────────────┘
                       │
         ┌─────────────┼─────────────┐
         │             │             │
    ┌────▼────┐   ┌────▼────┐   ┌────▼────┐
    │  Edge   │   │  Edge   │   │  Edge   │
    │ Agent   │   │ Agent   │   │ Agent   │
    │(DC-1)   │   │(DC-2)   │   │(DC-N)   │
    └────┬────┘   └────┬────┘   └────┬────┘
         │             │             │
    ┌────▼─────────────▼─────────────▼────┐
    │           Data Center Hardware       │
    │  ┌──────┐ ┌──────┐ ┌──────┐        │
    │  │Servers│ │Network│ │ CDU  │        │
    │  │Redfish│ │SNMP  │ │Modbus│        │
    │  └──────┘ └──────┘ └──────┘        │
    └─────────────────────────────────────┘
```

## Tech Stack

- **Edge Agent**: Go, gRPC, Modbus, Redfish, SNMP
- **Cloud Backend**: Go, Gin, gRPC, PostgreSQL, TimescaleDB, JWT
- **Frontend**: React, TypeScript, Vite, Tailwind CSS, Recharts
- **Deployment**: Docker, Docker Compose

## Quick Start

### Prerequisites
- Docker & Docker Compose
- Node.js 18+ (frontend development)
- Go 1.22+ (backend development)

### Development Setup

```bash
# 1. Clone repository
git clone https://github.com/arthurzhang888/cloudmanager.git
cd cloudmanager

# 2. Start infrastructure (database, message queue)
make dev-up

# 3. Run database migrations
make migrate

# 4. Start API Gateway
cd cloud-backend/api-gateway
go run cmd/main.go

# 5. Start frontend dev server (in another terminal)
make web-dev
```

Visit http://localhost:3000 for the Web Console.

### Production Deployment

```bash
# Build and start production environment
make prod-build
make prod-up
```

## Project Structure

```
cloudmanager/
├── edge-agent/              # Edge data collection agent
│   ├── cmd/agent/           # Main entry point
│   ├── internal/
│   │   ├── config/          # Configuration management
│   │   ├── collector/       # Data collection (Modbus)
│   │   └── discovery/       # Device discovery
│   └── proto/               # gRPC protocol definitions
├── cloud-backend/
│   ├── api-gateway/         # API Gateway (Gin + JWT)
│   │   ├── cmd/main.go
│   │   ├── internal/
│   │   │   ├── handlers/    # HTTP handlers
│   │   │   ├── middleware/  # JWT, CORS, logging
│   │   │   ├── alerter/     # Alert rule engine
│   │   │   ├── models/      # Shared data models
│   │   │   └── router/      # Route configuration
│   │   └── go.mod
│   └── shared/              # Shared code (db, etc.)
├── web-console/             # Web Console (React + Vite)
│   ├── src/
│   │   ├── pages/
│   │   │   ├── Dashboard.tsx    # Main dashboard with charts
│   │   │   ├── Servers.tsx
│   │   │   └── Cooling.tsx
│   │   ├── api/
│   │   └── types/
│   └── Dockerfile
├── docker-compose.yml
├── docker-compose.prod.yml
└── Makefile
```

## API Documentation

### Authentication

All API endpoints (except `/health` and `/auth/login`, `/auth/register`) require JWT authentication via the `Authorization: Bearer <token>` header.

### Endpoints

#### Health
- `GET /health` - Health check
- `GET /health/ready` - Readiness check

#### Authentication
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/logout` - Logout (requires JWT)
- `GET /api/v1/auth/me` - Get current user (requires JWT)
- `POST /api/v1/auth/refresh` - Refresh token (requires JWT)
- `POST /api/v1/auth/change-password` - Change password (requires JWT)

#### Assets (requires JWT)
- `GET /api/v1/assets/datacenters` - List data centers
- `GET /api/v1/assets/datacenters/:id` - Get data center
- `GET /api/v1/assets/servers` - List servers
- `GET /api/v1/assets/servers/:id` - Get server
- `GET /api/v1/assets/agents` - List edge agents
- `GET /api/v1/assets/cooling-devices` - List cooling devices

#### Telemetry (requires JWT)
- `GET /api/v1/telemetry/readings` - Get sensor readings
- `GET /api/v1/telemetry/pue` - Get PUE history
- `GET /api/v1/dashboard/stats` - Get dashboard statistics

#### Alerts (requires JWT)
- `GET /api/v1/alerts` - List alerts
- `GET /api/v1/alerts/:id` - Get alert
- `POST /api/v1/alerts/:id/acknowledge` - Acknowledge alert
- `POST /api/v1/alerts/:id/resolve` - Resolve alert

#### Alert Rules (requires admin role)
- `GET /api/v1/alert-rules` - List alert rules
- `POST /api/v1/alert-rules` - Create alert rule
- `GET /api/v1/alert-rules/:id` - Get alert rule
- `PUT /api/v1/alert-rules/:id` - Update alert rule
- `DELETE /api/v1/alert-rules/:id` - Delete alert rule

## Alert Rule Engine

The alert rule engine evaluates rules against incoming metrics every 30 seconds.

### Rule Configuration

```json
{
  "name": "High Supply Temperature",
  "description": "Trigger when CDU supply temperature exceeds 25°C for 5 minutes",
  "enabled": true,
  "severity": "warning",
  "category": "cooling",
  "conditions": [
    {
      "metric": "supply_temperature",
      "operator": ">",
      "threshold": 25.0,
      "duration": 300
    }
  ],
  "actions": [
    { "type": "email", "config": { "recipients": ["ops@example.com"] } },
    { "type": "slack", "config": { "channel": "#alerts" } }
  ]
}
```

### Supported Operators
- `>` - Greater than
- `>=` - Greater than or equal
- `<` - Less than
- `<=` - Less than or equal
- `==` - Equal
- `!=` - Not equal

### Actions
- `email` - Send email notification
- `webhook` - Call HTTP webhook
- `slack` - Send Slack message

## Dashboard Charts

The dashboard provides real-time visualization of:

1. **PUE Trend** - Power Usage Effectiveness over time
2. **Power Consumption** - IT power vs Cooling power
3. **Cooling Temperature** - Supply, return, and ambient temperature
4. **Server Utilization** - CPU and memory usage
5. **Alert History** - Alert counts by severity (last 7 days)
6. **Server Health** - Distribution of healthy/warning/critical/offline servers
7. **Cooling Efficiency** - CDU efficiency metrics

## Database Schema

### Core Tables
- `users` - User accounts with role-based access
- `data_centers` - Data center information
- `edge_agents` - Edge Agent management
- `servers` - Server assets
- `cooling_devices` - Liquid cooling devices (CDU, cooling towers)
- `sensor_readings` - Sensor time-series data (TimescaleDB hypertable)
- `alerts` - Alert records
- `pue_calculations` - PUE calculation history

## Configuration

### Edge Agent Example

```yaml
# edge-agent/config.yaml
agent_id: ""
hostname: "edge-dc1-01"
bootstrap_token: "your-token"
cloud_endpoint: "cloudmanager.example.com:443"
datacenter_id: "dc-1"
discovery:
  redfish_ranges:
    - "192.168.1.0/24"
  snmp_ranges:
    - "192.168.2.0/24"
  interval_sec: 3600
collection:
  interval_sec: 30
```

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | API Gateway port | 8080 |
| `DB_HOST` | PostgreSQL host | localhost |
| `DB_PORT` | PostgreSQL port | 5432 |
| `DB_USER` | Database user | cloudmanager |
| `DB_PASSWORD` | Database password | changeme |
| `DB_NAME` | Database name | cloudmanager |
| `JWT_SECRET` | JWT signing secret | default-secret-change-in-production |

## Development Commands

```bash
# View all available commands
make help

# Development environment
make dev-up        # Start infrastructure
make dev-down      # Stop infrastructure
make dev-logs      # View logs

# Testing
make test          # Run all tests
make test-edge     # Run Edge Agent tests

# Code generation
make proto         # Generate protobuf code

# Building
make build         # Build all services
make prod-build    # Build production images
```

## Contributing

1. Fork the project
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'feat: add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Create a Pull Request

## License

MIT License

## Acknowledgments

Inspired by [Supermicro SuperCloud Composer](https://www.supermicro.org.cn/en/solutions/management-software/supercloud-composer)
