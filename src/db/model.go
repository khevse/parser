package db

import (
	"database/sql"
)

func CreateDBTable() error {

	tx, err := TxBegin()
	if err != nil {
		return err
	}

	query := `
        DROP TABLE IF EXISTS 'price';

        CREATE TABLE IF NOT EXISTS 'price' (
          "id" INTEGER PRIMARY KEY AUTOINCREMENT,
          "number" VARCHAR(20) NOT NULL,
          "name" VARCHAR(500) NOT NULL,
          "amount" DECIMAL UNSIGNED NULL,
          "number_group" VARCHAR(20) NULL
        );
        `

	_, err = Exec(query, tx)
	if err != nil {
		TxRollback(tx)
		return err
	}

	err = TxCommit(tx)
	if err != nil {
		TxRollback(tx)
		return err
	}

	return nil
}

func AddRow(tx *sql.Tx, number string, name string, amount float32, number_group *string) error {

	args := []interface{}{number, name}

	if amount != 0 {
		args = append(args, amount)
	} else {
		args = append(args, sql.NullFloat64{})
	}

	if number_group != nil {
		args = append(args, *number_group)
	} else {
		args = append(args, sql.NullString{})
	}

	_, err := Exec("insert into price(number, name, amount, number_group) select ?, ?, ?, ?", tx, args...)

	return err
}
