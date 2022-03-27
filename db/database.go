package db

import (
	"database/sql"
	"fmt"
	"strconv"

	// "mwlog/ui"

	"log"
	"os"
)

var sqlDb *sql.DB

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
	Country   string
	Signal    string
	Format    int
	Remarks   string
	Receiver  int
	Antenna   int
	Latitude  float64
	Longitude float64
	Distance  float64
	Bearing   float64
	Sunstatus string
}

// OpenDB opens the logging database, creating it if needed
func OpenDB() {
	var err error
	if _, err = os.Stat("mwlog.db"); err != nil {
		file, err := os.Create("mwlog.db")
		if err != nil {
			log.Fatalln(err.Error())
		}
		err = file.Close()
		if err != nil {
			log.Fatalln(err.Error())
		}
		if sqlDb, err = sql.Open("sqlite3", "mwlog.db"); err != nil {
			log.Fatalln(err.Error())
		}
		if err := createDDL(); err != nil {
			log.Fatalln(err.Error())
		}
		err = sqlDb.Close()
		if err != nil {
			log.Fatalln(err.Error())
		}
	}
	if sqlDb, err = sql.Open("sqlite3", "mwlog.db"); err != nil {
		log.Fatalln(err.Error())
	}
}

// GetAllMWList retrieves all station information into a Rows[] structure
func GetAllMWList() *sql.Rows {
	readSQL := `SELECT frequency, station, city, state, country, power_day, power_night, distance, bearing
	 FROM mwlist ORDER by cast(frequency as number), station;`
	rows, err := sqlDb.Query(readSQL)
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
	rows, err := sqlDb.Query(readSQL)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return rows
}

// GetFormatByID returns an Name for ID
func GetFormatByID(id int) string {
	var value string

	readSQL := "SELECT value FROM formats WHERE ID = ?"

	rows, err := sqlDb.Query(readSQL, id)
	if err != nil {
		log.Println(err.Error())
		return err.Error()
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&value)
		if err != nil {
			log.Println(err.Error())
			return err.Error()
		}
		return value
	}
	log.Println(err.Error())
	return err.Error()
}

// GetFormatIDByName returns an ID for Name
func GetFormatIDByName(name string) int {
	var value string

	readSQL := `SELECT id FROM formats WHERE Name = ?`
	rows, err := sqlDb.Query(readSQL, name)
	if err != nil {
		log.Println(err.Error())
		return -1
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&value)
		if err != nil {
			log.Println(err.Error())
			return -1
		}
		id, _ := strconv.Atoi(value)
		return id
	}
	log.Println(err.Error())
	return -1
}

// GetAllFormats returns a pointer to SQL rows
func GetAllFormats() *sql.Rows {
	readSQL := `SELECT id, name FROM formats ORDER BY name;`
	rows, err := sqlDb.Query(readSQL)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return rows
}

// GetReceiverByID returns an Name for ID
func GetReceiverByID(id int) string {
	var value string

	readSQL := `SELECT name FROM receiver WHERE ID = ?`

	rows, err := sqlDb.Query(readSQL, id)
	if err != nil {
		log.Println(err.Error())
		return err.Error()
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&value)
		if err != nil {
			log.Println(err.Error())
			return err.Error()
		}
		return value
	}
	log.Println(err.Error())
	return err.Error()
}

// GetReceiverIDByName returns an ID for Name
func GetReceiverIDByName(name string) int {
	var value string

	readSQL := `SELECT id FROM receivers WHERE Name = ?`
	rows, err := sqlDb.Query(readSQL, name)
	if err != nil {
		log.Println(err.Error())
		return -1
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&value)
		if err != nil {
			log.Println(err.Error())
			return -1
		}
		id, _ := strconv.Atoi(value)
		return id
	}
	log.Println(err.Error())
	return -1
}

// GetAllReceivers returns a pointer to SQL rows
func GetAllReceivers() *sql.Rows {
	readSQL := `SELECT id, name FROM receivers ORDER BY name;`
	rows, err := sqlDb.Query(readSQL)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return rows
}

// GetAntennaByID returns an Name for ID
func GetAntennaByID(id int) string {
	var value string

	readSQL := `SELECT name FROM antennas WHERE ID = ?`

	rows, err := sqlDb.Query(readSQL, id)
	if err != nil {
		log.Println(err.Error())
		return err.Error()
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&value)
		if err != nil {
			log.Println(err.Error())
			return err.Error()
		}
		return value
	}
	log.Println(err.Error())
	return err.Error()
}

