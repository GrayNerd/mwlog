package db

import (
	"database/sql"
	"fmt"
	// "mwlog/ui"

	"log"
	"os"
)

// const (
// 	colID = iota
// 	colDate
// 	colTime
// 	colStation
// 	colFrequency
// 	colCitt
// 	colProv
// 	colCountry
// 	colSignal
// 	colRemarks
// 	colRcvr
// 	colAnt
// )

var sqldb *sql.DB

// LogEntry is the structure of the logging table
type LogEntry struct {
	ID        int
	Dt        string
	Tm        string
	Station   string
	Frequency string
	City      string
	Prov      string
	Cnty      string
	Signal    string
	Remarks   string
	Rcvr      int
	Ant       int
	Latitude  float64
	Longitude float64
	Distance  float64
	Bearing   float64
	Sunrise   string
	Sunset    string
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
		if sqldb, err = sql.Open("sqlite3", "mwlog.db"); err != nil {
			log.Fatalln(err.Error())
		}
		if err := createDDL(); err != nil {
			log.Fatalln(err.Error())
		}
		sqldb.Close()
	}
	if sqldb, err = sql.Open("sqlite3", "mwlog.db"); err != nil {
		log.Fatalln(err.Error())
	}
}

// GetAllFCC retrieves all station information into a Rows[] structure
func GetAllFCC() *sql.Rows {
	readSQL := `SELECT station, frequency, city, prov, country, power, pattern, class
	 FROM fcc ORDER by frequency, station;`
	rows, err := sqldb.Query(readSQL)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return rows
}

// GetFCCByCall retrieves the fcc data for a specified station
func GetFCCByCall(station string) *sql.Rows {
	readSQL := fmt.Sprintf(`SELECT id, station, frequency, city, prov, country, 
								   latitude, longitude, distance, bearing
							  FROM fcc 
							  WHERE station = upper("%v");`, station)
	rows, err := sqldb.Query(readSQL)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return rows
}

// AddLogging saves a log entry to the loggings table
func AddLogging(l LogEntry) int {
	s, err := sqldb.Prepare(`Insert into loggings (date, time, station, frequency, city, province, 
								 country, signal, remarks, receiver, antenna, latitude, longitude, distance, bearing, sunrise, sunset) 
								 values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`)
	if err != nil {
		log.Println(err.Error())
	}
	defer s.Close()
	_, err = s.Exec(l.Dt, l.Tm, l.Station, l.Frequency, l.City, l.Prov, l.Cnty, l.Signal, l.Remarks, l.Rcvr, l.Ant,
					l.Latitude, l.Longitude, l.Distance, l.Bearing, l.Sunrise, l.Sunset)
	if err != nil {
		log.Println(err.Error())
	}
	var id int
	rows, err := sqldb.Query("Select last_insert_rowid()")
	if err != nil {
		log.Println(err.Error())
	}
	defer rows.Close()
	rows.Next()
	rows.Scan(&id)
	return id
}

// UpdateLogging updates an existing logging record
func UpdateLogging(l LogEntry) {
	s, err := sqldb.Prepare(`update loggings 
	set date = ?, time = ?, station = ?, frequency = ?, city = ?, province = ?, 
									country = ?, signal = ?, remarks = ?, receiver = ?, antenna = ?,
									latitude = ?, longitude = ?, distance = ?, bearing = ?, sunrise = ?, sunset = ?
									where id = ?`)
	if err != nil {
		log.Println(err.Error())
	}
	defer s.Close()
	_, err = s.Exec(l.Dt, l.Tm, l.Station, l.Frequency, l.City, l.Prov, l.Cnty, l.Signal, l.Remarks, l.Rcvr, l.Ant, 
		l.Latitude, l.Longitude, l.Distance, l.Bearing, l.Sunrise, l.Sunset, l.ID)
	if err != nil {
		log.Println(err.Error())
	}
}

// GetLogBookStore fills the liststore for the logbook page
func GetLogBookStore() (*sql.Rows, error) {

	rows, err := sqldb.Query(`select id, date, time, station, frequency, city, province, country, signal, remarks 
	from loggings 
	order by date, time`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return rows, nil
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
		rows.Scan(&l.ID, &l.Dt, &l.Tm, &l.Station, &l.Frequency, &l.City, &l.Prov, &l.Cnty, &l.Signal, &l.Remarks, &l.Rcvr, &l.Ant, 
			&l.Latitude, &l.Longitude, &l.Distance, &l.Bearing, &l.Sunrise, &l.Sunset)
	} else {
		return nil, fmt.Errorf("Unable to retrieve logging by id")
	}

	return &l, nil

}
