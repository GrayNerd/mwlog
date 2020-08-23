package main

import (
	"database/sql"
	"log"
	"os"
)

var sqldb *sql.DB

func openDB() {
	var err error
	var c bool = false
	if _, err = os.Stat("mwlog.db"); err != nil {
		file, err := os.Create("mwlog.db")
		if err != nil {
			log.Fatal(err)
		}
		file.Close()
		log.Println("mwlog.db created")
		c = true
	}
	if sqldb, err = sql.Open("sqlite3", "mwlog.db"); err != nil {
		log.Fatal(err)
	}
	if c == true {
		if err := createDDL(); err != nil {
			log.Fatal(err)
		}
	}
}

func createDDL() error {
	if err := createFCCTable(sqldb); err != nil {
		return err
	}
	return nil
}

func createFCCTable(db *sql.DB) error {
	createFCCTableSQL := `CREATE TABLE fcc (
					"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
					"call" TEXT,
					"frequency" INTEGER,
					"city" TEXT,
					"prov" TEXT,
					"country" TEXT,
					"power" TEXT,
					"pattern" TEXT,
					"class" TEXT		
				  );`

	log.Println("Create fcc table...")
	statement, err := db.Prepare(createFCCTableSQL)
	if err != nil {
		return err
	}
	statement.Exec()
	log.Println("fcc table created")
	return nil
}
