package alerter

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/cloudmanager/cloud-backend/api-gateway/internal/models"
)

// MemoryMetricsStore is an in-memory implementation of MetricsStore
type MemoryMetricsStore struct {
	mu      sync.RWMutex
	metrics map[string][]MetricValue
}

// NewMemoryMetricsStore creates a new in-memory metrics store
func NewMemoryMetricsStore() *MemoryMetricsStore {
	store := &MemoryMetricsStore{
		metrics: make(map[string][]MetricValue),
	}
	// Seed with some mock data
	store.seedMockData()
	return store
}

// GetMetric returns metrics for a given time range
func (s *MemoryMetricsStore) GetMetric(ctx context.Context, metricName string, sourceType, sourceID string, from, to time.Time) ([]MetricValue, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	key := s.metricKey(metricName, sourceType, sourceID)
	values, exists := s.metrics[key]
	if !exists {
		return nil, errors.New("metric not found")
	}

	var result []MetricValue
	for _, v := range values {
		if v.Timestamp.After(from) && v.Timestamp.Before(to) {
			result = append(result, v)
		}
	}
	return result, nil
}

// GetLatestMetric returns the most recent metric value
func (s *MemoryMetricsStore) GetLatestMetric(ctx context.Context, metricName string, sourceType, sourceID string) (MetricValue, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	key := s.metricKey(metricName, sourceType, sourceID)
	values, exists := s.metrics[key]
	if !exists || len(values) == 0 {
		return MetricValue{}, errors.New("metric not found")
	}

	return values[len(values)-1], nil
}

// StoreMetric stores a new metric value
func (s *MemoryMetricsStore) StoreMetric(metricName, sourceType, sourceID string, value float64, tags map[string]string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := s.metricKey(metricName, sourceType, sourceID)
	s.metrics[key] = append(s.metrics[key], MetricValue{
		Timestamp: time.Now(),
		Value:     value,
		Tags:      tags,
	})
}

func (s *MemoryMetricsStore) metricKey(metricName, sourceType, sourceID string) string {
	return metricName + ":" + sourceType + ":" + sourceID
}

// seedMockData adds mock data for testing
func (s *MemoryMetricsStore) seedMockData() {
	now := time.Now()

	// CDU metrics
	for i := 0; i < 60; i++ {
		s.metrics["supply_temperature:cdu:cdu-001"] = append(
			s.metrics["supply_temperature:cdu:cdu-001"],
			MetricValue{
				Timestamp: now.Add(-time.Duration(60-i) * time.Minute),
				Value:     22.0 + float64(i)*0.1,
				Tags:      map[string]string{"unit": "celsius"},
			},
		)
	}

	// High temperature for testing alerts
	s.metrics["supply_temperature:cdu:cdu-002"] = []MetricValue{
		{Timestamp: now.Add(-5 * time.Minute), Value: 26.5, Tags: map[string]string{"unit": "celsius"}},
		{Timestamp: now.Add(-4 * time.Minute), Value: 27.0, Tags: map[string]string{"unit": "celsius"}},
		{Timestamp: now.Add(-3 * time.Minute), Value: 27.5, Tags: map[string]string{"unit": "celsius"}},
		{Timestamp: now.Add(-2 * time.Minute), Value: 28.0, Tags: map[string]string{"unit": "celsius"}},
		{Timestamp: now.Add(-1 * time.Minute), Value: 28.5, Tags: map[string]string{"unit": "celsius"}},
		{Timestamp: now, Value: 29.0, Tags: map[string]string{"unit": "celsius"}},
	}

	// Server metrics
	s.metrics["cpu_temperature:server:server-001"] = []MetricValue{
		{Timestamp: now, Value: 65.0, Tags: map[string]string{"unit": "celsius"}},
	}
}

// MemoryAlertStore is an in-memory implementation of AlertStore
type MemoryAlertStore struct {
	mu     sync.RWMutex
	alerts map[string]*models.Alert
}

// NewMemoryAlertStore creates a new in-memory alert store
func NewMemoryAlertStore() *MemoryAlertStore {
	return &MemoryAlertStore{
		alerts: make(map[string]*models.Alert),
	}
}

// CreateAlert creates a new alert
func (s *MemoryAlertStore) CreateAlert(alert *models.Alert) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.alerts[alert.ID] = alert
	return nil
}

// GetActiveAlert checks if there's an active alert for a rule and source
func (s *MemoryAlertStore) GetActiveAlert(ruleID, sourceType, sourceID string) (*models.Alert, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Note: In a real implementation, we'd need to track which rule created each alert
	// For now, we'll just check if there's any active alert for this source
	for _, alert := range s.alerts {
		if alert.SourceType == sourceType && alert.SourceID == sourceID && alert.Status == "active" {
			return alert, nil
		}
	}
	return nil, errors.New("no active alert found")
}

// UpdateAlert updates an existing alert
func (s *MemoryAlertStore) UpdateAlert(alert *models.Alert) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.alerts[alert.ID] = alert
	return nil
}

// GetAllAlerts returns all alerts
func (s *MemoryAlertStore) GetAllAlerts() []*models.Alert {
	s.mu.RLock()
	defer s.mu.RUnlock()

	alerts := make([]*models.Alert, 0, len(s.alerts))
	for _, alert := range s.alerts {
		alerts = append(alerts, alert)
	}
	return alerts
}

// GetAlertsByStatus returns alerts filtered by status
func (s *MemoryAlertStore) GetAlertsByStatus(status string) []*models.Alert {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var alerts []*models.Alert
	for _, alert := range s.alerts {
		if alert.Status == status {
			alerts = append(alerts, alert)
		}
	}
	return alerts
}
