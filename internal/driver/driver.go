package driver

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DB struct {
	SQL *sql.DB
}

var dbConn = DB{}

const (
	MaxOpenDbConn     = 10
	MaxIdleDbCOnn     = 5
	MaxDbConnLifeTime = 5 * time.Minute
)

func ConnectSql(dsn string) (*DB, error) {
	db, err := NewConnection(dsn)
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(MaxOpenDbConn)
	db.SetMaxIdleConns(MaxIdleDbCOnn)
	db.SetConnMaxLifetime(MaxDbConnLifeTime)
	dbConn.SQL = db
	if testDb(db); err != nil {
		return nil, err
	}
	return &dbConn, nil
}
func testDb(db *sql.DB) error {
	err := db.Ping()
	if err != nil {
		return err
	}
	return nil
}
func NewConnection(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
