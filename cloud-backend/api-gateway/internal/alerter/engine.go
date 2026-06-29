package alerter

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/cloudmanager/cloud-backend/api-gateway/internal/models"
	"github.com/google/uuid"
)

// MetricValue represents a single metric reading
type MetricValue struct {
	Timestamp time.Time
	Value     float64
	Tags      map[string]string
}

// MetricsStore interface for querying metrics
type MetricsStore interface {
	GetMetric(ctx context.Context, metricName string, sourceType, sourceID string, from, to time.Time) ([]MetricValue, error)
	GetLatestMetric(ctx context.Context, metricName string, sourceType, sourceID string) (MetricValue, error)
}

// AlertStore interface for storing alerts
type AlertStore interface {
	CreateAlert(alert *models.Alert) error
	GetActiveAlert(ruleID, sourceType, sourceID string) (*models.Alert, error)
	UpdateAlert(alert *models.Alert) error
}

// Rule represents an alert rule for evaluation
type Rule struct {
	ID          string
	Name        string
	Description string
	Enabled     bool
	Severity    string
	Category    string
	Conditions  []Condition
	Actions     []Action
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Condition represents a single condition in a rule
type Condition struct {
	Metric    string
	Operator  string
	Threshold float64
	Duration  int
}

// Action represents an action to take when rule triggers
type Action struct {
	Type   string
	Config map[string]interface{}
}

// Engine evaluates alert rules against metrics
type Engine struct {
	rules       map[string]*Rule
	rulesMu     sync.RWMutex
	metrics     MetricsStore
	alerts      AlertStore
	evalInterval time.Duration
	stopCh      chan struct{}
	wg          sync.WaitGroup

	// Track condition state for duration-based rules
	conditionState map[string]*conditionTracker
	stateMu        sync.Mutex
}

// conditionTracker tracks how long a condition has been true
type conditionTracker struct {
	RuleID      string
	SourceType  string
	SourceID    string
	ConditionID string
	FirstSeen   time.Time
	LastValue   float64
}

// NewEngine creates a new alert rule engine
func NewEngine(metrics MetricsStore, alerts AlertStore) *Engine {
	return &Engine{
		rules:          make(map[string]*Rule),
		metrics:        metrics,
		alerts:         alerts,
		evalInterval:   30 * time.Second,
		stopCh:         make(chan struct{}),
		conditionState: make(map[string]*conditionTracker),
	}
}

// Start begins the rule evaluation loop
func (e *Engine) Start() {
	e.wg.Add(1)
	go e.evaluationLoop()
	log.Println("Alert rule engine started")
}

// Stop stops the rule evaluation loop
func (e *Engine) Stop() {
	close(e.stopCh)
	e.wg.Wait()
	log.Println("Alert rule engine stopped")
}

// AddRule adds a rule to the engine
func (e *Engine) AddRule(rule *Rule) {
	e.rulesMu.Lock()
	defer e.rulesMu.Unlock()
	e.rules[rule.ID] = rule
	log.Printf("Added alert rule: %s (%s)", rule.Name, rule.ID)
}

// RemoveRule removes a rule from the engine
func (e *Engine) RemoveRule(ruleID string) {
	e.rulesMu.Lock()
	defer e.rulesMu.Unlock()
	delete(e.rules, ruleID)
	log.Printf("Removed alert rule: %s", ruleID)
}

// UpdateRule updates an existing rule
func (e *Engine) UpdateRule(rule *Rule) {
	e.rulesMu.Lock()
	defer e.rulesMu.Unlock()
	e.rules[rule.ID] = rule
	log.Printf("Updated alert rule: %s (%s)", rule.Name, rule.ID)
}

// GetRules returns all rules
func (e *Engine) GetRules() []*Rule {
	e.rulesMu.RLock()
	defer e.rulesMu.RUnlock()

	rules := make([]*Rule, 0, len(e.rules))
	for _, r := range e.rules {
		rules = append(rules, r)
	}
	return rules
}

// evaluationLoop runs the periodic rule evaluation
func (e *Engine) evaluationLoop() {
	defer e.wg.Done()

	ticker := time.NewTicker(e.evalInterval)
	defer ticker.Stop()

	// Run initial evaluation
	e.evaluateAllRules()

	for {
		select {
		case <-ticker.C:
			e.evaluateAllRules()
		case <-e.stopCh:
			return
		}
	}
}

// evaluateAllRules evaluates all enabled rules
func (e *Engine) evaluateAllRules() {
	e.rulesMu.RLock()
	rules := make([]*Rule, 0, len(e.rules))
	for _, r := range e.rules {
		if r.Enabled {
			rules = append(rules, r)
		}
	}
	e.rulesMu.RUnlock()

	for _, rule := range rules {
		e.evaluateRule(rule)
	}
}

// evaluateRule evaluates a single rule against current metrics
func (e *Engine) evaluateRule(rule *Rule) {
	ctx := context.Background()

	// For each condition in the rule, check if it's violated
	for i, condition := range rule.Conditions {
		// Get all sources that have this metric
		sources := e.getSourcesForMetric(ctx, condition.Metric)

		for _, source := range sources {
			metric, err := e.metrics.GetLatestMetric(ctx, condition.Metric, source.Type, source.ID)
			if err != nil {
				continue
			}

			conditionMet := e.evaluateCondition(condition, metric.Value)
			conditionKey := fmt.Sprintf("%s:%s:%s:%d", rule.ID, source.Type, source.ID, i)

			if conditionMet {
				e.handleConditionMet(rule, condition, conditionKey, source.Type, source.ID, metric.Value)
			} else {
				e.clearConditionState(conditionKey)
			}
		}
	}
}

// Source identifies a metric source
type Source struct {
	Type string
	ID   string
}

// getSourcesForMetric returns all sources that have a given metric
// This is a simplified version - in production, this would query the metrics store
func (e *Engine) getSourcesForMetric(ctx context.Context, metric string) []Source {
	// Mock sources for common metrics
	switch metric {
	case "supply_temperature", "return_temperature", "flow_rate", "pressure":
		return []Source{
			{Type: "cdu", ID: "cdu-001"},
			{Type: "cdu", ID: "cdu-002"},
		}
	case "cpu_temperature", "power_consumption":
		return []Source{
			{Type: "server", ID: "server-001"},
			{Type: "server", ID: "server-002"},
		}
	case "pue":
		return []Source{
			{Type: "datacenter", ID: "dc-001"},
		}
	}
	return nil
}

// evaluateCondition checks if a condition is met
func (e *Engine) evaluateCondition(condition Condition, value float64) bool {
	switch condition.Operator {
	case ">":
		return value > condition.Threshold
	case ">=":
		return value >= condition.Threshold
	case "<":
		return value < condition.Threshold
	case "<=":
		return value <= condition.Threshold
	case "==":
		return value == condition.Threshold
	case "!=":
		return value != condition.Threshold
	default:
		return false
	}
}

// handleConditionMet handles when a condition is met
func (e *Engine) handleConditionMet(rule *Rule, condition Condition, conditionKey, sourceType, sourceID string, value float64) {
	e.stateMu.Lock()
	defer e.stateMu.Unlock()

	tracker, exists := e.conditionState[conditionKey]
	if !exists {
		// First time condition is met
		tracker = &conditionTracker{
			RuleID:      rule.ID,
			SourceType:  sourceType,
			SourceID:    sourceID,
			ConditionID: conditionKey,
			FirstSeen:   time.Now(),
			LastValue:   value,
		}
		e.conditionState[conditionKey] = tracker
		return
	}

	tracker.LastValue = value

	// Check if condition has been met for the required duration
	duration := time.Duration(condition.Duration) * time.Second
	if time.Since(tracker.FirstSeen) >= duration {
		// Trigger alert
		e.triggerAlert(rule, condition, sourceType, sourceID, value)
		// Reset tracker to avoid duplicate alerts
		delete(e.conditionState, conditionKey)
	}
}

// clearConditionState clears the condition state when condition is no longer met
func (e *Engine) clearConditionState(conditionKey string) {
	e.stateMu.Lock()
	defer e.stateMu.Unlock()
	delete(e.conditionState, conditionKey)
}

// triggerAlert creates an alert
func (e *Engine) triggerAlert(rule *Rule, condition Condition, sourceType, sourceID string, value float64) {
	// Check if there's already an active alert for this rule and source
	existingAlert, err := e.alerts.GetActiveAlert(rule.ID, sourceType, sourceID)
	if err == nil && existingAlert != nil {
		// Update existing alert with new information
		existingAlert.Description = fmt.Sprintf("%s: current value %.2f, threshold %.2f",
			condition.Metric, value, condition.Threshold)
		e.alerts.UpdateAlert(existingAlert)
		return
	}

	// Create new alert
	alert := &models.Alert{
		ID:           uuid.New().String(),
		Severity:     rule.Severity,
		Category:     rule.Category,
		SourceType:   sourceType,
		SourceID:     sourceID,
		SourceName:   fmt.Sprintf("%s-%s", sourceType, sourceID),
		Title:        rule.Name,
		Description:  fmt.Sprintf("%s: current value %.2f, threshold %.2f", condition.Metric, value, condition.Threshold),
		Status:       "active",
		CreatedAt:    time.Now(),
	}

	if err := e.alerts.CreateAlert(alert); err != nil {
		log.Printf("Failed to create alert: %v", err)
		return
	}

	log.Printf("Alert triggered: %s for %s/%s (value: %.2f)",
		rule.Name, sourceType, sourceID, value)

	// Execute actions
	for _, action := range rule.Actions {
		e.executeAction(action, alert)
	}
}

// executeAction executes an alert action
func (e *Engine) executeAction(action Action, alert *models.Alert) {
	switch action.Type {
	case "email":
		log.Printf("[Action] Sending email for alert %s", alert.ID)
	case "webhook":
		log.Printf("[Action] Calling webhook for alert %s", alert.ID)
	case "slack":
		log.Printf("[Action] Sending Slack notification for alert %s", alert.ID)
	default:
		log.Printf("[Action] Unknown action type: %s", action.Type)
	}
}
