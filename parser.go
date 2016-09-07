package main

import (
	"./src/data"
	"./src/db"
	"./src/parser"
	"database/sql"
	"fmt"
	"sync"
)

func main() {

	db.Init(db.TYPE_sqlite3, "database.db", nil)
	defer db.Close()

	var err error

	err = db.CreateDBTable()
	if err != nil {
		panic(err.Error())
	}

	host := "https://www.invitro.ru"
	mainPage := parser.DownloadHtml(host, "/analizes/for-doctors/")

	groupTx, err := db.TxBegin()
	if err != nil {
		panic(err.Error())
	}

	var wg sync.WaitGroup

	groups_list := data.ReadPriceGroups(&mainPage)
	for i, _ := range groups_list {
		wg.Add(1)

		go func(group *data.Group) {
			defer wg.Done()

			groupPage := parser.DownloadHtml(host, group.Href)
			groupPrice, err := data.ReadPrice(group.Number, group.Name, groupPage)
			if err != nil {
				panic(err.Error())
			}

			err = writePrise(groupTx, groupPrice)
			if err != nil {
				panic(err.Error())
			}

		}(&groups_list[i])

		if i == 8 {
			break // TODO разобраться с шаблоном, т.к. он немного отличается от остальных тем что отсутвует тэг tbody
		}
	}

	wg.Wait()

	if err == nil {
		err = db.TxCommit(groupTx)
	}

	if err != nil {
		fmt.Println(err.Error())
		db.TxRollback(groupTx)
	}
}

func writePrise(tx *sql.Tx, list []*data.Price) error {

	var err error

	for i, _ := range list {
		val := list[i]

		var parentNumber *string = nil
		if val.Parent != nil {
			parentNumber = &val.Parent.Number
		}

		err = db.AddRow(tx, val.Number, val.Name, val.Amount, parentNumber)
		if err == nil && len(val.Children) > 0 {
			err = writePrise(tx, val.Children)
		}

		if err != nil {
			db.TxRollback(tx)
			fmt.Println(err.Error())
			break
		}
	}

	return err
}
