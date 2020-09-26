package main

import (
	"mwlog/db"
	"mwlog/ui"
)

func loadFCCData(callsign string) error {
	rows := db.GetFCCByCall(callsign)
	var id int
	var freq, city, prov, cnty string
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&id, &callsign, &freq, &city, &prov, &cnty)

		if c, err := ui.GetEntry("lg_callsign"); err == nil {
			c.SetText(callsign)
		}
		if f, err := ui.GetEntry("lg_frequency"); err == nil {
			f.SetText(freq)
		}
		if c, err := ui.GetEntry("lg_city"); err == nil {
			c.SetText(city)
		}
		if p, err := ui.GetEntry("lg_province"); err == nil {
			p.SetText(prov)
		}
		if c, err := ui.GetEntry("lg_country"); err == nil {
			c.SetText(cnty)
		}


	}
	return nil
}
