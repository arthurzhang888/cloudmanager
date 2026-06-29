package collector

import (
	"encoding/binary"
	"fmt"
	"math"
	"time"

	"github.com/goburrow/modbus"
)

// handlerWithConnectClose is the interface for handlers that support Connect/Close
type handlerWithConnectClose interface {
	Connect() error
	Close() error
}

// modbusClientAdapter wraps modbus.Client to add our interface methods
type modbusClientAdapter struct {
	client  modbus.Client
	handler handlerWithConnectClose
}

func (a *modbusClientAdapter) Connect() error {
	return a.handler.Connect()
}

func (a *modbusClientAdapter) Close() error {
	return a.handler.Close()
}

func (a *modbusClientAdapter) ReadHoldingRegisters(address, quantity uint16) ([]byte, error) {
	return a.client.ReadHoldingRegisters(address, quantity)
}

func (a *modbusClientAdapter) ReadInputRegisters(address, quantity uint16) ([]byte, error) {
	return a.client.ReadInputRegisters(address, quantity)
}

// ModbusCollector handles Modbus communication for cooling devices
type ModbusCollector struct {
	client  *modbusClientAdapter
	slaveID byte
	address string
}

// NewModbusCollector creates a new Modbus collector for TCP connection
func NewModbusCollector(address string, slaveID byte) (*ModbusCollector, error) {
	handler := modbus.NewTCPClientHandler(address)
	handler.SlaveId = slaveID
	handler.Timeout = 10 * time.Second

	return &ModbusCollector{
		client: &modbusClientAdapter{
			client:  modbus.NewClient(handler),
			handler: handler,
		},
		slaveID: slaveID,
		address: address,
	}, nil
}

// NewModbusCollectorRTU creates a new Modbus collector for RTU connection
func NewModbusCollectorRTU(device string, baudRate int, slaveID byte) (*ModbusCollector, error) {
	handler := modbus.NewRTUClientHandler(device)
	handler.BaudRate = baudRate
	handler.DataBits = 8
	handler.Parity = "N"
	handler.StopBits = 1
	handler.SlaveId = slaveID
	handler.Timeout = 10 * time.Second

	return &ModbusCollector{
		client: &modbusClientAdapter{
			client:  modbus.NewClient(handler),
			handler: handler,
		},
		slaveID: slaveID,
		address: device,
	}, nil
}

// Connect establishes the Modbus connection
func (c *ModbusCollector) Connect() error {
	return c.client.Connect()
}

// Close closes the Modbus connection
func (c *ModbusCollector) Close() error {
	return c.client.Close()
}

// ReadUint16 reads a single uint16 register
func (c *ModbusCollector) ReadUint16(address uint16) (uint16, error) {
	data, err := c.client.ReadInputRegisters(address, 1)
	if err != nil {
		return 0, fmt.Errorf("failed to read register %d: %w", address, err)
	}
	return binary.BigEndian.Uint16(data), nil
}

// ReadInt16 reads a single int16 register
func (c *ModbusCollector) ReadInt16(address uint16) (int16, error) {
	data, err := c.client.ReadInputRegisters(address, 1)
	if err != nil {
		return 0, fmt.Errorf("failed to read register %d: %w", address, err)
	}
	return int16(binary.BigEndian.Uint16(data)), nil
}

// ReadFloat32 reads two registers as a float32 (IEEE 754)
func (c *ModbusCollector) ReadFloat32(address uint16) (float32, error) {
	data, err := c.client.ReadInputRegisters(address, 2)
	if err != nil {
		return 0, fmt.Errorf("failed to read registers %d-%d: %w", address, address+1, err)
	}
	bits := binary.BigEndian.Uint32(data)
	return math.Float32frombits(bits), nil
}

// ReadUint32 reads two registers as uint32
func (c *ModbusCollector) ReadUint32(address uint16) (uint32, error) {
	data, err := c.client.ReadInputRegisters(address, 2)
	if err != nil {
		return 0, fmt.Errorf("failed to read registers %d-%d: %w", address, address+1, err)
	}
	return binary.BigEndian.Uint32(data), nil
}

// RegisterDef defines a register mapping for a metric
type RegisterDef struct {
	Address  uint16  `json:"address"`
	DataType string  `json:"data_type"` // "uint16", "int16", "uint32", "float32"
	Scale    float64 `json:"scale"`
	Unit     string  `json:"unit"`
}

// CoolingDeviceConfig holds configuration for a cooling device
type CoolingDeviceConfig struct {
	Name        string                 `json:"name"`
	DeviceType  string                 `json:"device_type"` // "cdu" or "cooling_tower"
	RegisterMap map[string]RegisterDef `json:"register_map"`
}

