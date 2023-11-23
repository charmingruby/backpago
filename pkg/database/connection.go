package database

import (
	"database/sql"
	"fmt"

	"github.com/charmingruby/backpago/configs"
	_ "github.com/lib/pq"
)

func NewConnection() (*sql.DB, error) {

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		configs.Env.DBHost,
		configs.Env.DBPort,
		configs.Env.DBUser,
		configs.Env.DBPassword,
		configs.Env.DBName,
	)

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil

}
