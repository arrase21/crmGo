package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/arrase21/crm-users/internal/config"
	"github.com/arrase21/crm-users/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(cfg *config.PostgresConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.SSLMode,
		cfg.TimeZone,
	)
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:      gormLogger,
		PrepareStmt: true,
	})
	if err != nil {
		return nil, err
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	return db, nil
}

// automigrate
func Automigrate(db *gorm.DB) error {
	log.Println("🔄 Running database migrations...")
	err := db.AutoMigrate(
		&domain.User{},
		&domain.Permission{},
		&domain.PermissionAction{},
		&domain.Role{},
		&domain.RolePermission{},
		&domain.UserRole{},
	)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %s", err)
	}
	log.Println("✅ database migrations complete")
	return nil
}

// close database conection
func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
