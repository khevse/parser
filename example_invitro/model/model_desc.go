package model

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/khevse/parser/db"
	"gopkg.in/mgo.v2/bson"
)

const TABLE_DESC = "INVITRO_PRICE_DESCRIPTION"

type ModelDesc struct {
	DB        db.DB
	once      sync.Once
	insertReq *sql.Stmt
}

func NewModelDesc(database db.DB) *ModelDesc {
	return &ModelDesc{
		DB: database,
	}
}

func (m *ModelDesc) AddItem(positionId interface{}, section, content string) (interface{}, error) {
	if err := m.internalInit(); err != nil {
		return nil, err
	}

	if m.DB.Driver() == db.DriverProstgres {
		return add(m.DB, m.insertReq, positionId, section, content)

	} else if m.DB.Driver() == db.DriverMongo {
		newId := bson.NewObjectId()
		args := bson.M{
			"_id":         newId,
			"position_id": positionId,
			"section":     section,
			"content":     content,
		}

		return add(m.DB, TABLE_DESC, args)
	}

	return nil, errors.New("Wrong DB driver")
}

func (m *ModelDesc) internalInit() (retval error) {

	m.once.Do(func() {
		switch m.DB.Driver() {
		case db.DriverProstgres:
			retval = m.internalInitPG()
		case db.DriverMongo:
			retval = m.internalInitMongo()
		default:
			panic("Unknown db driver.")
		}
	})

	return
}

func (m *ModelDesc) internalInitPG() error {

	tableName := m.DB.CollectionName(TABLE_DESC)

	if has, err := m.DB.HasTable(TABLE_DESC); err != nil {
		return err
	} else if !has {
		query := fmt.Sprintf(`
            CREATE TABLE %s (
              id SERIAL PRIMARY KEY,
              position_id INTEGER NOT NULL,
              section VARCHAR(100) NOT NULL,
              content TEXT NOT NULL
            );`,
			tableName)

		if err := m.DB.CreateTable(query); err != nil {
			return err
		}
	} else if has {
		if err := m.DB.ClearTable(TABLE_DESC); err != nil {
			log.Println(err)
			return err
		}
	}

	{
		query := fmt.Sprintf(`
				INSERT INTO %s(position_id, section, content)
				        values($1, $2, $3)
				         RETURNING id`, tableName)
		if val, err := m.DB.Connection().(*sql.DB).Prepare(query); err != nil {
			return err
		} else {
			m.insertReq = val
		}
	}

	return nil
}

func (m *ModelDesc) internalInitMongo() error {

	if has, err := m.DB.HasTable(TABLE_DESC); err != nil {
		return err
	} else if !has {
		if err := m.DB.CreateTable(TABLE_DESC); err != nil {
			return err
		}
	} else if has {
		if err := m.DB.ClearTable(TABLE_DESC); err != nil {
			return err
		}
	}

	return nil
}
