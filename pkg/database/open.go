package database

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func OpenDB(driver string, dsn string) (*sql.DB, error) {
	log.Println(driver, dsn, "------------------------")
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	if err := db.PingContext(context.Background()); err != nil {
		log.Println(err, "-----")
		return nil, err
	}

	return db, nil

}
