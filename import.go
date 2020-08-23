package main

import (
	// "bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	// "os"
	"strings"

	// "github.com/anaskhan96/soup"
	"github.com/gotk3/gotk3/gtk"
)

const (
	colFreq = iota
	colCall
	colCity
	colProv
	colCountry
	colPower
	colPattern
	colClass
)

func importFCC(win *gtk.Window) {
	log.Println("download started...")
	res, err := http.Get("https://transition.fcc.gov/fcc-bin/amq?call=&arn=&state=&city=&freq=530&fre2=1700&type=0&facid=&class=&list=4&NextTab=Results+to+Next+Page%2FTab&dist=&dlat2=&mlat2=&slat2=&NS=N&dlon2=&mlon2=&slon2=&EW=W&size=9")
	if err != nil {
		log.Fatalln(err)
	}
	robots, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	rows := strings.Split(string(robots), "\n")
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
		call := columns[1]
		if callExists(call) == true {
			pattern := getCurrentPattern(call)
			pattern = updatePattern(pattern, columns[5], columns[6])
			power := getCurrentPower(call)
			power = updatePower(power, columns[14], power)

			s, err := sqldb.Prepare(fmt.Sprintf("update fcc set pattern = ?, power = ? where call = ?"))
			_, err = s.Exec(pattern, power, call)
			if err != nil {
				log.Fatalln(err.Error())
			}
		} else {
			var frequency int
			fmt.Sscanf(columns[2], "%d", &frequency)

			city := strings.Trim(columns[10], " ")
			stateprov := strings.Trim(columns[11], " ")
			country := strings.Trim(columns[12], " ")
			class := strings.Trim(columns[7], " ")

			power, pattern := "", ""
			power = updatePower(power, (columns[14])[0:5], columns[6])
			pattern = updatePattern(pattern, columns[5], columns[6])

			err = insertStation(sqldb, call, frequency, city, stateprov, country, power, pattern, class)
			if err != nil {
				log.Fatalln(err.Error())
			}
			// fmt.Printf("*%v* *%v* *%v* *%v* *%v* *%v* *%v*\n",
			// 	call, frequency, pattern, city, stateprov, country, power)
		}
	}
	dialog := gtk.MessageDialogNew(win, gtk.DIALOG_DESTROY_WITH_PARENT, gtk.MESSAGE_INFO, gtk.BUTTONS_OK, "Loading completed.")
	dialog.SetTitle("Information")
	dialog.Run()
	dialog.Hide()
}

func insertStation(db *sql.DB, callLetters string, frequency int, city string, prov string,
	country string, power string, pattern string, ch string) error {
	log.Printf("Inserting station record ... %v %s\n", frequency, callLetters)
	insertStationSQL := `INSERT INTO fcc(call, frequency, city, prov, country, power, pattern, class) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	statement, err := db.Prepare(insertStationSQL)
	if err != nil {
		return err
	}
	if _, err = statement.Exec(callLetters, frequency, city, prov, country, power, pattern, ch); err != nil {
		return err
	}
	return nil
}

func callExists(call string) bool {
	sql := fmt.Sprintf("select count(*) from fcc where call = '%s'", call)
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

func getCurrentPower(call string) string {
	sql := fmt.Sprintf("select power from fcc where call = '%s'", call)
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
func getCurrentPattern(call string) string {
	sql := fmt.Sprintf("select pattern from fcc where call = '%s'", call)
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

func updatePattern(pattern, newPat string, opTime string) string {
	p := strings.Split(pattern, "/")
	l := len(p)

	switch strings.Trim(opTime, " ") {
	case "Daytime":
		pattern = newPat + "//"
		if l > 1 {
			pattern += p[1]
		} else {
		}
	case "Unlimited":
		pattern = p[0]
	case "Nighttime":
		if l > 0 {
			pattern = p[0] + "//" + newPat
		} else {
			pattern = "//" + newPat
		}
	}
	return pattern
}
func updatePower(power string, newPow string, opTime string) string {
	p := strings.Split(power, "/")
	l := len(p)

	switch strings.Trim(opTime, " ") {
	case "Daytime":
		power = newPow + "/"
		if l > 1 {
			power += p[1]
		}
	case "Unlimited":
		power = newPow + power
	case "Nighttime":
		if l > 0 {
			power = p[0] + "/" + newPow
		} else {
			power = "/" + newPow
		}
	}
	return power
}
