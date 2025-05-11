package driver

import (
	"database/sql"
	"time"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

//DB holds the Database connection pool
type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

const maxOpenDbConn = 10
const maxIdleDbConn = 5
const maxDbLifeTime = 5 * time.Minute

//Creates the DB connection for PotstGres
func ConnectSQL(connString string) (*DB, error) {
	d, err := NewDataBase(connString)
	if err != nil {
		panic(err)
	}
	d.SetMaxOpenConns(maxOpenDbConn)
	d.SetMaxIdleConns(maxIdleDbConn)
	d.SetConnMaxLifetime(maxDbLifeTime)

	dbConn.SQL = d

	err = testDB(d)
	if err != nil {
		return nil, err
	}
	return dbConn, nil

}

//Tries to ping the DB
func testDB(d *sql.DB) error{
	err := d.Ping()
	if err != nil {
		return err
	}
	return nil
}

//Creates a new Database connection for the application
func NewDataBase(connString string) (*sql.DB, error){
	db, err := sql.Open("pgx", connString)
	if (err != nil) {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}