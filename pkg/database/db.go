package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"stlab.itechart-group.com/go/food_delivery/authentication_service/pkg/logging"
)

type PostgresDB struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
	logger   logging.Logger
}

func NewPostgresDB(database PostgresDB) (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		database.Username, database.Password, database.Host, database.Port, database.DBName, database.SSLMode))
	if err != nil {

		return nil, fmt.Errorf("error connecting to database:%s", err)
	}
	err = db.Ping()
	if err != nil {
		database.logger.Errorf("DB ping error:%s", err)
		return nil, err
	}
	_, err = db.Exec(USER_SCHEMA)
	if err != nil {
		database.logger.Errorf("Error executing initial migration into users:%s", err)
		return nil, fmt.Errorf("error executing initial migration into users:%s", err)
	}
	return db, nil
}

const USER_SCHEMA = `
	CREATE TABLE IF NOT EXISTS users (
		id serial not null primary key ,
		email varchar(225) NOT NULL UNIQUE,
		password varchar(225) NOT NULL,
	    role varchar(50) NOT NULL,
	  	created_at date NOT NULL,
	    deleted bool NOT NULL
	);	
`
