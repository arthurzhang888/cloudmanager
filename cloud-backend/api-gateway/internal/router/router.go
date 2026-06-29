package router

import (
	"os"

	"github.com/cloudmanager/cloud-backend/api-gateway/internal/alerter"
	"github.com/cloudmanager/cloud-backend/api-gateway/internal/handlers"
	"github.com/cloudmanager/cloud-backend/api-gateway/internal/middleware"
	"github.com/cloudmanager/cloud-backend/api-gateway/internal/repository"
	"github.com/cloudmanager/cloud-backend/shared/db"
	"github.com/gin-gonic/gin"
)

// SetupRouter configures all routes for the API gateway
func SetupRouter(postgresDB *db.PostgresDB, alertEngine *alerter.Engine, alertStore alerter.AlertStore) *gin.Engine {
	r := gin.New()

	// Global middleware
	r.Use(middleware.ErrorHandler())
	r.Use(middleware.RequestLogger())
	r.Use(middleware.Cors())
	r.Use(gin.Recovery())

	// Initialize repositories
	var userRepo *repository.UserRepository
	if postgresDB != nil {
		userRepo = repository.NewUserRepository(postgresDB)
	}

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler()
	authHandler := handlers.NewAuthHandler(userRepo)
	assetHandler := handlers.NewAssetHandler(postgresDB)
	telemetryHandler := handlers.NewTelemetryHandler()
	alertHandler := handlers.NewAlertHandler(alertEngine, alertStore)

	// JWT configuration
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default-secret-change-in-production"
	}

	jwtConfig := middleware.JWTConfig{
		Secret: jwtSecret,
		SkipPaths: []string{
			"/health",
			"/health/ready",
			"/api/v1/auth/login",
			"/api/v1/auth/register",
		},
	}

	// Health check routes (no auth required)
	r.GET("/health", healthHandler.HealthCheck)
	r.GET("/health/ready", healthHandler.ReadinessCheck)

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Auth routes (no JWT required)
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/register", authHandler.Register)
			auth.POST("/logout", middleware.JWTAuth(jwtConfig), authHandler.Logout)
			auth.GET("/me", middleware.JWTAuth(jwtConfig), authHandler.GetCurrentUser)
			auth.POST("/refresh", middleware.JWTAuth(jwtConfig), authHandler.RefreshToken)
			auth.POST("/change-password", middleware.JWTAuth(jwtConfig), authHandler.ChangePassword)
		}

		// Protected routes (JWT required)
		protected := v1.Group("")
		protected.Use(middleware.JWTAuth(jwtConfig))
		{
			// Dashboard
			protected.GET("/dashboard/stats", telemetryHandler.GetDashboardStats)

			// Asset routes
			assets := protected.Group("/assets")
			{
				// Data Centers
				assets.GET("/datacenters", assetHandler.ListDataCenters)
				assets.GET("/datacenters/:id", assetHandler.GetDataCenter)

				// Edge Agents
				assets.GET("/agents", assetHandler.ListEdgeAgents)
				assets.GET("/agents/:id", assetHandler.GetEdgeAgent)

				// Servers
				assets.GET("/servers", assetHandler.ListServers)
				assets.GET("/servers/:id", assetHandler.GetServer)

				// Cooling Devices
				assets.GET("/cooling-devices", assetHandler.ListCoolingDevices)
				assets.GET("/cooling-devices/:id", assetHandler.GetCoolingDevice)
			}

			// Telemetry routes
			telemetry := protected.Group("/telemetry")
			{
				telemetry.GET("/readings", telemetryHandler.GetSensorReadings)
				telemetry.GET("/pue", telemetryHandler.GetPueHistory)
			}

			// Alert routes
			alerts := protected.Group("/alerts")
			{
				alerts.GET("", alertHandler.ListAlerts)
				alerts.GET("/:id", alertHandler.GetAlert)
				alerts.POST("/:id/acknowledge", alertHandler.AcknowledgeAlert)
				alerts.POST("/:id/resolve", alertHandler.ResolveAlert)
			}

			// Alert rules (admin only)
			rules := protected.Group("/alert-rules")
			rules.Use(middleware.RequireRole("admin"))
			{
				rules.GET("", alertHandler.ListAlertRules)
				rules.POST("", alertHandler.CreateAlertRule)
				rules.GET("/:id", alertHandler.GetAlertRule)
				rules.PUT("/:id", alertHandler.UpdateAlertRule)
				rules.DELETE("/:id", alertHandler.DeleteAlertRule)
			}
		}
	}

	return r
}
