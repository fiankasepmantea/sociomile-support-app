package model

import (
    "time"
    "gorm.io/gorm"
)

type Tenant struct {
    ID        uint64         `gorm:"primaryKey;column:id"`
    Name      string         `gorm:"column:name;unique;not null"`
    CreatedAt time.Time      `gorm:"column:created_at"`
    UpdatedAt time.Time      `gorm:"column:updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (Tenant) TableName() string {
    return "tenants"
}
