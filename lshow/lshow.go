package lshow

import (
	"fmt"
	"log"
	"mwlog/db"
	"mwlog/ui"
)

const (
	colFreq = iota
	colCall
	colCity
	colState
	colCountry
	colPower
	colPattern
	colClass
)

// LoadLS loads the data into a liststore
 func LoadLS() {

	ls := ui.GetListStore("liststore")

	rows := db.GetAllMWList()
	defer rows.Close()

	for rows.Next() {

		var call string
		var freq int
		var city string
		var state string
		var country string
		var power string
		var pattern string
		var class string
		rows.Scan(&call, &freq, &city, &state, &country, &power, &pattern, &class)
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
		if err := ls.SetValue(i, colState, state); err != nil {
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
