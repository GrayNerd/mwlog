package db

import "fmt"

func createDDL() error {
	if err := createMWListTable(); err != nil {
		return err
	}
	if err := createLoggingTable(); err != nil {
		return err
	}
	if err := createAudioTable(); err != nil {
		return err
	}
	if err := createChannelsTable(); err != nil {
		return err
	}
	if err := createMWListTable(); err != nil {
		return err
	}
	return nil
}

// CreateMWListTable creates the mwlist table in the database in it doesn't already exist
func createMWListTable() error {
	s := `CREATE TABLE IF NOT EXISTS "mwlist" ( 
			"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT, 
			"station" TEXT, 
			"frequency" TEXT, 
			"city" TEXT, 
			"state" TEXT, 
			"country" TEXT, 
			"language" TEXT, 
			"power_day" TEXT, 
			"power_night" TEXT, 
			"latitude" REAL,
			"longitude" REAL,
			"distance" NUMERIC,
			"bearing" NUMERIC );`
	err := createTable(s)
	if err != nil {
		return err
	}
	return nil
}

func createLoggingTable() error {
	s := `CREATE TABLE IF NOT EXISTS "loggings" ( 
					"id" INTEGER NOT NULL UNIQUE, 
					"date" TEXT NOT NULL, 
					"time" TEXT NOT NULL, 
					"station" TEXT NOT NULL, 
					"frequency" TEXT NOT NULL, 
					"city" TEXT NOT NULL, 
					"state" TEXT NOT NULL, 
					"country" TEXT NOT NULL, 
					"signal" TEXT NOT NULL, 
					"format" INTEGER, 
					"remarks" BLOB NOT NULL, 
					"receiver" INTEGER NOT NULL, 
					"antenna" INTEGER NOT NULL, 
					"latitude" REAL, 
					"longitude" REAL, 
					"distance" REAL, 
					"bearing" REAL, 
					"sunstatus" TEXT, 
					PRIMARY KEY("ID") 
				);`

	err := createTable(s)
	return err
}

func createAudioTable() error {
	s := `CREATE TABLE IF NOT EXISTS "audio" (
				"id"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
				"sound"	BLOB NOT NULL
				)`
	err := createTable(s)
	return err
}
// createChannelsTable creates the Channels table
func createChannelsTable() error {
	s := `CREATE TABLE IF NOT EXISTS "channels" (
			"id" 		INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
			"frequency"	TEXT NOT NULL UNIQUE,
			"class"		TEXT,
			"daytime"	TEXT,
			"nighttime"	TEXT
			)`
	err := createTable(s)
	return err
}

// createReceiversTable creates the Receivers table
func createReceiversTable() error {
	s := `CREATE TABLE IF NOT EXISTS "receivers" (
					"ID"	INTEGER NOT NULL UNIQUE,
					"Name"	TEXT NOT NULL UNIQUE,	PRIMARY KEY("ID" AUTOINCREMENT)
				)`
	err := createTable(s)
	return err
}

// createAntennasTable creates the Antennas table
func createAntennasTable() error {
	s := `CREATE TABLE IN NOT EXITST "antennas" (
					"ID"	INTEGER NOT NULL UNIQUE,
					"Name"	TEXT NOT NULL UNIQUE,	PRIMARY KEY("ID" AUTOINCREMENT)
				)`
	err := createTable(s)
	return err
}

//createTable create the table specified in the ddl
func createTable(ddl string) error {
	stmt, err := sqlDb.Prepare(ddl)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	return err
}

// truncateTable truncates the named table
func truncateTable(t string) error {
	stmt, err := sqlDb.Prepare(fmt.Sprintf("delete from %s", t))
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	return err
}
