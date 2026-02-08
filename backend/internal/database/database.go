package database

import (
	"fmt"
	"log"

	"github.com/guatom999/self-boardcast/internal/config"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

func DatabaseConnect(cfg *config.Config) *sqlx.DB {

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		panic(err)
	}

	if err := db.Ping(); err != nil {
		log.Printf("Error pinging database: %v", err)
		panic(err)
	}

	db.SetMaxOpenConns(20)

	return db

}
