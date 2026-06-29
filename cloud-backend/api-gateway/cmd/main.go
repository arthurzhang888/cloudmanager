package main

import (
	"log"
	"os"

	"github.com/cloudmanager/cloud-backend/api-gateway/internal/alerter"
	"github.com/cloudmanager/cloud-backend/api-gateway/internal/router"
	"github.com/cloudmanager/cloud-backend/shared/db"
)

func main() {
	// Get configuration from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "cloudmanager"
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "changeme"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "cloudmanager"
	}

	// Build database connection string
	dbURL := "postgres://" + dbUser + ":" + dbPassword + "@" + dbHost + ":" + dbPort + "/" + dbName + "?sslmode=disable"

	// Connect to database
	postgresDB, err := db.NewPostgresDB(dbURL)
	if err != nil {
		log.Printf("Warning: Failed to connect to database: %v", err)
		log.Println("Continuing without database connection...")
	} else {
		defer postgresDB.Close()
		log.Println("Connected to PostgreSQL database")
	}

	// Initialize alert rule engine
	metricsStore := alerter.NewMemoryMetricsStore()
	alertStore := alerter.NewMemoryAlertStore()
	alertEngine := alerter.NewEngine(metricsStore, alertStore)

	// Seed with sample rules
	seedAlertRules(alertEngine)

	// Start alert engine
	alertEngine.Start()
	defer alertEngine.Stop()

	// Setup router with alert engine
	r := router.SetupRouter(postgresDB, alertEngine, alertStore)

	// Start server
	log.Printf("Starting API Gateway on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// seedAlertRules adds sample alert rules for demonstration
func seedAlertRules(engine *alerter.Engine) {
	rules := []*alerter.Rule{
		{
			ID:          "rule-high-supply-temp",
			Name:        "High Supply Temperature",
			Description: "Trigger when CDU supply temperature exceeds 25°C for 5 minutes",
			Enabled:     true,
			Severity:    "warning",
			Category:    "cooling",
			Conditions: []alerter.Condition{
				{
					Metric:    "supply_temperature",
					Operator:  ">",
					Threshold: 25.0,
					Duration:  300, // 5 minutes
				},
			},
			Actions: []alerter.Action{
				{Type: "email", Config: map[string]interface{}{"recipients": []string{"ops@example.com"}}},
				{Type: "slack", Config: map[string]interface{}{"channel": "#alerts"}},
			},
		},
		{
			ID:          "rule-critical-supply-temp",
			Name:        "Critical Supply Temperature",
			Description: "Trigger when CDU supply temperature exceeds 30°C immediately",
			Enabled:     true,
			Severity:    "critical",
			Category:    "cooling",
			Conditions: []alerter.Condition{
				{
					Metric:    "supply_temperature",
					Operator:  ">",
					Threshold: 30.0,
					Duration:  0, // Immediate
				},
			},
			Actions: []alerter.Action{
				{Type: "email", Config: map[string]interface{}{"recipients": []string{"ops@example.com", "manager@example.com"}}},
				{Type: "webhook", Config: map[string]interface{}{"url": "https://pager.example.com/trigger"}},
			},
		},
		{
			ID:          "rule-high-cpu-temp",
			Name:        "High CPU Temperature",
			Description: "Trigger when server CPU temperature exceeds 80°C for 10 minutes",
			Enabled:     true,
			Severity:    "warning",
			Category:    "hardware",
			Conditions: []alerter.Condition{
				{
					Metric:    "cpu_temperature",
					Operator:  ">",
					Threshold: 80.0,
					Duration:  600, // 10 minutes
				},
			},
			Actions: []alerter.Action{
				{Type: "email", Config: map[string]interface{}{"recipients": []string{"ops@example.com"}}},
			},
		},
	}

	for _, rule := range rules {
		engine.AddRule(rule)
	}
	log.Printf("Seeded %d alert rules", len(rules))
}