// MetricReading represents a single sensor reading
type MetricReading struct {
	MetricName string    `json:"metric_name"`
	Value      float64   `json:"value"`
	Unit       string    `json:"unit"`
	Timestamp  time.Time `json:"timestamp"`
}

// CoolingDeviceReader handles reading metrics from a cooling device
type CoolingDeviceReader struct {
	collector *ModbusCollector
	config    CoolingDeviceConfig
}

// NewCoolingDeviceReader creates a new reader for a cooling device
func NewCoolingDeviceReader(collector *ModbusCollector, config CoolingDeviceConfig) *CoolingDeviceReader {
	return &CoolingDeviceReader{
		collector: collector,
		config:    config,
	}
}

// ReadAll reads all configured metrics from the device
func (r *CoolingDeviceReader) ReadAll() ([]MetricReading, error) {
	readings := make([]MetricReading, 0, len(r.config.RegisterMap))
	timestamp := time.Now()

	for metricName, regDef := range r.config.RegisterMap {
		value, err := r.readRegister(regDef)
		if err != nil {
			// Log error but continue reading other metrics
			continue
		}

		readings = append(readings, MetricReading{
			MetricName: metricName,
			Value:      value,
			Unit:       regDef.Unit,
			Timestamp:  timestamp,
		})
	}

	return readings, nil
}

// ReadMetric reads a single metric
func (r *CoolingDeviceReader) ReadMetric(metricName string) (MetricReading, error) {
	regDef, exists := r.config.RegisterMap[metricName]
	if !exists {
		return MetricReading{}, fmt.Errorf("metric %s not found in register map", metricName)
	}

	value, err := r.readRegister(regDef)
	if err != nil {
		return MetricReading{}, err
	}

	return MetricReading{
		MetricName: metricName,
		Value:      value,
		Unit:       regDef.Unit,
		Timestamp:  time.Now(),
	}, nil
}

func (r *CoolingDeviceReader) readRegister(def RegisterDef) (float64, error) {
	switch def.DataType {
	case "uint16":
		val, err := r.collector.ReadUint16(def.Address)
		if err != nil {
			return 0, err
		}
		return float64(val) * def.Scale, nil

	case "int16":
		val, err := r.collector.ReadInt16(def.Address)
		if err != nil {
			return 0, err
		}
		return float64(val) * def.Scale, nil

	case "uint32":
		val, err := r.collector.ReadUint32(def.Address)
		if err != nil {
			return 0, err
		}
		return float64(val) * def.Scale, nil

	case "float32":
		val, err := r.collector.ReadFloat32(def.Address)
		if err != nil {
			return 0, err
		}
		return float64(val) * def.Scale, nil

	default:
		return 0, fmt.Errorf("unsupported data type: %s", def.DataType)
	}
}

// Common CDU register mappings
var DefaultCDURegisterMap = map[string]RegisterDef{
	"supply_temperature": {
		Address:  0,
		DataType: "uint16",
		Scale:    0.1,
		Unit:     "°C",
	},
	"return_temperature": {
		Address:  1,
		DataType: "uint16",
		Scale:    0.1,
		Unit:     "°C",
	},
	"flow_rate": {
		Address:  2,
		DataType: "uint16",
		Scale:    1.0,
		Unit:     "L/min",
	},
	"pressure": {
		Address:  3,
		DataType: "uint16",
		Scale:    0.01,
		Unit:     "bar",
	},
	"pump_speed": {
		Address:  4,
		DataType: "uint16",
		Scale:    1.0,
		Unit:     "RPM",
	},
	"power": {
		Address:  10,
		DataType: "uint32",
		Scale:    1.0,
		Unit:     "W",
	},
}

// Common Cooling Tower register mappings
var DefaultCoolingTowerRegisterMap = map[string]RegisterDef{
	"ambient_temperature": {
		Address:  0,
		DataType: "uint16",
		Scale:    0.1,
		Unit:     "°C",
	},
	"wet_bulb_temperature": {
		Address:  1,
		DataType: "uint16",
		Scale:    0.1,
		Unit:     "°C",
	},
	"fan_speed": {
		Address:  2,
		DataType: "uint16",
		Scale:    1.0,
		Unit:     "RPM",
	},
	"water_flow": {
		Address:  3,
		DataType: "uint16",
		Scale:    1.0,
		Unit:     "m³/h",
	},
	"power": {
		Address:  10,
		DataType: "uint32",
		Scale:    1.0,
		Unit:     "W",
	},
}
