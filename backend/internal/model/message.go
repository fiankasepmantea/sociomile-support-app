package model
import "time"

type SenderType string
const (SenderCustomer SenderType = "customer"; SenderAgent SenderType = "agent")
type Message struct {
    ID            uint        `gorm:"primaryKey" json:"id"`
    TenantID      uint        `gorm:"not null;index" json:"tenant_id"`
    ConversationID uint       `gorm:"not null;index:idx_messages_conversation_created" json:"conversation_id"`
    SenderType    SenderType  `gorm:"type:message_sender_type;not null" json:"sender_type"`
    SenderUserID  *uint       `gorm:"index" json:"sender_user_id,omitempty"`
    Body          string      `gorm:"type:text;not null" json:"body"`
    CreatedAt     time.Time   `gorm:"index:idx_messages_conversation_created" json:"created_at"`
    DeletedAt     *time.Time  `gorm:"index" json:"deleted_at,omitempty"`
}
