package db

import (
	"database/sql"
	"fmt"
	"net/http"
	// "sync"

	"io/ioutil"
	"log"
	"strings"

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
	res, err := http.Get("https://transition.fcc.gov/fcc-bin/amq?call=&arn=&state=&city=&freq=530&fre2=1700&type=0&facid=&class=&list=4&NextTab=Results+to+Next+Page%2FTab&dist=&dlat2=&mlat2=&slat2=&NS=N&dlon2=&mlon2=&slon2=&EW=W&size=9")
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
		callsign := strings.Trim(columns[1], " ")
		if callExists(callsign) == true {
			pattern := formatPattern(callsign, columns[5], columns[6])
			power := formatPower(callsign, strings.Trim((columns[14])[0:5], " "), columns[6])

			s, err := sqldb.Prepare(fmt.Sprintf("update fcc set pattern = ?, power = ? where callsign = ?"))
			_, err = s.Exec(pattern, power, callsign)
			if err != nil {
				log.Fatalln(err.Error())
			}
		} else {
			// var frequency int
			// fmt.Sscanf(columns[2], "%d", &frequency)
			freq := strings.Trim(columns[2][0:5], " ")
			city := strings.Trim(columns[10], " ")
			stateprov := strings.Trim(columns[11], " ")
			country := strings.Trim(columns[12], " ")
			class := strings.Trim(columns[7], " ")

			var power, pattern string
			power = formatPower(power, strings.Trim((columns[14])[0:5], " "), columns[6])
			pattern = formatPattern(pattern, strings.Trim(columns[5], " "), columns[6])

			err = insertStation(sqldb, callsign, freq, city, stateprov, country, power, pattern, class)
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

func insertStation(db *sql.DB, callsign string, frequency string, city string, prov string,
	country string, power string, pattern string, ch string) error {
	log.Printf("Inserting station record ... %v %s\n", frequency, callsign)
	insertStationSQL := `INSERT INTO fcc(callsign, frequency, city, prov, country, power, pattern, class) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	statement, err := db.Prepare(insertStationSQL)
	if err != nil {
		return err
	}
	if _, err = statement.Exec(callsign, frequency, city, prov, country, power, pattern, ch); err != nil {
		return err
	}
	return nil
}

func callExists(callsign string) bool {
	sql := fmt.Sprintf("select count(*) from fcc where callsign = '%s'", callsign)
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

func getCurrentPower(callsign string) string {
	sql := fmt.Sprintf("select power from fcc where callsign = '%s'", callsign)
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
func getCurrentPattern(callsign string) string {
	sql := fmt.Sprintf("select pattern from fcc where callsign = '%s'", callsign)
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

func formatPattern(callsign string, newPat string, opTime string) string {
	pattern := getCurrentPattern(callsign)

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
func formatPower(callsign string, newPow string, opTime string) string {

	power := getCurrentPower(callsign)

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
