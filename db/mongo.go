package db

import (
	"strings"

	"gopkg.in/mgo.v2"
)

type Mongo struct {
	DB

	conn   *mgo.Session
	dbName string
}

func NewMongo(info *mgo.DialInfo, dbName string) (*Mongo, error) {

	conn, err := mgo.DialWithInfo(info)
	if err != nil {
		return nil, err
	}

	conn.SetMode(mgo.Monotonic, true)

	m := Mongo{
		conn:   conn,
		dbName: dbName,
	}

	return &m, nil
}

func (m *Mongo) Close() error {
	m.conn.Close()
	return nil
}

func (m *Mongo) Driver() driver {
	return DriverMongo
}

func (m *Mongo) Schema() string {
	return m.dbName
}

func (m *Mongo) Connection() interface{} {
	return m.conn
}

func (m *Mongo) CollectionName(table string) string {
	return table
}

func (m *Mongo) HasTable(name string) (bool, error) {

	list, err := m.conn.DB(m.dbName).CollectionNames()
	if err != nil {
		return false, err
	}

	for _, tableName := range list {
		if strings.EqualFold(tableName, name) {
			return true, nil
		}
	}

	return false, nil
}

func (m *Mongo) CreateTable(name string) error {
	return m.conn.DB(m.dbName).C(name).Create(new(mgo.CollectionInfo))
}

func (m *Mongo) ClearTable(name string) error {
	_, err := m.conn.DB(m.dbName).C(name).RemoveAll(nil)
	return err
}
