package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"orderService/configs"
	"os"
	"time"
)

func Connect(cnf configs.Database) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s",
		cnf.User,
		cnf.Password,
		cnf.Host,
		cnf.Port,
		cnf.Name,
		cnf.Schema,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: time.Second,
				LogLevel:      logger.Info,
				Colorful:      true,
			},
		),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database")
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get DB from GORM")
	}

	sqlDB.SetMaxOpenConns(cnf.MaxOpenConnection)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, err
}
