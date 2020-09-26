package lshow

import (
	// "database/sql"
	"fmt"
	"log"
	"mwlog/db"
	"mwlog/ui"

	// "github.com/gotk3/gotk3/gtk"
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

// LoadLS loads the data into a liststore
func LoadLS() {

	ls, err := ui.GetListStore("liststore")
	if err != nil {
		log.Fatalln(err)
	}

	rows := db.GetAllFCC()
	defer rows.Close()

	for rows.Next() {

		var call string
		var freq int
		var city string
		var prov string
		var country string
		var power string
		var pattern string
		var class string
		rows.Scan(&call, &freq, &city, &prov, &country, &power, &pattern, &class)

		i := ls.Append()
		// Set the contents of the tree store row that the iterator represents
		f := fmt.Sprintf("%4d", freq)
		if err := ls.SetValue(i, colFreq, f); err != nil {
			log.Fatal("Unable set value:", err)
		}
		if err := ls.SetValue(i, colCall, call); err != nil {
			log.Fatal("Unable set value:", err)
		}
		if err := ls.SetValue(i, colCity, city); err != nil {
			log.Fatal("Unable set value:", err)
		}
		if err := ls.SetValue(i, colProv, prov); err != nil {
			log.Fatal("Unable set value:", err)
		}
		if err := ls.SetValue(i, colCountry, country); err != nil {
			log.Fatal("Unable set value:", err)
		}
		if err := ls.SetValue(i, colPower, power); err != nil {
			log.Fatal("Unable set value:", err)
		}
		if err := ls.SetValue(i, colPattern, pattern); err != nil {
			log.Fatal("Unable set value:", err)
		}
		if err := ls.SetValue(i, colClass, class); err != nil {
			log.Fatal("Unable set value:", err)
		}
	}
}
