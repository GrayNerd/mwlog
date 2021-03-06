package db

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	"log"
	"strconv"

	"github.com/gotk3/gotk3/gtk"
)

// ImportMWList downloads and imports the mwlist.gov MW station data
func ImportMWList() {

	// dialog := gtk.MessageDialogNew(win, gtk.DIALOG_DESTROY_WITH_PARENT, gtk.MESSAGE_INFO, gtk.BUTTONS_CLOSE, "%s", "")
	// dialog.SetTitle("Importing MWList data")
	// dialog.SetSizeRequest(300, 200)
	// box, _ := dialog.GetMessageArea()
	// btn, _ := dialog.GetWidgetForResponse(gtk.RESPONSE_CLOSE)
	// btn.SetSensitive(false)
	// label, _ := gtk.LabelNew("Downloading MWList data...\n")
	// box.Add(label)
	// btn.Connect("clicked", func() {
	// 	dialog.Destroy()
	// })
	// glib.IdleAdd(dialog.ShowNow)

	// var res *http.Response
	// res, err := http.Get("https://transition.fcc.gov/fcc-bin/amq?call=&arn=&state=&city=&freq=530&fre2=1700&type=0&facid=&class=&list=4&NextTab=Results+to+Next+Page%2FTab&dist=10000&dlat2=50&mlat2=30&slat2=&NS=N&dlon2=104&mlon2=30&slon2=&EW=W&size=9")
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// defer res.Body.Close()
	// robots, err := ioutil.ReadAll(res.Body)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// rows := strings.Split(string(robots), "\n")

	// load from file code

	dialog, _ := gtk.FileChooserDialogNewWith2Buttons("Select Import File",
		nil, gtk.FILE_CHOOSER_ACTION_OPEN,
		"Cancel", gtk.RESPONSE_CANCEL,
		"Open", gtk.RESPONSE_ACCEPT)

	result := dialog.Run()
	if result == gtk.RESPONSE_CANCEL {
		dialog.Close()
		return
	}

	filename := dialog.GetFilename()
	dialog.Close()

	f, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer f.Close()

	createMWListTable()
	truncateTable("mwlist")

	r := csv.NewReader(f)
	r.Comma = ';'
	r.FieldsPerRecord = -1
	startLine := false
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if len(row) == 0 {
			continue
		}
		if err != nil {
			log.Fatalln(err.Error())
		}
		if startLine == false {
			if strings.Contains(row[0], "======") {
				startLine = true
			}
			continue
		}

		freq := row[0]
		country := row[1]
		language := row[2]
		station := row[3]
		address := row[5]
		power := row[9]
		city, state, powerDay, powerNight := parseAddress(address, power)

		latitude, _ := strconv.ParseFloat(row[6], 64)
		longitude, _ := strconv.ParseFloat(row[7], 64)

		distance, _ := strconv.Atoi(row[11])
		bearing, _ := strconv.Atoi(row[12])

		if !(country == "USA" || country == "CAN" || country == "MEX") {
			continue
		}
		f, _ := strconv.Atoi(freq)
		if f < 530 {
			continue
		}
		if f > 1710 {
			break
		}

		if callExists(station) {
			ss := strings.FieldsFunc(row[5], func(r rune) bool {
				if r == '(' || r == ')' {
					return true
				}
				return false
			})
			if len(ss) > 2 {
				if ss[1] == "D" { //set day power
					setPower(station, "power_day", power)
				}
				if ss[1] == "N" { //set night power
					setPower(station, "power_night", power)
				}
			}

		} else {
			err = insertStation(station, freq, city, state, country, language, powerDay, powerNight,
				latitude, longitude, distance, bearing)
			if err != nil {
				log.Println(err.Error())
			}
		}
	}
}

func parseAddress(address, power string) (city, state, powerDay, powerNight string) {
	ss := strings.FieldsFunc(address, func(r rune) bool {
		if r == '(' || r == ')' {
			return true
		}
		return false
	})

	if len(ss) > 1 {
		city = ss[0][:len(ss[0])-1]
		state = ss[len(ss)-1]
	} else {
		city = address
		state = ""
	}
	if len(ss) > 2 { // has day/night
		if ss[1] == "D" {
			powerDay = power
			powerNight = "off"
		} else {
			powerNight = power
		}
	} else {
		powerDay = power
		powerNight = ""
	}
	return
}

func insertStation(station, frequency, city, state, country, language string,
	powerDay, powerNight string, latitude, longitude float64, distance, bearing int) error {

	log.Printf("Inserting station record ... %v %s\n", frequency, station)
	q := `INSERT INTO mwlist(station, frequency, city, state, country, language, power_day, power_night, 
										 latitude, longitude, distance, bearing) 
							VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	stmt, err := sqldb.Prepare(q)
	if err != nil {
		return err
	}
	if _, err = stmt.Exec(station, frequency, city, state, country, language, powerDay, powerNight,
		latitude, longitude, distance, bearing); err != nil {
		return err
	}
	return nil
}

func callExists(station string) bool {
	sql := fmt.Sprintf("select count(*) from mwlist where station = '%s'", station)
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
func setPower(station, field, power string) {
	q := fmt.Sprintf(`update mwlist set %s = "%s" where station = "%s"`, field, power, station)
	stmt, err := sqldb.Prepare(q)
	if err != nil {
		log.Println(err.Error())
		return
	}
	stmt.Exec(q)
}

// func getCurrentPower(station string) string {
// 	sql := fmt.Sprintf("select power from mwlist where station = '%s'", station)
// 	row, err := sqldb.Query(sql)
// 	if err != nil {
// 		log.Println(err.Error())
// 	}
// 	defer row.Close()
// 	var power string
// 	row.Next()
// 	row.Scan(&power)
// }

// func getCurrentPattern(station string) string {
// 	sql := fmt.Sprintf("select pattern from mwlist where station = '%s'", station)
// 	row, err := sqldb.Query(sql)
// 	if err != nil {
// 		log.Println(err.Error())
// 	}
// 	defer row.Close()
// 	var pattern string
// 	row.Next()
// 	row.Scan(&pattern)
// 	return pattern
// }

// func formatPattern(station string, newPat string, opTime string) string {
// 	pattern := getCurrentPattern(station)
// 	p := strings.Split(pattern, "/")
// 	l := len(p)
// 	switch strings.Trim(opTime, " ") {
// 	case "Daytime":
// 		pattern = newPat
// 		if l > 1 {
// 			pattern += "/" + p[1]
// 		} else {
// 		}
// 	case "Unlimited":
// 		pattern = p[0]
// 	case "Nighttime":
// 		if l > 0 {
// 			pattern = p[0] + "/" + newPat
// 		} else {
// 			pattern = "/" + newPat
// 		}
// 	}
// 	return pattern
// }

// func formatPower(station string, newPow string, opTime string) string {
// 	power := getCurrentPower(station)
// 	p := strings.Split(power, "/")
// 	l := len(p)
// 	switch strings.Trim(opTime, " ") {
// 	case "Daytime":
// 		power = newPow
// 		if l > 1 {
// 			power += "/" + p[1]
// 		}
// 	case "Unlimited":
// 		power = newPow + " U"
// 	case "Nighttime":
// 		if l > 0 {
// 			power = p[0] + "/" + newPow
// 		} else {
// 			power = "/" + newPow
// 		}
// 	}
// 	return power
// }
