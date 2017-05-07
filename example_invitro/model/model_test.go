package model

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"testing"

	driver "github.com/khevse/parser/db"
	"github.com/ncw/swift"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func TestFiles(t *testing.T) {
	log.SetFlags(log.Lshortfile)

	conn := swift.Connection{
		UserName: "test:tester",
		ApiKey:   "testing",
		AuthUrl:  "http://127.0.0.1:12345/auth/v1.0",
	}

	if err := conn.Authenticate(); err != nil {
		t.Fatal(err)
	}

	const CONTAINER_NAME = "swift"

	m := NewModelFiles(&conn, CONTAINER_NAME)

	tmpFileName := "tmp"
	tmpFileData := []byte("model")
	contentType := http.DetectContentType(tmpFileData)
	if err := m.conn.ObjectPutBytes(CONTAINER_NAME, tmpFileName, tmpFileData, contentType); err != nil {
		t.Fatal(err)
	}

	if name, err := m.Add([]byte("hello")); err != nil {
		t.Fatal(err)
	} else if list, err := m.conn.ObjectNames(m.container, nil); err != nil {
		t.Fatal(err)
	} else {
		names := make(map[string]bool)
		for _, n := range list {
			names[n] = false
		}

		if _, ok := names[tmpFileName]; !ok {
			t.Fatal("Fail remove")
		} else if _, ok := names[name]; !ok {
			t.Fatal("Fail add")
		}

		if err := m.removeAll(); err != nil {
			t.Fatal(err)
		} else if list, err := m.conn.ObjectNames(m.container, nil); err != nil {
			t.Fatal(err)
		} else {
			names = make(map[string]bool)
			for _, n := range list {
				names[n] = false
			}

			if _, ok := names[tmpFileName]; !ok {
				t.Fatal("Fail remove")
			} else if _, ok := names[name]; ok {
				t.Fatal("remove")
			}

		}
	}

	if data, err := m.Get(tmpFileName); err != nil {
		t.Error(err)
	} else if !bytes.Equal(tmpFileData, data) {
		t.Errorf("'%s' != '%s'", string(tmpFileData), string(data))
	}
}

func TestPG(t *testing.T) {

	db := getPgDriver()

	defer func() {
		conn := db.Connection().(*sql.DB)
		query := fmt.Sprintf(`DROP SCHEMA IF EXISTS %s CASCADE`, db.Schema())
		if _, err := conn.Exec(query); err != nil {
			t.Error(err)
		}
	}()

	{
		m := NewModelPrice(db)
		if err := m.internalInit(); err != nil {
			t.Error(err)
		}
	}

	{
		m := NewModelDesc(db)
		if err := m.internalInit(); err != nil {
			t.Error(err)
		}
	}

	{
		mPriсe := NewModelPrice(db)
		mDesc := NewModelDesc(db)

		if idGroup, err := mPriсe.AddGroup(0, 1, "Group #1"); err != nil {
			t.Error(err)
		} else if idItem, err := mPriсe.AddItem(idGroup.(int64), "Item #1", "10", 122.55); err != nil {
			t.Error(err)
		} else if idGroup == 0 {
			t.Error("Fail")
		} else if idItem == 0 {
			t.Error("Fail")
		} else if idDesc, err := mDesc.AddItem(idItem.(int64), "MAIN", "TEXT"); err != nil {
			t.Error(err)
		} else if idDesc == 0 {
			t.Error("Fail")
		}
	}
}

func TestInternalInitMоngo(t *testing.T) {

	db := getMongoDriver()

	defer func() {
		conn := db.Connection().(*mgo.Session)
		if err := conn.DB(db.Schema()).DropDatabase(); err != nil {
			t.Error(err)
		}
	}()

	{
		m := NewModelPrice(db)
		if err := m.internalInit(); err != nil {
			t.Error(err)
		}
	}

	{
		m := NewModelDesc(db)
		if err := m.internalInit(); err != nil {
			t.Error(err)
		}
	}

	{
		mPriсe := NewModelPrice(db)
		mDesc := NewModelDesc(db)

		if idGroup, err := mPriсe.AddGroup(0, 1, "Group #1"); err != nil {
			t.Error(err)
		} else if idItem, err := mPriсe.AddItem(idGroup.(bson.ObjectId), "Item #1", "10", 122.55); err != nil {
			t.Error(err)
		} else if idGroup == 0 {
			t.Error("Fail")
		} else if idItem == 0 {
			t.Error("Fail")
		} else if idDesc, err := mDesc.AddItem(idItem.(bson.ObjectId), "MAIN", "TEXT"); err != nil {
			t.Error(err)
		} else if idDesc == 0 {
			t.Error("Fail")
		}
	}
}

func getPgDriver() driver.DB {
	log.SetFlags(log.Lshortfile)

	const SCHEMA = "invitro"

	db, err := driver.NewPostgres("host=localhost port=5432 user=postgres dbname=postgres sslmode=disable", SCHEMA)
	if err != nil {
		panic(err)
	}

	return db
}

func getMongoDriver() driver.DB {
	log.SetFlags(log.Lshortfile)

	const DATABASE = "localhost"

	info := mgo.DialInfo{
		Addrs:    []string{"localhost:27017"},
		Database: DATABASE,
	}

	db, err := driver.NewMongo(&info, DATABASE)
	if err != nil {
		panic(err)
	}

	return db
}
