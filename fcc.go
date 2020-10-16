package main

import (
	"fmt"
	"mwlog/db"
	"mwlog/ui"
)

func loadFCCData(station string) error {
	var id uint
	var latitude, longitude, distance, bearing float64
	var freq, city, prov, cnty string
	rows := db.GetFCCByCall(station)
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&id, &station, &freq, &city, &prov, &cnty, &latitude, &longitude, &distance, &bearing)
		ui.GetEntry("logging_station").SetText(station)
		ui.GetEntry("logging_frequency").SetText(freq)
		ui.GetEntry("logging_city").SetText(city)
		ui.GetEntry("logging_province").SetText(prov)
		ui.GetEntry("logging_country").SetText(cnty)
		ui.GetLabel("logging_latitude").SetText(fmt.Sprintf("%.3f", latitude))
		ui.GetLabel("logging_longitude").SetText(fmt.Sprintf("%.3f", longitude))
		ui.GetLabel("logging_distance").SetText(fmt.Sprintf("%.0f", distance))
		ui.GetLabel("logging_bearing").SetText(fmt.Sprintf("%.0f", bearing))
	}
	return nil
}
