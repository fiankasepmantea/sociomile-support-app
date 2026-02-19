package model

import (
	"time"

	"gorm.io/gorm"
)

type ConversationStatus string

const (
	StatusOpen     ConversationStatus = "open"
	StatusAssigned ConversationStatus = "assigned"
	StatusClosed   ConversationStatus = "closed"
)

type Conversation struct {
	ID              uint
	TenantID        uint
	Tenant          Tenant
	CustomerID      uint
	Customer        Customer
	AssignedAgentID *uint
	AssignedAgent   *User
	Status          ConversationStatus
	LastMessageAt   *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt
	Messages        []Message
}


