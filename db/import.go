package db

import (
	"database/sql"
	"fmt"
	"net/http"

	// "sync"

	"io/ioutil"
	"log"
	"strings"

	"strconv"
	// "github.com/gotk3/gotk3/glib"
	// "github.com/gotk3/gotk3/gtk"
)

// ImportFCC downloads and imports the fcc.gov MW station data
func ImportFCC() {

	// dialog := gtk.MessageDialogNew(win, gtk.DIALOG_DESTROY_WITH_PARENT, gtk.MESSAGE_INFO, gtk.BUTTONS_CLOSE, "%s", "")
	// dialog.SetTitle("Importing FCC data")
	// dialog.SetSizeRequest(300, 200)
	// box, _ := dialog.GetMessageArea()
	// btn, _ := dialog.GetWidgetForResponse(gtk.RESPONSE_CLOSE)
	// btn.SetSensitive(false)
	// label, _ := gtk.LabelNew("Downloading FCC data...\n")
	// box.Add(label)
	// btn.Connect("clicked", func() {
	// 	dialog.Destroy()
	// })
	// glib.IdleAdd(dialog.ShowNow)

	var res *http.Response
	res, err := http.Get("https://transition.fcc.gov/fcc-bin/amq?call=&arn=&state=&city=&freq=530&fre2=1700&type=0&facid=&class=&list=4&NextTab=Results+to+Next+Page%2FTab&dist=10000&dlat2=50&mlat2=30&slat2=&NS=N&dlon2=104&mlon2=30&slon2=&EW=W&size=9")
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()
	robots, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	// txt, _ := label.GetText()
	// label.SetText(txt + "Importing FCC Data...\n")
	// glib.IdleAdd(dialog.ShowNow)
	rows := strings.Split(string(robots), "\n")
	// load from file code
	// fmt.Printf("%s", string(robots))
	// f, err := os.Open("AM Query Results -- Audio Division (FCC) USA.txt")
	// if err != nil {
	// 	log.Fatalln(err.Error())
	// }
	// defer f.Close()
	// buf := bytes.NewBuffer(make([]byte, 0))
	// if _, err = buf.ReadFrom(f); err != nil {
	// 	log.Fatalln(err.Error())
	// }
	// doc := string(buf.Bytes())
	// if err != nil {
	// 	log.Fatalln(err.Error())
	// }
	// _ = doc
	// links := doc.Find("pre", "class", "listtext").FindAll("span")
	for _, row := range rows {
		if len(row) < 2 {
			break
		}
		columns := strings.Split(row, "|")
		station := strings.Trim(columns[1], " ")
		if callExists(station) == true {
			pattern := formatPattern(station, columns[5], columns[6])
			power := formatPower(station, strings.Trim((columns[14])[0:5], " "), columns[6])

			s, err := sqldb.Prepare(fmt.Sprintf("update fcc set pattern = ?, power = ? where station = ?"))
			_, err = s.Exec(pattern, power, station)
			if err != nil {
				log.Fatalln(err.Error())
			}
		} else {
			freq := strings.Trim(columns[2][0:5], " ")
			city := strings.Trim(columns[10], " ")
			stateprov := strings.Trim(columns[11], " ")
			country := strings.Trim(columns[12], " ")
			class := strings.Trim(columns[7], " ")

			latHour, _ := strconv.Atoi(strings.Trim(columns[20], " "))
			latMin, _ := strconv.Atoi(strings.Trim(columns[21], " "))
			latSec, _ := strconv.Atoi(columns[22][0:2])
			longHour, _ := strconv.Atoi(strings.Trim(columns[24], " "))
			longMin, _ := strconv.Atoi(strings.Trim(columns[25], " "))
			longSec, _ := strconv.Atoi(columns[26][0:2])

			latitude := float64(latHour) + (float64(latMin)/60.0) + (float64(latSec)/3600.0)
			longitude := (float64(longHour) + (float64(longMin)/60.0) + (float64(longSec)/3600.0)) * -1.0

			distance, _ := strconv.ParseFloat(strings.Trim(columns[28], " km"), 64)
			bearing, _ := strconv.ParseFloat(strings.Trim(columns[30], " deg"), 64)
			
			var power, pattern string
			power = formatPower(power, strings.Trim((columns[14])[0:5], " "), columns[6])
			pattern = formatPattern(pattern, strings.Trim(columns[5], " "), columns[6])

			err = insertStation(sqldb, station, freq, city, stateprov, country, power, pattern, class,
				latitude, longitude, distance, bearing)
			if err != nil {
				log.Println(err)
			}
		}
	}
	// txt, _ = label.GetText()
	// label.SetText(txt + "Completed.")
	// btn.SetSensitive(true)
	// glib.IdleAdd(dialog.ShowNow)
	//	dialog.Destroy()
}

func insertStation(db *sql.DB, station string, frequency string, city string, prov string,
	country string, power string, pattern string, ch string,
	latitude, longitude, distance, bearing float64) error {

	log.Printf("Inserting station record ... %v %s\n", frequency, station)
	insertStationSQL := `INSERT INTO fcc(station, frequency, city, prov, country, power, pattern, class, 
										 latitude, longitude, distance, bearing) 
							VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	statement, err := db.Prepare(insertStationSQL)
	if err != nil {
		return err
	}
	if _, err = statement.Exec(station, frequency, city, prov, country, power, pattern, ch,
		latitude, longitude, distance, bearing); err != nil {
		return err
	}
	return nil
}

func callExists(station string) bool {
	sql := fmt.Sprintf("select count(*) from fcc where station = '%s'", station)
	row, err := sqldb.Query(sql)
	if err != nil {
		log.Println(err.Error())
	}
	defer row.Close()

	var val int
	row.Next()
	row.Scan(&val)

	if val != 0 {
		return true
	}
	return false
}

func getCurrentPower(station string) string {
	sql := fmt.Sprintf("select power from fcc where station = '%s'", station)
	row, err := sqldb.Query(sql)
	if err != nil {
		log.Println(err.Error())
	}
	defer row.Close()

	var power string
	row.Next()
	row.Scan(&power)

	return power
}
func getCurrentPattern(station string) string {
	sql := fmt.Sprintf("select pattern from fcc where station = '%s'", station)
	row, err := sqldb.Query(sql)
	if err != nil {
		log.Println(err.Error())
	}
	defer row.Close()

	var pattern string
	row.Next()
	row.Scan(&pattern)

	return pattern
}

func formatPattern(station string, newPat string, opTime string) string {
	pattern := getCurrentPattern(station)

	p := strings.Split(pattern, "/")
	l := len(p)

	switch strings.Trim(opTime, " ") {
	case "Daytime":
		pattern = newPat
		if l > 1 {
			pattern += "/" + p[1]
		} else {
		}
	case "Unlimited":
		pattern = p[0]
	case "Nighttime":
		if l > 0 {
			pattern = p[0] + "/" + newPat
		} else {
			pattern = "/" + newPat
		}
	}
	return pattern
}
func formatPower(station string, newPow string, opTime string) string {

	power := getCurrentPower(station)

	p := strings.Split(power, "/")
	l := len(p)

	switch strings.Trim(opTime, " ") {
	case "Daytime":
		power = newPow
		if l > 1 {
			power += "/" + p[1]
		}
	case "Unlimited":
		power = newPow + " U"
	case "Nighttime":
		if l > 0 {
			power = p[0] + "/" + newPow
		} else {
			power = "/" + newPow
		}
	}
	return power
}
