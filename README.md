# CloudManager

企业级数据中心基础设施管理平台，支持多地大型数据中心（1000+ 服务器）的统一管理，具备可组合基础设施能力和液冷系统深度集成。

## 系统架构

```
┌─────────────────────────────────────────────────────────┐
│                    云端管控平台                           │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐   │
│  │ API Gateway│ │ Asset    │ │ Telemetry│ │ Web      │   │
│  │           │ │ Service  │ │ Service  │ │ Console  │   │
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
    │           数据中心硬件               │
    │  ┌──────┐ ┌──────┐ ┌──────┐        │
    │  │服务器 │ │交换机 │ │ CDU  │        │
    │  │Redfish│ │SNMP  │ │Modbus│        │
    │  └──────┘ └──────┘ └──────┘        │
    └─────────────────────────────────────┘
```

## 技术栈

- **边缘 Agent**: Go, gRPC, Modbus, Redfish, SNMP
- **云端后端**: Go, Gin, gRPC, PostgreSQL, TimescaleDB, Redis, NATS
- **前端**: React, TypeScript, Vite, Tailwind CSS, Recharts
- **部署**: Docker, Docker Compose

## 核心功能

### MVP 已实现
- [x] 项目架构和开发环境
- [x] 数据库设计（资产、液冷设备、时序数据、告警）
- [x] gRPC 通信协议定义
- [x] Edge Agent 配置模块
- [x] Web 控制台（Dashboard、Servers、Cooling、Topology）

### 进行中
- [ ] Modbus 采集模块
- [ ] Asset Service 数据库层
- [ ] 云端 API 服务
- [ ] 告警规则引擎

### 计划中
- [ ] 资源池化和编排
- [ ] 历史数据分析和预测
- [ ] 多租户支持

## 快速开始

### 环境要求
- Docker & Docker Compose
- Node.js 18+ (前端开发)
- Go 1.22+ (后端开发)

### 启动开发环境

```bash
# 1. 克隆项目
git clone <repository-url>
cd cloudmanager

# 2. 启动基础设施（数据库、消息队列等）
make dev-up

# 3. 运行数据库迁移
make migrate

# 4. 启动前端开发服务器
make web-dev
```

访问 http://localhost:3000 查看 Web 控制台。

### 生产部署

```bash
# 构建并启动生产环境
make prod-build
make prod-up
```

## 项目结构

```
cloudmanager/
├── edge-agent/              # 边缘数据采集 Agent
│   ├── cmd/agent/           # 主程序入口
│   ├── internal/            # 内部模块
│   │   ├── config/          # 配置管理
│   │   ├── discovery/       # 设备发现（Redfish/SNMP）
│   │   ├── collector/       # 数据采集（Modbus）
│   │   └── uploader/        # 数据上报
│   └── proto/               # gRPC 协议定义
├── cloud-backend/           # 云端后端服务
│   ├── api-gateway/         # API 网关
│   ├── asset-service/       # 资产管理服务
│   ├── telemetry-service/   # 遥测数据处理
│   └── shared/              # 共享代码
├── web-console/             # Web 控制台（React）
│   ├── src/
│   │   ├── components/      # UI 组件
│   │   ├── pages/           # 页面
│   │   ├── api/             # API 客户端
│   │   └── types/           # TypeScript 类型
│   └── Dockerfile
├── migrations/              # 数据库迁移脚本
├── docker-compose.yml       # 开发环境配置
├── docker-compose.prod.yml  # 生产环境配置
└── Makefile                 # 常用命令
```

## 数据库 Schema

### 核心表
- `data_centers` - 数据中心信息
- `edge_agents` - 边缘 Agent 管理
- `servers` - 服务器资产
- `cooling_devices` - 液冷设备（CDU、冷却塔）
- `sensor_readings` - 传感器时序数据（TimescaleDB hypertable）
- `alerts` - 告警记录
- `pue_calculations` - PUE 计算历史

## API 文档

详见 [docs/api.md](docs/api.md)

## 配置说明

### Edge Agent 配置示例

```yaml
# edge-agent/config.yaml
agent_id: ""  # 首次启动为空，bootstrap 后自动填充
hostname: "edge-dc1-01"
bootstrap_token: "your-bootstrap-token"
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

## 开发指南

### 常用命令

```bash
# 查看所有可用命令
make help

# 开发环境
make dev-up        # 启动基础设施
make dev-down      # 停止基础设施
make dev-logs      # 查看日志

# 测试
make test          # 运行所有测试
make test-edge     # 运行 Edge Agent 测试

# 代码生成
make proto         # 生成 protobuf 代码

# 构建
make build         # 构建所有服务
make prod-build    # 构建生产镜像
```

### 提交代码

```bash
git add .
git commit -m "feat: your changes"
```

## 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'feat: add amazing feature'`)
4. 推送分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

## 许可证

MIT License

## 联系方式

如有问题或建议，请提交 Issue 或联系维护者。
