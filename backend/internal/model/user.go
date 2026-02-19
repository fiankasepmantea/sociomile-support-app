package model

import "time"

type Role string

const (
    RoleAdmin Role = "admin"
    RoleAgent Role = "agent"
)

type User struct {
	ID           uint      `json:"id"`
	TenantID     uint      `json:"tenant_id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-" gorm:"column:password_hash"`
	Role         string    `json:"role"`
	Name         string    `json:"name"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}


