package db

func createDDL() error {
	if err := createFCCTable(); err != nil {
		return err
	}
	if err := createLoggingTable(); err != nil {
		return err
	}
	if err := createAudioTable(); err != nil {
		return err
	}
	if err := createChannelTable(); err != nil {
		return err
	}
	return nil
}

func createFCCTable() error {
	s := `CREATE TABLE IF NOT EXISTS fcc (
			"id"		integer NOT NULL PRIMARY KEY AUTOINCREMENT,
			"station"	TEXT,
			"frequency"	TEXT,
			"city"		TEXT,
			"prov"		TEXT,
			"country"	TEXT,
			"power"		TEXT,
			"pattern"	TEXT,
			"class"		TEXT,
			"latitude"  REAL,
			"longitude" REAL,
			"distance"  REAL,
			"bearing"   REAL
			);`
	err := createTable(s)
	if err != nil {
		return err
	}
	return nil
}

func createLoggingTable() error {
	s := `CREATE TABLE IF NOT EXISTS "loggings" (
			"ID"		INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
			"date"		TEXT NOT NULL,
			"time"		TEXT NOT NULL,
			"station"	TEXT NOT NULL,
			"frequency"	TEXT NOT NULL,
			"city"		TEXT NOT NULL,
			"province"	TEXT,
			"country"	TEXT NOT NULL,
			"signal"	TEXT NOT NULL,
			"remarks"	TEXT NOT NULL,
			"receiver"	INTEGER NOT NULL,
			"antenna"	INTEGER NOT NULL
			)`

	err := createTable(s)
	if err != nil {
		return err
	}
	return nil
}

func createAudioTable() error {
	s := `CREATE TABLE IF NOT EXISTS "audio" (
				"id"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
				"sound"	BLOB NOT NULL
				)`
	err := createTable(s)
	if err != nil {
		return err
	}
	return nil
}
func createChannelTable() error {
	s := `CREATE TABLE IF NOT EXISTS "channel" (
			"id" 		INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
			"frequency"	TEXT NOT NULL UNIQUE,
			"class"		TEXT,
			"daytime"	TEXT,
			"nighttime"	TEXT
			)`
	err := createTable(s)
	if err != nil {
		return err
	}
	return nil
}

func createTable(ddl string) error {
	stmt, err := sqldb.Prepare(ddl)
	if err != nil {
		return err
	}
	defer stmt.Close()
	if _, err = stmt.Exec(); err != nil {
		return err
	}
	return nil
}
