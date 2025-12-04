package database

import (
	"fmt"
	"log"
	"lomi-backend/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var SqlxDB *sqlx.DB

// ConnectSqlxDB creates a sqlx database connection for the wallet system
func ConnectSqlxDB(cfg *config.Config) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode,
	)

	var err error
	SqlxDB, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database with sqlx: ", err)
	}

	// Connection pooling
	SqlxDB.SetMaxIdleConns(10)
	SqlxDB.SetMaxOpenConns(100)

	log.Println("✅ Connected to PostgreSQL Database (sqlx)")
}
