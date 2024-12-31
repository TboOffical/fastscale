package main

import (
	"database/sql"
	"log"
	"os"
)

func createTables(db *sql.DB) {
	//load the tables file; this is included with the project and it should be there by default

	content, err := os.ReadFile("./tables.sql")
	if err != nil {
		log.Fatalln("Tables file was not found, re-download the application?")
	}

	_, err = db.Exec(string(content))
	if err != nil {
		log.Fatalln("Was not able to create essential database tables: ", err.Error())
	}
}

func getDatabaseVariable(name string, db *sql.DB) (string, error) {
	rows := db.QueryRow("SELECT text FROM operations WHERE variable=?", name)

	var variableContent string
	err := rows.Scan(variableContent)
	if err != nil {
		return "", err
	}

	return variableContent, nil
}

func setDatabaseVariable(name string, content string, db *sql.DB) error {
	_, err := db.Exec("INSERT INTO operations (variable, content) VALUES (?, ?) ON CONFLICT(variable) DO UPDATE SET content=? ", name, content, content)
	if err != nil {
		return err
	}
	return nil
}
