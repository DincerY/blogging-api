package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Init() error {
	var err error
	connStr := "host=localhost port=5432 user=postgres password=123456 dbname=BloggingDb sslmode=disable"

	DB, err = sql.Open("postgres", connStr)

	if err != nil {
		return fmt.Errorf("veritabanı bağlantı hatası:", err)
	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("veritabanı erişim hatası:", err)
	}
	return nil
}
