package collector

import (
	"math"
	"testing"
)

func TestRegisterDefScale(t *testing.T) {
	tests := []struct {
		name     string
		regDef   RegisterDef
		rawValue uint16
		expected float64
	}{
		{
			name: "temperature with 0.1 scale",
			regDef: RegisterDef{
				Address:  0,
				DataType: "uint16",
				Scale:    0.1,
				Unit:     "°C",
			},
			rawValue: 250,
			expected: 25.0,
		},
		{
			name: "pressure with 0.01 scale",
			regDef: RegisterDef{
				Address:  3,
				DataType: "uint16",
				Scale:    0.01,
				Unit:     "bar",
			},
			rawValue: 230,
			expected: 2.3,
		},
		{
			name: "flow rate with 1.0 scale",
			regDef: RegisterDef{
				Address:  2,
				DataType: "uint16",
				Scale:    1.0,
				Unit:     "L/min",
			},
			rawValue: 100,
			expected: 100.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := float64(tt.rawValue) * tt.regDef.Scale
			if math.Abs(result-tt.expected) > 0.0001 {
				t.Errorf("expected %f, got %f", tt.expected, result)
			}
		})
	}
}

func TestCoolingDeviceConfig(t *testing.T) {
	config := CoolingDeviceConfig{
		Name:       "CDU-Test",
		DeviceType: "cdu",
		RegisterMap: map[string]RegisterDef{
			"supply_temp": {
				Address:  0,
				DataType: "uint16",
				Scale:    0.1,
				Unit:     "°C",
			},
		},
	}

	if config.Name != "CDU-Test" {
		t.Errorf("expected name CDU-Test, got %s", config.Name)
	}

	if config.DeviceType != "cdu" {
		t.Errorf("expected device type cdu, got %s", config.DeviceType)
	}

	reg, exists := config.RegisterMap["supply_temp"]
	if !exists {
		t.Error("expected supply_temp to exist in register map")
	}

	if reg.Address != 0 {
		t.Errorf("expected address 0, got %d", reg.Address)
	}

	if reg.Scale != 0.1 {
		t.Errorf("expected scale 0.1, got %f", reg.Scale)
	}
}

func TestDefaultCDURegisterMap(t *testing.T) {
	// Verify default CDU register map has expected metrics
	expectedMetrics := []string{
		"supply_temperature",
		"return_temperature",
		"flow_rate",
		"pressure",
		"pump_speed",
		"power",
	}

	for _, metric := range expectedMetrics {
		if _, exists := DefaultCDURegisterMap[metric]; !exists {
			t.Errorf("expected metric %s in DefaultCDURegisterMap", metric)
		}
	}
}

func TestDefaultCoolingTowerRegisterMap(t *testing.T) {
	// Verify default cooling tower register map has expected metrics
	expectedMetrics := []string{
		"ambient_temperature",
		"wet_bulb_temperature",
		"fan_speed",
		"water_flow",
		"power",
	}

	for _, metric := range expectedMetrics {
		if _, exists := DefaultCoolingTowerRegisterMap[metric]; !exists {
			t.Errorf("expected metric %s in DefaultCoolingTowerRegisterMap", metric)
		}
	}
}

func TestMetricReading(t *testing.T) {
	reading := MetricReading{
		MetricName: "supply_temperature",
		Value:      25.5,
		Unit:       "°C",
	}

	if reading.MetricName != "supply_temperature" {
		t.Errorf("expected metric name supply_temperature, got %s", reading.MetricName)
	}

	if reading.Value != 25.5 {
		t.Errorf("expected value 25.5, got %f", reading.Value)
	}

	if reading.Unit != "°C" {
		t.Errorf("expected unit °C, got %s", reading.Unit)
	}
}
