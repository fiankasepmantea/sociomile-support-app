package model


import ("database/sql/driver"; "encoding/json"; "time")


type JSONB map[string]interface{}
func (j JSONB) Value() (driver.Value, error) { return json.Marshal(j) }
func (j *JSONB) Scan(v interface{}) error {
    if v == nil { *j = make(JSONB); return nil }
    return json.Unmarshal(v.([]byte), j)
}
type ActivityLog struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    TenantID    uint      `gorm:"not null;index" json:"tenant_id"`
    Type        string    `gorm:"size:100;not null;index" json:"type"`
    EntityType  string    `gorm:"size:100;not null;index:idx_activity_logs_entity" json:"entity_type"`
    EntityID    int64     `gorm:"not null;index:idx_activity_logs_entity" json:"entity_id"`
    ActorUserID *uint     `gorm:"index" json:"actor_user_id,omitempty"`
    Payload     JSONB     `gorm:"type:jsonb;not null;default:'{}'" json:"payload"`
    CreatedAt   time.Time `gorm:"index:idx_activity_logs_created" json:"created_at"`
}
