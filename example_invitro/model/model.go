package model

import (
	"database/sql"
	"log"

	"github.com/khevse/parser/db"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func add(database db.DB, req interface{}, args ...interface{}) (interface{}, error) {

	var id interface{}

	if database.Driver() == db.DriverProstgres {
		var newId int64
		if err := req.(*sql.Stmt).QueryRow(args...).Scan(&newId); err != nil {
			log.Println(err)
			return nil, err
		} else {
			id = newId
		}

	} else {
		newId := bson.NewObjectId()
		arg := args[0].(bson.M)
		arg["_id"] = newId

		conn := database.Connection().(*mgo.Session)
		if err := conn.DB(database.Schema()).C(req.(string)).Insert(arg); err != nil {
			log.Println(err)
			return nil, err
		} else {
			id = newId
		}
	}

	return id, nil
}
