package database

import (
	"database/sql"
	"log"
)

func GetComponents() *sql.Rows {
	db := Connect()

	result, err := db.Query("SELECT * FROM component")
	if err != nil {
		panic(err.Error())
	}

	return result
}

func GetComponentParents(componentId string) *sql.Rows {
	db := Connect()

	query := "SELECT component.id, component.name FROM component JOIN dependency ON dependency.componentA = component.id WHERE componentB = " + componentId

	result, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}

	return result
}

func CreateComponent(componentName string) int64 {
	if componentName != "" {
		db := Connect()

		res, err := db.Exec("INSERT INTO component (`name`) VALUES (?)", componentName)

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
