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

const TABLE_PRICE = "INVITRO_PRICE"

type ModelPrice struct {
	DB             db.DB
	once           sync.Once
	insertReqGroup *sql.Stmt
	insertReqItem  *sql.Stmt
}

func NewModelPrice(database db.DB) *ModelPrice {
	return &ModelPrice{
		DB: database,
	}
}

func (m *ModelPrice) AddGroup(parentId interface{}, level int, name string) (interface{}, error) {
	return m.add(parentId, level, name, "", 0)
}

func (m *ModelPrice) AddItem(parentId interface{}, name, priceCode string, amount float64) (interface{}, error) {
	return m.add(parentId, 0, name, priceCode, amount)
}

func (m *ModelPrice) add(parentId interface{}, level int, name, priceCode string, amount float64) (interface{}, error) {
	if err := m.internalInit(); err != nil {
		return nil, err
	}

	isGroup := level > 0

	if m.DB.Driver() == db.DriverProstgres {

		var parentIdInt int64
		if parentId != nil {
			switch parentId.(type) {
			case int64:
				parentIdInt = parentId.(int64)
			case int:
				parentIdInt = int64(parentId.(int))
			default:
				log.Fatalf("Wrong type of the patent id: %T (%v)", parentId, parentId)
			}
		}

		if parentIdInt == 0 {
			parentId = new(sql.NullInt64)
		} else {
			parentId = parentIdInt
		}

		if isGroup {
			return add(m.DB, m.insertReqGroup, parentId, level, name)
		} else {
			return add(m.DB, m.insertReqItem, parentId, name, priceCode, amount)
		}

	} else if m.DB.Driver() == db.DriverMongo {
		newId := bson.NewObjectId()
		args := bson.M{
			"_id":       newId,
			"parent_id": parentId,
			"name":      name,
		}

		if isGroup {
			args["level"] = level
			args["price_code"] = nil
			args["amount"] = nil
		} else {
			args["level"] = nil
			args["price_code"] = priceCode
			args["amount"] = amount
		}

		if parentId == 0 {
			args["parent_id"] = nil
		}

		return add(m.DB, TABLE_PRICE, args)
	}

	return nil, errors.New("Wrong DB driver")
}

func (m *ModelPrice) internalInit() (retval error) {

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

func (m *ModelPrice) internalInitPG() error {

	tableName := m.DB.CollectionName(TABLE_PRICE)

	if has, err := m.DB.HasTable(TABLE_PRICE); err != nil {
		return err
	} else if !has {
		query := fmt.Sprintf(`
            CREATE TABLE %s (
              id SERIAL PRIMARY KEY,
              parent_id INTEGER NULL,
              level INTEGER NULL,
              name VARCHAR(500) NOT NULL,
              price_code VARCHAR(20) NULL,
              amount DECIMAL NULL
            );`,
			tableName)

		if err := m.DB.CreateTable(query); err != nil {
			return err
		}
	} else if has {
		if err := m.DB.ClearTable(TABLE_PRICE); err != nil {
			log.Println(err)
			return err
		}
	}

	{
		query := fmt.Sprintf(`
				INSERT INTO %s(parent_id, level, name, price_code, amount)
				        values($1, $2, $3, NULL, NULL)
				         RETURNING id`, tableName)
		if val, err := m.DB.Connection().(*sql.DB).Prepare(query); err != nil {
			return err
		} else {
			m.insertReqGroup = val
		}
	}
	{
		query := fmt.Sprintf(`
				INSERT INTO %s(parent_id, level, name, price_code, amount)
				        values($1, 0, $2, $3, $4)
				         RETURNING id`, tableName)
		if val, err := m.DB.Connection().(*sql.DB).Prepare(query); err != nil {
			return err
		} else {
			m.insertReqItem = val
		}
	}

	return nil
}

func (m *ModelPrice) internalInitMongo() (retval error) {

	if has, err := m.DB.HasTable(TABLE_PRICE); err != nil {
		return err
	} else if !has {
		if err := m.DB.CreateTable(TABLE_PRICE); err != nil {
			return err
		}
	} else if has {
		if err := m.DB.ClearTable(TABLE_PRICE); err != nil {
			return err
		}
	}

	return nil
}
