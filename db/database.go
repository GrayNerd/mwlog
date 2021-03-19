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
// 	colCity
// 	colState
// 	colCountry
// 	colSignal
// 	colRemarks
// 	colRcvr
// 	colAnt
// )

var sqldb *sql.DB

// Channel is the structure of the channels table
type Channel struct {
	ID        int
	Frequency string
	Class     string
	Daytime   string
	Nighttime string
}

// LogRecord is the structure of the logging table
type LogRecord struct {
	ID        int
	Dt        string
	Tm        string
	Station   string
	Frequency string
	City      string
	State     string
	Cnty      string
	Signal    string
	Format    int
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

// GetAllMWList retrieves all station information into a Rows[] structure
func GetAllMWList() *sql.Rows {
	readSQL := `SELECT frequency, station, city, state, country, power_day, power_night, distance, bearing
	 FROM mwlist ORDER by cast(frequency as number), station;`
	rows, err := sqldb.Query(readSQL)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return rows
}

// GetMWListByCall retrieves the mwlist data for a specified station
func GetMWListByCall(station string) *sql.Rows {
	readSQL := fmt.Sprintf(`SELECT id, station, frequency, city, state, country, latitude, longitude, distance, bearing
							  FROM mwlist 
							  WHERE station = upper("%v");`, station)
	rows, err := sqldb.Query(readSQL)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return rows
}

// AddLogging saves a log entry to the loggings table
func AddLogging(l *LogRecord) int {
	_, err := sqldb.Exec(`Insert into loggings (date, time, station, frequency, city, state, country, 
		signal, format, remarks, receiver, antenna, latitude, longitude, distance, bearing, sunrise, sunset) 
								 values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
								 l.Dt, l.Tm, l.Station, l.Frequency, l.City, l.State, l.Cnty, l.Signal, l.Format, l.Remarks, l.Rcvr, l.Ant, l.Latitude, l.Longitude, l.Distance, l.Bearing, l.Sunrise, l.Sunset)
	if err != nil {
		log.Println(err.Error())
	}
	log.Println("Add logging", l.ID)
	var id int
	rows, err := sqldb.Query("Select last_insert_rowid()")
	if err != nil {
		log.Println(err.Error())
	}
	defer rows.Close()
	rows.Next()
	rows.Scan(&id)
	l.ID = id
	return id
}

// UpdateLogging updates an existing logging record
func UpdateLogging(l *LogRecord) {
	s, err := sqldb.Prepare(`update loggings 
	set date = ?, time = ?, station = ?, frequency = ?, city = ?, state = ?, 
									country = ?, signal = ?, format = ?, remarks = ?, receiver = ?, antenna = ?,
									latitude = ?, longitude = ?, distance = ?, bearing = ?, sunrise = ?, sunset = ?
									where id = ?`)
	if err != nil {
		log.Println(err.Error())
	}
	defer s.Close()
	_, err = s.Exec(l.Dt, l.Tm, l.Station, l.Frequency, l.City, l.State, l.Cnty, l.Signal, l.Format, l.Remarks, l.Rcvr, l.Ant,
		l.Latitude, l.Longitude, l.Distance, l.Bearing, l.Sunrise, l.Sunset, l.ID)
	if err != nil {
		log.Println(err.Error())
	}
}

// GetLogBookStore fills the liststore for the logbook page
func GetLogBookStore() (*sql.Rows, error) {

	rows, err := sqldb.Query(`select cast(id as text), date, time, station, frequency, city, state, country, signal, remarks 
	from loggings 
	order by date, time`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return rows, nil
}

// GetLoggingForFreq fills the liststore for the channels logging section
func GetLoggingForFreq(freq string) (*sql.Rows, error) {

	rows, err := sqldb.Query(`select id, station, city, state, country, format, min(date) as firstheard, count(*) as times
								from loggings 
								where frequency = ?
								group by station
								order by station desc`, freq)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return rows, nil
}

// DeleteLogging delete an entry in the logging table specified ID
func DeleteLogging(id int) {
	q := "delete from loggings where id = '?'"
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
func GetLoggingByID(id int) (*LogRecord, error) {
	var l LogRecord

	q := `select id, date, time, station, frequency, city, state, country, signal, format, remarks, 
			receiver, antenna, latitude, longitude, distance, bearing, sunrise, sunset
		 from loggings where id = ?`
	rows, err := sqldb.Query(q, id)
	if err != nil {
		log.Println(err.Error())
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&l.ID, &l.Dt, &l.Tm, &l.Station, &l.Frequency, &l.City, &l.State, &l.Cnty, &l.Signal,
				  &l.Format, &l.Remarks, &l.Rcvr, &l.Ant,
				&l.Latitude, &l.Longitude, &l.Distance, &l.Bearing, &l.Sunrise, &l.Sunset)
	} else {
		return nil, fmt.Errorf("Unable to retrieve logging by id")
	}

	return &l, nil

}

// GetLoggingLocations returns a *sql.Rows dataset of station, lat, long
func GetLoggingLocations() *sql.Rows {
	q := `select station, latitude, longitude from loggings`
	rows, err := sqldb.Query(q)
	if err != nil {
		log.Println(err.Error())
	}
	return rows
}

// GetChannel get the channel info
func GetChannel(freq string) (*Channel, error) {
	var ch Channel

	q := "select id, class, daytime, nighttime from channel where frequency = ?"
	rows, err := sqldb.Query(q, freq)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var id int
	var class, daytime, nighttime string
	if rows.Next() {
		rows.Scan(&id, &class, &daytime, &nighttime)
		ch.ID = id
		ch.Frequency = freq
		ch.Class = class
		ch.Daytime = daytime
		ch.Nighttime = nighttime
	} else { // must be new entry
		return nil, fmt.Errorf("channel entry not found: %v", freq)
	}
	return &ch, nil
}

// SaveChannel saves the channel info
func SaveChannel(ch *Channel) error {
	if ch.ID < 1 {
		q := `insert into channel (frequency, class, daytime, nighttime) 
				values(?,?,?,?)`
		stmt, err := sqldb.Prepare(q)
		if err != nil {
			return err
		}
		defer stmt.Close()
		_, err = stmt.Exec(ch.Frequency, ch.Class, ch.Daytime, ch.Nighttime)
		if err != nil {
			return err
		}
	} else {
		q := "update channel set frequency = ?, class = ?, daytime = ?, nighttime = ? where id = ?"
		stmt, err := sqldb.Prepare(q)
		if err != nil {
			return err
		}
		defer stmt.Close()
		_, err = stmt.Exec(ch.Frequency, ch.Class, ch.Daytime, ch.Nighttime, ch.ID)
		if err != nil {
			return err
		}
	}
	return nil
}