// GetAntennaIDByName returns an ID for Name
func GetAntennaIDByName(name string) int {
	var value string

	readSQL := `SELECT id FROM antennas WHERE Name = ?`
	rows, err := sqlDb.Query(readSQL, name)
	if err != nil {
		log.Println(err.Error())
		return -1
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&value)
		if err != nil {
			log.Println(err.Error())
			return -1
		}
		id, _ := strconv.Atoi(value)
		return id
	}
	log.Println(err.Error())
	return -1
}

// GetAllAntennas returns a pointer to SQL rows
func GetAllAntennas() *sql.Rows {
	readSQL := `SELECT id, name FROM antennas ORDER BY name;`
	rows, err := sqlDb.Query(readSQL)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return rows
}



// AddLogging saves a log entry to the loggings table
func AddLogging(l LogRecord) (int, error) {
	_, err := sqlDb.Exec(`Insert into loggings (date, time, station, frequency, city, state, country, 
		signal, format, remarks, receiver, antenna, latitude, longitude, distance, bearing, sunstatus) 
								 values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
		l.Dt, l.Tm, l.Station, l.Frequency, l.City, l.State, l.Country, l.Signal, l.Format, l.Remarks, l.Receiver, l.Antenna, l.Latitude, l.Longitude, l.Distance, l.Bearing, l.Sunstatus)
	if err != nil {
		return -1, err
	}
	log.Println("Add logging", l.ID)
	var id int
	rows, err := sqlDb.Query("Select last_insert_rowid()")
	if err != nil {
		return -1, err
	}
	defer rows.Close()
	if rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			return -1, err
		}
		return id, nil
	}
	return -1, err
}

// UpdateLogging updates an existing logging record
func UpdateLogging(l *LogRecord) {
	s, err := sqlDb.Prepare(`update loggings 
	set date = ?, time = ?, station = ?, frequency = ?, city = ?, state = ?, 
									country = ?, signal = ?, format = ?, remarks = ?, receiver = ?, antenna = ?,
									latitude = ?, longitude = ?, distance = ?, bearing = ?, sunstatus = ?
									where id = ?`)
	if err != nil {
		log.Println(err.Error())
	}
	defer s.Close()
	_, err = s.Exec(l.Dt, l.Tm, l.Station, l.Frequency, l.City, l.State, l.Country, l.Signal, l.Format, l.Remarks, l.Receiver, l.Antenna,
		l.Latitude, l.Longitude, l.Distance, l.Bearing, l.Sunstatus, l.ID)
	if err != nil {
		log.Println(err.Error())
	}
}

// GetLogBookStore fills the ListStore for the logbook page
func GetLogBookStore() (*sql.Rows, error) {

	rows, err := sqlDb.Query(`select cast(id as text), date, time, station, frequency, city, state, country, signal, remarks 
	from loggings 
	order by date, time`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return rows, nil
}

// GetLoggingForFreq fills the ListStore for the channels logging section
func GetLoggingForFreq(freq string) (*sql.Rows, error) {

	rows, err := sqlDb.Query(`select id, station, city, state, country, format, 
       										min(date) as firstheard, count(*) as times
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
	q := "delete from loggings where id = ?"
	stmt, err := sqlDb.Prepare(q)
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
func GetLoggingByID(id int) (LogRecord, error) {
	var l LogRecord

	q := `select id, date, time, station, frequency, city, state, country, signal, 
					format, remarks, receiver, antenna, 
					latitude, longitude, distance, bearing, sunstatus
		 from loggings where id = ?`

	rows, err := sqlDb.Query(q, id)
	if err != nil {
		log.Println(err.Error())
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&l.ID, &l.Dt, &l.Tm, &l.Station, &l.Frequency, &l.City, &l.State, &l.Country, &l.Signal, &l.Format, &l.Remarks, &l.Receiver, &l.Antenna, &l.Latitude, &l.Longitude, &l.Distance, &l.Bearing, &l.Sunstatus)
		if err != nil {
			return LogRecord{}, err
		}
	} else {
		return LogRecord{}, fmt.Errorf("unable to retrieve logging by id")
	}

	return l, nil
}

// GetLoggingLocations returns a *sql.Rows dataset of station, lat, long
func GetLoggingLocations() *sql.Rows {
	q := `select station, latitude, longitude from loggings`
	rows, err := sqlDb.Query(q)
	if err != nil {
		log.Println(err.Error())
	}
	return rows
}

// GetChannel get the channel info
func GetChannel(freq string) (*Channel, error) {
	var ch Channel

	q := "select id, class, daytime, nighttime from channel where frequency = ?"
	rows, err := sqlDb.Query(q, freq)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var id int
	var class, daytime, nighttime string
	if rows.Next() {
		err := rows.Scan(&id, &class, &daytime, &nighttime)
		if err != nil {
			return nil, err
		}
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
		stmt, err := sqlDb.Prepare(q)
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
		stmt, err := sqlDb.Prepare(q)
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
