package db

import (
	"fmt"
	"strings"

	_ "github.com/lib/pq"

	"database/sql"
)

type Postgres struct {
	DB

	conn   *sql.DB
	schema string
}

func NewPostgres(connectionString, schema string) (*Postgres, error) {

	conn, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	conn.SetMaxOpenConns(1)

	p := Postgres{
		conn:   conn,
		schema: schema,
	}

	return &p, nil
}

func (p *Postgres) Close() error {
	return p.conn.Close()
}

func (p *Postgres) Driver() driver {
	return DriverProstgres
}

func (p *Postgres) Schema() string {
	return p.schema
}

func (p *Postgres) Connection() interface{} {
	return p.conn
}

func (p *Postgres) CollectionName(table string) string {
	tableName := table
	if s := p.Schema(); len(s) > 0 {
		tableName = fmt.Sprintf("%s.%s", s, table)
	}
	return tableName
}

func (p *Postgres) HasTable(name string) (bool, error) {

	tableName := name
	if len(p.schema) > 0 {
		tableName = fmt.Sprintf("%s.%s", p.schema, name)

		query := fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS %s`, p.schema)
		_, err := p.conn.Exec(query)
		if err != nil {
			return false, err
		}

	}

	query := fmt.Sprintf(`SELECT 1 FROM %s WHERE 1=0`, tableName)

	_, err := p.conn.Exec(query)
	if err != nil {
		if strings.HasSuffix(err.Error(), "does not exist") {
			return false, nil
		} else {
			return false, err
		}
	}

	return true, nil
}

func (p *Postgres) CreateTable(query string) error {

	_, err := p.conn.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) ClearTable(name string) error {

	_, err := p.conn.Exec("DELETE FROM " + p.CollectionName(name))
	return err
}
