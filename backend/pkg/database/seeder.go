package database

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Seeder struct {
	db *gorm.DB
}

func NewSeeder(db *gorm.DB) *Seeder {
	return &Seeder{db: db}
}

func (s *Seeder) Run() error {
	log.Println("🌱 Running seeders...")
	if err := s.seedTenants(); err != nil {
		return err
	}
	if err := s.seedUsers(); err != nil {
		return err
	}
	log.Println("✅ Seeders completed")
	return nil
}

func (s *Seeder) seedTenants() error {
	tenants := []struct{ Name string }{
		{Name: "test-tenant"},
		{Name: "demo-company"},
	}
	for _, tenant := range tenants {
		var id int64
		err := s.db.Raw(`INSERT INTO tenants (name) VALUES ($1) ON CONFLICT (name) DO NOTHING RETURNING id`, tenant.Name).Scan(&id).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			var exists bool
			s.db.Raw("SELECT EXISTS(SELECT 1 FROM tenants WHERE name = $1)", tenant.Name).Scan(&exists)
			if exists {
				log.Printf("⏭️  Tenant already exists: %s", tenant.Name)
				continue
			}
			return err
		}
		if id > 0 {
			log.Printf("✅ Seeded tenant: %s (ID: %d)", tenant.Name, id)
		}
	}
	return nil
}

func (s *Seeder) seedUsers() error {
	plainPassword := "123456"

	passwordHashBytes, err := bcrypt.GenerateFromPassword(
		[]byte(plainPassword),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	passwordHash := string(passwordHashBytes)

	users := []struct {
		TenantName, Email, Role, Name, PasswordHash string
	}{
		{"test-tenant", "admin@gmail.com", "admin", "Test Admin", passwordHash},
		{"test-tenant", "agent@gmail.com", "agent", "Test Agent", passwordHash},
	}

	for _, user := range users {
		var tenantID int64

		err := s.db.Raw(
			"SELECT id FROM tenants WHERE name = $1",
			user.TenantName,
		).Scan(&tenantID).Error

		if err != nil {
			continue
		}

		var exists bool
		s.db.Raw(
			`SELECT EXISTS(
				SELECT 1 FROM users
				WHERE email = $1 AND tenant_id = $2
			)`,
			user.Email,
			tenantID,
		).Scan(&exists)

		if exists {
			log.Printf("⏭️  User already exists: %s", user.Email)
			continue
		}

		var userID int64
		err = s.db.Raw(`
			INSERT INTO users
			(tenant_id, email, password_hash, role, name, is_active)
			VALUES ($1, $2, $3, $4, $5, true)
			RETURNING id
		`,
			tenantID,
			user.Email,
			user.PasswordHash,
			user.Role,
			user.Name,
		).Scan(&userID).Error

		if err != nil {
			return err
		}

		log.Printf("✅ Seeded user: %s (%s)", user.Email, user.Role)
	}

	return nil
}

