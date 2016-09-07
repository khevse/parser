package db

import (
	_ "github.com/mattn/go-sqlite3"

	"database/sql"
	"database/sql/driver"

	"errors"
	"log"
	"os"
)

const (
	TYPE_sqlite3  string = "sqlite3"  // https://github.com/mattn/go-sqlite3
	TYPE_postgres string = "postgres" // https://github.com/lib/pq
)

var (
	conection *sql.DB     = nil // DB connection
	logger    *log.Logger = log.New(os.Stdout, "\r\n", 0)
	dialect   string
)

func Init(db_type string, settings string, logger_ref *log.Logger) {

	var err error

	dialect = db_type

	switch dialect {
	case TYPE_sqlite3, TYPE_postgres:
		conection, err = sql.Open(dialect, settings)
	default:
		err = errors.New("Unknown type of the database.")
	}

	if err != nil {
		logger.Panicln(err.Error())
	}

	if logger_ref != nil {
		logger = logger_ref
	}
}

func Close() {
	conection.Close()
}

func Exec(query string, tx *sql.Tx, args ...interface{}) (retval driver.Result, err error) {

	if tx == nil {
		retval, err = conection.Exec(query, args...)
	} else {
		retval, err = tx.Exec(query, args...)
	}

	if err != nil {
		logger.Printf("Error: %s:\n%s", query, err.Error())
	}

	return
}

func TxBegin() (tx *sql.Tx, err error) {

	tx, err = conection.Begin()
	if err != nil {
		logger.Printf("Create transaction: %s", err.Error())
	}

	return
}

func TxCommit(tx *sql.Tx) error {

	err := tx.Commit()
	if err != nil {
		logger.Printf("Commit transaction: %s", err.Error())
	}

	return err
}

func TxRollback(tx *sql.Tx) error {

	err := tx.Rollback()
	if err != nil {
		logger.Printf("Rollback transaction: %s", err.Error())
	}

	return err
}
