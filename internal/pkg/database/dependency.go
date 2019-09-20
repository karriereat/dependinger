package database

import (
	"log"
)

func CreateDependency(componentA int, componentB int) int64 {

	if componentA != 0 && componentB != 0 {
		db := Connect()

		res, err := db.Exec("INSERT INTO dependency (`componentA`, `componentB`) VALUES (?,?)", componentA, componentB)

		if err == nil {
			id, err := res.LastInsertId()
			if err == nil {
				return id
			}
		}
	}

	if err != nil {
		log.Fatal(err.Error())
	}

	return 0

}
