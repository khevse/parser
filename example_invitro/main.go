package main

import (
	"log"
	"os"
	"time"

	"github.com/khevse/parser/db"
	"github.com/khevse/parser/engine"
	"github.com/khevse/parser/example_invitro/controller"
	"github.com/khevse/parser/example_invitro/model"
	"github.com/khevse/parser/workers"
	"github.com/ncw/swift"
	"gopkg.in/mgo.v2"
)

func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds | log.Lshortfile)

	pgConnection := os.Getenv("PG_CONN")
	pgSchema := os.Getenv("PG_SCHEMA")
	mongoAddress := os.Getenv("MONGO_ADDRESS")
	mongoDBName := os.Getenv("MONGO_DBNAME")
	swiftUser := os.Getenv("SWIFT_USER")
	swiftApiKey := os.Getenv("SWIFT_API_KEY")
	swiftAuthUrl := os.Getenv("SWIFT_AUTH_URL")
	swiftContainer := os.Getenv("SWIFT_CONTAINER")

	var database db.DB
	if len(pgConnection) > 0 {
		database = getPgDriver(pgConnection, pgSchema)
	} else if len(mongoAddress) > 0 {
		database = getMongoDriver(mongoAddress, mongoDBName)
	}

	swiftConn := getSwiftConnection(swiftUser, swiftApiKey, swiftAuthUrl)

	modelPrice := model.NewModelPrice(database)
	modelDesc := model.NewModelDesc(database)
	modelFiles := model.NewModelFiles(swiftConn, swiftContainer)

	d := new(workers.Dispatcher)
	d.Run()
	defer d.Close()

	e, err := engine.NewEngine("https://www.invitro.ru/analizes/for-doctors", d)
	if err != nil {
		log.Fatal(err)
	}

	res, err := e.Parse(controller.TOCHandler)
	if err != nil {
		log.Fatal(err)
	} else {
		for item := range res {
			target := item.(*controller.Target)

			if err := target.Save(modelPrice, modelDesc, modelFiles); err != nil {
				log.Println("ERROR:", err)
				return
			} else {
				log.Println("LOG:", "level:", target.Level, "name:", target.Name, "amount:", target.Amount)
			}
		}
	}

}

func getPgDriver(conn, schema string) db.DB {
	// Parameters example:
	// - conn   - "host=localhost port=5432 user=postgres dbname=postgres sslmode=disable"
	// - schema - "invitro"

	db, err := db.NewPostgres(conn, schema)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func getMongoDriver(address, dbname string) db.DB {
	// Parameters example:
	// - address - "localhost:27017"
	// - dbname  - "invitro"

	info := mgo.DialInfo{
		Addrs:    []string{address},
		Database: dbname,
	}

	db, err := db.NewMongo(&info, dbname)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func getSwiftConnection(user, apiKey, url string) *swift.Connection {
	// Parameters example:
	// - user   - "test:tester"
	// - apiKey - "testing"
	// - url    - "http://127.0.0.1:12345/auth/v1.0",

	conn := swift.Connection{
		UserName: user,
		ApiKey:   apiKey,
		AuthUrl:  url,
	}

	var lastError error
	for i := 0; i < 20; i++ {
		if err := conn.Authenticate(); err != nil {
			lastError = err
			log.Println("WAIT swift:", err, "(", user, apiKey, url, ")")
		} else {
			lastError = nil
			break
		}

		time.Sleep(time.Second * 5)
	}

	if lastError != nil {
		log.Fatal(lastError)
	}

	return &conn
}
