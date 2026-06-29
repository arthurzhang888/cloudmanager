package repository

import (
	"context"
	"errors"
	"time"

	"github.com/cloudmanager/cloud-backend/shared/db"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// User represents a user in the system
type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Never serialize password
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserRepository handles user data access
type UserRepository struct {
	db *db.PostgresDB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *db.PostgresDB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *User) error {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	query := `
		INSERT INTO users (id, username, email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err = r.db.Exec(ctx, query,
		user.ID,
		user.Username,
		user.Email,
		string(hashedPassword),
		user.Role,
		user.CreatedAt,
		user.UpdatedAt,
	)

	return err
}

// GetByID gets a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id string) (*User, error) {
	query := `
		SELECT id, username, email, role, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := &User{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetByUsername gets a user by username
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*User, error) {
	query := `
		SELECT id, username, email, password_hash, role, created_at, updated_at
		FROM users
		WHERE username = $1
	`

	user := &User{}
	var passwordHash string
	err := r.db.QueryRow(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&passwordHash,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	user.Password = passwordHash
	return user, nil
}

// GetByEmail gets a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, username, email, role, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	user := &User{}
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// ValidateCredentials validates user credentials
func (r *UserRepository) ValidateCredentials(ctx context.Context, username, password string) (*User, error) {
	user, err := r.GetByUsername(ctx, username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Clear password before returning
	user.Password = ""
	return user, nil
}

// UpdatePassword updates a user's password
func (r *UserRepository) UpdatePassword(ctx context.Context, userID, newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := `
		UPDATE users
		SET password_hash = $1, updated_at = $2
		WHERE id = $3
	`

	_, err = r.db.Exec(ctx, query, string(hashedPassword), time.Now(), userID)
	return err
}

// Update updates user information
func (r *UserRepository) Update(ctx context.Context, user *User) error {
	user.UpdatedAt = time.Now()

	query := `
		UPDATE users
		SET username = $1, email = $2, role = $3, updated_at = $4
		WHERE id = $5
	`

	_, err := r.db.Exec(ctx, query,
		user.Username,
		user.Email,
		user.Role,
		user.UpdatedAt,
		user.ID,
	)

	return err
}

// Delete deletes a user
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

// List lists all users
func (r *UserRepository) List(ctx context.Context) ([]*User, error) {
	query := `
		SELECT id, username, email, role, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := &User{}
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, rows.Err()
}

// CreateUsersTable creates the users table if it doesn't exist
func (r *UserRepository) CreateUsersTable(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY,
			username VARCHAR(32) UNIQUE NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			role VARCHAR(20) NOT NULL DEFAULT 'user',
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)
	`
	_, err := r.db.Exec(ctx, query)
	return err
}

// SeedAdminUser creates an admin user if no users exist
func (r *UserRepository) SeedAdminUser(ctx context.Context) error {
	// Check if any users exist
	var count int
	err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return nil // Users already exist
	}

	// Create admin user
	admin := &User{
		Username: "admin",
		Email:    "admin@example.com",
		Password: "admin",
		Role:     "admin",
	}

	return r.Create(ctx, admin)
}
