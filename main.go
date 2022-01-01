package main

import (
	"fmt"
	"log"
)

func main() {
	var databaseName string
	var databaseNames []string

	db := connect()

	rows, error := db.Query("SHOW DATABASES")

	if error != nil {
		panic(error)
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&databaseName)

		if err != nil {
			log.Fatal(err)
		}

		databaseNames = append(databaseNames, databaseName)

	}
	fmt.Println(databaseNames)

	err := rows.Err()

	if err != nil {
		log.Fatal(err)
	}
}
