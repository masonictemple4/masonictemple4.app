package db

import (
	"fmt"
	"os"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"

	"github.com/masonictemple4/masonictemple4.app/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New(config *gorm.Config) (*gorm.DB, error) {

	var db *gorm.DB
	var err error

	// Treat "" and "false" as false
	// Treat everything else as true.
	useConn := os.Getenv("USE_CLOUD_SQL_CONNECTOR")

	if useConn == "" || useConn == "false" {
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_PORT"),
			os.Getenv("SSL_MODE"),
		)
		db, err = gorm.Open(postgres.Open(dsn), config)
		if err != nil {
			return nil, err
		}
	} else {
		instanceConnectionName := os.Getenv("DB_CONNECTION_NAME")
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s sslmode=%s sslrootcert=%s sslcert=%s, sslkey=%s",
			instanceConnectionName,
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
			os.Getenv("SSL_MODE"),
			os.Getenv("SSL_ROOT"),
			os.Getenv("SSL_CERT"),
			os.Getenv("SSL_KEY"),
		)

		db, err = gorm.Open(postgres.New(postgres.Config{
			DriverName: "cloudsqlpostgres",
			DSN:        dsn,
		}))

	}

	if err != nil {
		return nil, err
	}

	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Comment{},
		&models.Tag{},
		&models.Media{},
	)

	return err
}
