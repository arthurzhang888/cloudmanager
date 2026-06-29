package repository

import (
	"context"
	"fmt"

	"github.com/cloudmanager/cloud-backend/asset-service/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ServerRepository handles database operations for servers
type ServerRepository struct {
	db *pgxpool.Pool
}

// NewServerRepository creates a new server repository
func NewServerRepository(db *pgxpool.Pool) *ServerRepository {
	return &ServerRepository{db: db}
}

// Create inserts a new server
func (r *ServerRepository) Create(ctx context.Context, server *models.Server) error {
	query := `
		INSERT INTO servers (
			agent_id, redfish_endpoint, manufacturer, model, serial_number, sku,
			cpu_count, cpu_model, memory_gb, disk_count, total_disk_gb,
			status, power_state, health, raw_data
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING id, discovered_at, last_updated, created_at
	`

	rawData := []byte("{}")
	if server.RawData != nil {
		rawData = server.RawData
	}

	return r.db.QueryRow(ctx, query,
		server.AgentID,
		server.RedfishEndpoint,
		server.Manufacturer,
		server.Model,
		server.SerialNumber,
		server.SKU,
		server.CPUCount,
		server.CPUModel,
		server.MemoryGB,
		server.DiskCount,
		server.TotalDiskGB,
		server.Status,
		server.PowerState,
		server.Health,
		rawData,
	).Scan(&server.ID, &server.DiscoveredAt, &server.LastUpdated, &server.CreatedAt)
}

// CreateOrUpdate inserts a new server or updates if serial number exists
func (r *ServerRepository) CreateOrUpdate(ctx context.Context, server *models.Server) error {
	query := `
		INSERT INTO servers (
			agent_id, redfish_endpoint, manufacturer, model, serial_number, sku,
			cpu_count, cpu_model, memory_gb, disk_count, total_disk_gb,
			status, power_state, health, raw_data
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		ON CONFLICT (serial_number) DO UPDATE SET
			agent_id = EXCLUDED.agent_id,
			redfish_endpoint = EXCLUDED.redfish_endpoint,
			status = EXCLUDED.status,
			power_state = EXCLUDED.power_state,
			health = EXCLUDED.health,
			last_updated = NOW()
		RETURNING id, discovered_at, last_updated, created_at
	`

	rawData := []byte("{}")
	if server.RawData != nil {
		rawData = server.RawData
	}

	return r.db.QueryRow(ctx, query,
		server.AgentID,
		server.RedfishEndpoint,
		server.Manufacturer,
		server.Model,
		server.SerialNumber,
		server.SKU,
		server.CPUCount,
		server.CPUModel,
		server.MemoryGB,
		server.DiskCount,
		server.TotalDiskGB,
		server.Status,
		server.PowerState,
		server.Health,
		rawData,
	).Scan(&server.ID, &server.DiscoveredAt, &server.LastUpdated, &server.CreatedAt)
}

// GetByID retrieves a server by ID
func (r *ServerRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Server, error) {
	query := `
		SELECT id, agent_id, redfish_endpoint, manufacturer, model, serial_number, sku,
			cpu_count, cpu_model, memory_gb, disk_count, total_disk_gb,
			status, power_state, health, discovered_at, last_updated, created_at
		FROM servers WHERE id = $1
	`

	server := &models.Server{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&server.ID, &server.AgentID, &server.RedfishEndpoint, &server.Manufacturer,
		&server.Model, &server.SerialNumber, &server.SKU, &server.CPUCount,
		&server.CPUModel, &server.MemoryGB, &server.DiskCount, &server.TotalDiskGB,
		&server.Status, &server.PowerState, &server.Health,
		&server.DiscoveredAt, &server.LastUpdated, &server.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("server not found")
		}
		return nil, err
	}
	return server, nil
}

// List retrieves all servers with optional filtering
func (r *ServerRepository) List(ctx context.Context, filters map[string]interface{}) ([]models.Server, error) {
	query := `
		SELECT id, agent_id, redfish_endpoint, manufacturer, model, serial_number, sku,
			cpu_count, cpu_model, memory_gb, disk_count, total_disk_gb,
			status, power_state, health, discovered_at, last_updated, created_at
		FROM servers
		WHERE 1=1
	`
	args := []interface{}{}
	argCount := 1

	if agentID, ok := filters["agent_id"].(uuid.UUID); ok {
		query += fmt.Sprintf(" AND agent_id = $%d", argCount)
		args = append(args, agentID)
		argCount++
	}
	if status, ok := filters["status"].(string); ok && status != "" {
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, status)
		argCount++
	}

	query += " ORDER BY discovered_at DESC"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanServers(rows)
}

// ListByAgent retrieves servers by agent ID
func (r *ServerRepository) ListByAgent(ctx context.Context, agentID uuid.UUID) ([]models.Server, error) {
	query := `
		SELECT id, agent_id, redfish_endpoint, manufacturer, model, serial_number, sku,
			cpu_count, cpu_model, memory_gb, disk_count, total_disk_gb,
			status, power_state, health, discovered_at, last_updated, created_at
		FROM servers
		WHERE agent_id = $1
		ORDER BY discovered_at DESC
	`

	rows, err := r.db.Query(ctx, query, agentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanServers(rows)
}

// UpdateStatus updates server status
func (r *ServerRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status, powerState, health string) error {
	query := `
		UPDATE servers
		SET status = $1, power_state = $2, health = $3, last_updated = NOW()
		WHERE id = $4
	`
	_, err := r.db.Exec(ctx, query, status, powerState, health, id)
	return err
}

// Delete removes a server
func (r *ServerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM servers WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

// Count returns the total number of servers
func (r *ServerRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM servers").Scan(&count)
	return count, err
}

// CountByStatus returns server counts grouped by status
func (r *ServerRepository) CountByStatus(ctx context.Context) (map[string]int64, error) {
	query := `SELECT status, COUNT(*) FROM servers GROUP BY status`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]int64)
	for rows.Next() {
		var status string
		var count int64
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		result[status] = count
	}
	return result, rows.Err()
}

func (r *ServerRepository) scanServers(rows pgx.Rows) ([]models.Server, error) {
	var servers []models.Server
	for rows.Next() {
		var s models.Server
		if err := rows.Scan(
			&s.ID, &s.AgentID, &s.RedfishEndpoint, &s.Manufacturer, &s.Model,
			&s.SerialNumber, &s.SKU, &s.CPUCount, &s.CPUModel, &s.MemoryGB,
			&s.DiskCount, &s.TotalDiskGB, &s.Status, &s.PowerState, &s.Health,
			&s.DiscoveredAt, &s.LastUpdated, &s.CreatedAt,
		); err != nil {
			return nil, err
		}
		servers = append(servers, s)
	}
	return servers, rows.Err()
}
