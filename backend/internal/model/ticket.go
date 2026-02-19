package model

import "time"

type TicketStatus string

const (
	TicketStatusOpen   TicketStatus = "open"
	TicketStatusClosed TicketStatus = "closed"
)

type Ticket struct {
    ID             uint         `gorm:"primaryKey"`
    TenantID       uint
    ConversationID uint
    Title          string       `gorm:"type:varchar(255);not null"`
    Description    string       `gorm:"type:varchar(255);not null"`
    Status         TicketStatus `gorm:"type:varchar(20)"`
    CreatedAt      time.Time
    UpdatedAt      time.Time
}

