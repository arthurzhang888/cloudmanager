package handlers

import (
	"net/http"
	"time"

	"github.com/cloudmanager/cloud-backend/api-gateway/internal/alerter"
	"github.com/cloudmanager/cloud-backend/api-gateway/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AlertHandler handles alert-related requests
type AlertHandler struct {
	engine *alerter.Engine
	store  alerter.AlertStore
}

// NewAlertHandler creates a new alert handler
func NewAlertHandler(engine *alerter.Engine, store alerter.AlertStore) *AlertHandler {
	return &AlertHandler{
		engine: engine,
		store:  store,
	}
}

// Alert is an alias for models.Alert for backward compatibility
type Alert = models.Alert

// AlertRule is an alias for models.AlertRule for backward compatibility
type AlertRule = models.AlertRule

// AlertCondition is an alias for models.AlertCondition for backward compatibility
type AlertCondition = models.AlertCondition

// AlertAction is an alias for models.AlertAction for backward compatibility
type AlertAction = models.AlertAction

// ListAlerts returns a list of alerts
func (h *AlertHandler) ListAlerts(c *gin.Context) {
	// Parse query parameters
	status := c.Query("status")
	severity := c.Query("severity")
	category := c.Query("category")
	sourceType := c.Query("source_type")

	// Get alerts from store
	var alerts []*Alert
	if memStore, ok := h.store.(*alerter.MemoryAlertStore); ok {
		if status != "" {
			alerts = convertAlerts(memStore.GetAlertsByStatus(status))
		} else {
			alerts = convertAlerts(memStore.GetAllAlerts())
		}
	}

	// Filter by other parameters
	var filtered []*Alert
	for _, alert := range alerts {
		if severity != "" && alert.Severity != severity {
			continue
		}
		if category != "" && alert.Category != category {
			continue
		}
		if sourceType != "" && alert.SourceType != sourceType {
			continue
		}
		filtered = append(filtered, alert)
	}

	// If no alerts in store, return mock data for demo
	if len(filtered) == 0 {
		filtered = []*Alert{
			{
				ID:          uuid.New().String(),
				Severity:    "warning",
				Category:    "cooling",
				SourceType:  "cdu",
				SourceID:    "cdu-001",
				SourceName:  "CDU-Rack-A01",
				Title:       "Supply Temperature High",
				Description: "Supply temperature exceeded 25°C threshold",
				Status:      "active",
				CreatedAt:   time.Now().Add(-2 * time.Hour),
			},
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"alerts": filtered,
		"meta": gin.H{
			"total": len(filtered),
		},
	})
}

// convertAlerts converts alerter.Alert handlers to Alert handlers
func convertAlerts(alerts []*Alert) []*Alert {
	return alerts
}

// GetAlert returns a single alert by ID
func (h *AlertHandler) GetAlert(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid alert ID"})
		return
	}

	// TODO: Implement with actual database query
	c.JSON(http.StatusNotFound, gin.H{"error": "alert not found"})
}

// AcknowledgeAlert acknowledges an alert
func (h *AlertHandler) AcknowledgeAlert(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid alert ID"})
		return
	}

	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	_ = userID
	_ = username

	// TODO: Implement with actual database update
	c.JSON(http.StatusOK, gin.H{
		"message": "alert acknowledged",
		"alert_id": id,
	})
}

// ResolveAlert resolves an alert
func (h *AlertHandler) ResolveAlert(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid alert ID"})
		return
	}

	// TODO: Implement with actual database update
	c.JSON(http.StatusOK, gin.H{
		"message": "alert resolved",
		"alert_id": id,
	})
}

// ListAlertRules returns all alert rules
func (h *AlertHandler) ListAlertRules(c *gin.Context) {
	// Get rules from engine
	engineRules := h.engine.GetRules()

	// Convert to response format
	rules := make([]AlertRule, len(engineRules))
	for i, r := range engineRules {
		rules[i] = convertRuleToResponse(r)
	}

	c.JSON(http.StatusOK, gin.H{
		"rules": rules,
	})
}

// convertRuleToResponse converts an engine rule to response format
func convertRuleToResponse(r *alerter.Rule) AlertRule {
	conditions := make([]AlertCondition, len(r.Conditions))
	for i, c := range r.Conditions {
		conditions[i] = AlertCondition{
			Metric:    c.Metric,
			Operator:  c.Operator,
			Threshold: c.Threshold,
			Duration:  c.Duration,
		}
	}

	actions := make([]AlertAction, len(r.Actions))
	for i, a := range r.Actions {
		actions[i] = AlertAction{
			Type:   a.Type,
			Config: a.Config,
		}
	}

	return AlertRule{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Enabled:     r.Enabled,
		Severity:    r.Severity,
		Category:    r.Category,
		Conditions:  conditions,
		Actions:     actions,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

// GetAlertRule returns a single alert rule by ID
func (h *AlertHandler) GetAlertRule(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rule ID"})
		return
	}

	// TODO: Implement with actual database query
	c.JSON(http.StatusNotFound, gin.H{"error": "alert rule not found"})
}

// CreateAlertRule creates a new alert rule
func (h *AlertHandler) CreateAlertRule(c *gin.Context) {
	var req AlertRule
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate rule
	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "rule name is required"})
		return
	}
	if len(req.Conditions) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "at least one condition is required"})
		return
	}

	// Convert to engine rule
	conditions := make([]alerter.Condition, len(req.Conditions))
	for i, c := range req.Conditions {
		conditions[i] = alerter.Condition{
			Metric:    c.Metric,
			Operator:  c.Operator,
			Threshold: c.Threshold,
			Duration:  c.Duration,
		}
	}

	actions := make([]alerter.Action, len(req.Actions))
	for i, a := range req.Actions {
		actions[i] = alerter.Action{
			Type:   a.Type,
			Config: a.Config,
		}
	}

	rule := &alerter.Rule{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		Enabled:     req.Enabled,
		Severity:    req.Severity,
		Category:    req.Category,
		Conditions:  conditions,
		Actions:     actions,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Add to engine
	h.engine.AddRule(rule)

	c.JSON(http.StatusCreated, convertRuleToResponse(rule))
}

// UpdateAlertRule updates an existing alert rule
func (h *AlertHandler) UpdateAlertRule(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rule ID"})
		return
	}

	var req AlertRule
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert to engine rule
	conditions := make([]alerter.Condition, len(req.Conditions))
	for i, c := range req.Conditions {
		conditions[i] = alerter.Condition{
			Metric:    c.Metric,
			Operator:  c.Operator,
			Threshold: c.Threshold,
			Duration:  c.Duration,
		}
	}

	actions := make([]alerter.Action, len(req.Actions))
	for i, a := range req.Actions {
		actions[i] = alerter.Action{
			Type:   a.Type,
			Config: a.Config,
		}
	}

	rule := &alerter.Rule{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		Enabled:     req.Enabled,
		Severity:    req.Severity,
		Category:    req.Category,
		Conditions:  conditions,
		Actions:     actions,
		CreatedAt:   req.CreatedAt,
		UpdatedAt:   time.Now(),
	}

	// Update in engine
	h.engine.UpdateRule(rule)

	c.JSON(http.StatusOK, convertRuleToResponse(rule))
}

// DeleteAlertRule deletes an alert rule
func (h *AlertHandler) DeleteAlertRule(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rule ID"})
		return
	}

	// Remove from engine
	h.engine.RemoveRule(id)

	c.JSON(http.StatusOK, gin.H{"message": "alert rule deleted"})
}
