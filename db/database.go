package db

import (
	"database/sql"
	"fmt"
	"mwlog/ui"

	"log"
	"os"
)


// const (
// 	colID = iota 
// 	colDate
// 	colTime
// 	colCallsign
// 	colFrequency
// 	colCitt
// 	colProv
// 	colCountry
// 	colSignal 
// 	colProgramming
// 	colRcvr
// 	colAnt
// )

var sqldb *sql.DB

// LogEntry is the structure of the logging table
type LogEntry struct {
	ID          int
	Dt          string
	Tm          string
	Callsign    string
	Frequency   string
	City        string
	Prov        string
	Cnty        string
	Signal      string
	Programming string
	Rcvr        int
	Ant         int
}

// OpenDB opens the logging database, creating it if needed
func OpenDB() {
	var err error
	if _, err = os.Stat("mwlog.db"); err != nil {
		file, err := os.Create("mwlog.db")
		if err != nil {
			log.Fatalln(err.Error())
		}
		file.Close()
	}
	if sqldb, err = sql.Open("sqlite3", "mwlog.db"); err != nil {
		log.Fatalln(err.Error())
	}
	if err := createDDL(); err != nil {
		log.Fatalln(err.Error())
	}
}

func createDDL() error {
	if err := createFCCTable(sqldb); err != nil {
		return err
	}
	return nil
}

func createFCCTable(db *sql.DB) error {
	createFCCTableSQL := `CREATE TABLE IF NOT EXISTS fcc (
							"id"	integer NOT NULL PRIMARY KEY AUTOINCREMENT,
							"callsign"	TEXT,
							"frequency"	TEXT,
							"city"	TEXT,
							"prov"	TEXT,
							"country"	TEXT,
							"power"	TEXT,
							"pattern"	TEXT,
							"class"	TEXT
					  	);`

	stmt, err := db.Prepare(createFCCTableSQL)
	if err != nil {
		return err
	}
	defer stmt.Close()
	if _, err = stmt.Exec(); err != nil {
		return err
	}
	return nil
}

func createLoggingTable(db *sql.DB) error {
	createFCCTableSQL := `CREATE TABLE IF NOT EXISTS "loggings" (
				"ID"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
				"date"	TEXT NOT NULL,
				"time"	TEXT NOT NULL,
				"callsign"	TEXT NOT NULL,
				"frequency"	TEXT NOT NULL,
				"city"	TEXT NOT NULL,
				"province"	TEXT,
				"country"	TEXT NOT NULL,
				"signal"	TEXT NOT NULL,
				"programming"	TEXT NOT NULL,
				"receiver"	INTEGER NOT NULL,
				"antenna"	INTEGER NOT NULL
			)`

	log.Println("Create fcc table...")
	stmt, err := db.Prepare(createFCCTableSQL)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	return nil
}

// GetAllFCC retrieves all station information into a Rows[] structure
func GetAllFCC() *sql.Rows {
	readSQL := `SELECT callsign, frequency, city, prov, country, power, pattern, class
	 FROM fcc ORDER by frequency, callsign;`
	rows, err := sqldb.Query(readSQL)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer rows.Close()
	return rows
}

// GetFCCByCall retrieves the fcc data for a specified callsign
func GetFCCByCall(callsign string) *sql.Rows {
	readSQL := fmt.Sprintf(`SELECT id, callsign, frequency, city, prov, country
	 FROM fcc WHERE callsign = upper("%v");`, callsign)
	rows, err := sqldb.Query(readSQL)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return rows
}

// AddLogging saves a log entry to the loggings table
func AddLogging(l LogEntry) {
	s, err := sqldb.Prepare(`Insert into loggings (dt, tm, callsign, frequency, city, province, 
								 country, signal, programming, receiver, antenna) 
							values(?,?,?,?,?,?,?,?,?,?,?)`)
	if err != nil {
		log.Println(err.Error())
	}
	defer s.Close()
	_, err = s.Exec(l.Dt, l.Tm, l.Callsign, l.Frequency, l.City, l.Prov, l.Cnty, l.Signal, l.Programming, l.Rcvr, l.Ant)
	if err != nil {
		log.Println(err.Error())
	}
}

// UpdateLogging updates an existing logging record
func UpdateLogging(l LogEntry) {
	s, err := sqldb.Prepare(`update loggings 
								set dt = ?, tm = ?, callsign = ?, frequency = ?, city = ?, province = ?, 
									country = ?, signal = ?, programming = ?, receiver = ?, antenna = ?
								where id = ?`)
	if err != nil {
		log.Println(err.Error())
	}
	defer s.Close()
	_, err = s.Exec(l.Dt, l.Tm, l.Callsign, l.Frequency, l.City, l.Prov, l.Cnty, l.Signal, l.Programming, l.Rcvr, l.Ant, l.ID)
	if err != nil {
		log.Println(err.Error())
	}
}

// FillLogBookStore fills the liststore for the logbook page
func FillLogBookStore() {
	ls, err := ui.GetListStore("logbook_store")
	if err != nil {
		return
	}
	ls.Clear()

	rows, err := sqldb.Query(`select id, dt, tm, callsign, frequency, city, province, country, signal, programming 
								from loggings 
								order by dt, tm`)
	if err != nil {
		log.Println(err.Error())
	}
	defer rows .Close()
	var id uint
	var dt, tm, callsign, frequency, city, province, country, signal, programming string
	for rows.Next() {
		rows.Scan(&id, &dt, &tm, &callsign, &frequency, &city, &province, &country, &signal, &programming)
		log.Println(id, dt, tm, callsign)
		iter := ls.Append()
		if err = ls.SetValue(iter, 0, id); err != nil {
			log.Println(err.Error())
		}
		if err = ls.SetValue(iter, 1, dt); err != nil {
			log.Println(err.Error())
		}
		if err = ls.SetValue(iter, 2, tm); err != nil {
			log.Println(err.Error())
		}
		if err = ls.SetValue(iter, 3, callsign); err != nil {
			log.Println(err.Error())
		}
		if err = ls.SetValue(iter, 4, frequency); err != nil {
			log.Println(err.Error())
		}
		l := city + ", " + province + "  " + country
		if err = ls.SetValue(iter, 5, l); err != nil {
			log.Println(err.Error())
		}
		if err = ls.SetValue(iter, 6, signal); err != nil {
			log.Println(err.Error())
		}
		if err = ls.SetValue(iter, 7, programming); err != nil {
			log.Println(err.Error())
		}
	}

}

// DeleteLogging delete an entry in the logging table specified ID
func DeleteLogging(id uint) {
	q := "delete from loggings where id = ?"
	stmt, err := sqldb.Prepare(q)
	if err != nil {
		log.Println(err.Error())
	}
	defer stmt.Close()
	_, err = stmt.Exec(id)
	if err != nil {
		log.Println(err.Error())
	}
}

// GetLoggingByID retrieves a logging by ID
func GetLoggingByID(id uint) (*LogEntry, error) {
	var l LogEntry

	q := "select * from loggings where id = ?"
	rows, err := sqldb.Query(q, id)
	if err != nil {
		log.Println(err.Error())
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&l.ID, &l.Dt, &l.Tm, &l.Callsign, &l.Frequency, &l.City, &l.Prov, &l.Cnty, &l.Signal, &l.Programming, &l.Rcvr, &l.Ant)
	} else {
		return nil, fmt.Errorf("Unable to retrieve logging by id")
	}
	
	return &l, nil

}