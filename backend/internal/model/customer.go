package model
import "time"
type Customer struct {
    ID         uint       `gorm:"primaryKey" json:"id"`
    TenantID   uint       `gorm:"not null;index" json:"tenant_id"`
    ExternalID string     `gorm:"size:255;not null;uniqueIndex:idx_customers_tenant_external" json:"external_id"`
    Name       *string    `gorm:"size:255" json:"name,omitempty"`
    CreatedAt  time.Time  `json:"created_at"`
    DeletedAt  *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}
