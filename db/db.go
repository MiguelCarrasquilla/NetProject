package db

import (
	"database/sql"
	"dbconnection/config"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func Connect(cfg *config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("Error al abrir la base de datos: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("Error al conectar a la base de datos: %v", err)
	}

	log.Println("Conexión a la base de datos exitosa.")
	return db, nil
}
