package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresDB struct {
	DB *gorm.DB
}

func NewPostgresDB(dsn string) (*PostgresDB, error) {
	var db *gorm.DB
	var err error

	maxRetries := 10

	for i := 1; i <= maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			sqlDB, _ := db.DB()
			if pingErr := sqlDB.Ping(); pingErr == nil {
				log.Println("✅ Database connected")
				return &PostgresDB{DB: db}, nil
			}
		}

		log.Printf("⏳ Waiting database... attempt %d/%d\n", i, maxRetries)
		time.Sleep(3 * time.Second)
	}

	return nil, fmt.Errorf("failed connect database: %w", err)
}

func (p *PostgresDB) Close() {
	sqlDB, err := p.DB.DB()
	if err == nil {
		sqlDB.Close()
	}
}

func (p *PostgresDB) RunSeeder() error {
	seeder := NewSeeder(p.DB)
	return seeder.Run()
}
