package models

import "time"

// Alert represents an alert in the system
type Alert struct {
	ID             string     `json:"id"`
	Severity       string     `json:"severity"` // critical, warning, info
	Category       string     `json:"category"` // hardware, cooling, network, threshold
	SourceType     string     `json:"source_type"`
	SourceID       string     `json:"source_id"`
	SourceName     string     `json:"source_name"`
	Title          string     `json:"title"`
	Description    string     `json:"description"`
	Status         string     `json:"status"` // active, acknowledged, resolved
	AcknowledgedBy string     `json:"acknowledged_by,omitempty"`
	AcknowledgedAt *time.Time `json:"acknowledged_at,omitempty"`
	ResolvedAt     *time.Time `json:"resolved_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}

// AlertRule represents an alert rule
type AlertRule struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Enabled     bool                   `json:"enabled"`
	Severity    string                 `json:"severity"`
	Category    string                 `json:"category"`
	Conditions  []AlertCondition       `json:"conditions"`
	Actions     []AlertAction          `json:"actions"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// AlertCondition represents a condition in an alert rule
type AlertCondition struct {
	Metric    string  `json:"metric"`
	Operator  string  `json:"operator"` // >, <, ==, >=, <=, !=
	Threshold float64 `json:"threshold"`
	Duration  int     `json:"duration"` // seconds
}

// AlertAction represents an action to take when alert triggers
type AlertAction struct {
	Type   string                 `json:"type"` // email, webhook, slack
	Config map[string]interface{} `json:"config"`
}
