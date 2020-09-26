package main

import (
	// "fmt"
	"log"
	"time"

	"mwlog/db"
	"mwlog/ui"
	// "github.com/gotk3/gotk3/gtk"
)

var oneAndDone bool = false

func notebookSwitcher(pn int) {
	switch pn {
	case 1:
		initLogEntry()
	}
}
func initLogEntry() {
	if callsign, err := ui.GetEntry("lg_callsign"); err == nil {
		if s, _ := callsign.GetText(); len(s) == 0 {
			clearLogEntry()
			prefillLogEntry()
		}
	}
}

func prefillLogEntry() {
	dt, err := ui.GetEntry("lg_date")
	if err != nil {
		log.Println(err)
	}
	tm, err := ui.GetEntry("lg_time")
	if err != nil {
		log.Println(err)
	}

	currentTime := time.Now()
	dt.SetText(currentTime.Format("2006-01-02"))
	tm.SetText(currentTime.Format("1504"))
}

func clearLogEntry() {
	if dt, err := ui.GetEntry("lg_date"); err == nil {
		dt.SetText("")
	}
	if tm, err := ui.GetEntry("lg_time"); err == nil {
		tm.SetText("")
	}

	if callsign, err := ui.GetEntry("lg_callsign"); err == nil {
		callsign.SetText("")
	}
	if freq, err := ui.GetEntry("lg_frequency"); err == nil {
		freq.SetText("")
	}
	if loc, err := ui.GetEntry("lg_city"); err == nil {
		loc.SetText("")
	}
	if loc, err := ui.GetEntry("lg_province"); err == nil {
		loc.SetText("")
	}
	if loc, err := ui.GetEntry("lg_country"); err == nil {
		loc.SetText("")
	}

	if sig, err := ui.GetTextBuffer("lg_signal_buffer"); err == nil {
		sig.SetText("")
	}
	if prg, err := ui.GetTextBuffer("lg_programming_buffer"); err == nil {
		prg.SetText("")
	}

	if rcvr, err := ui.GetComboBox("lg_receiver"); err != nil {
		rcvr.SetActive(-1)
	}
	if ant, err := ui.GetComboBox("lg_antenna"); err == nil {
		ant.SetActive(-1)
	}
}

func validateDate() {
	dt, err := ui.GetEntry("lg_date")
	if err != nil {
		log.Println(err)
	}
	d, err := dt.GetText()

	loc, err := time.LoadLocation("Canada/Regina")
	if err != nil {
		log.Println(err)
	}

	_, err = time.ParseInLocation("2006-01-02", d, loc)
}

func validateCall() {
	if c, err := ui.GetEntry("lg_callsign"); err == nil {
		if callsign, err := c.GetText(); err != nil {
			log.Println(err)
		} else {
			if err = loadFCCData(callsign); err != nil {
				log.Println(err)
			}
		}
	}
}

func saveLogEntry(id uint) {
	var logging db.LogEntry
	if dt, err := ui.GetEntry("lg_date"); err == nil {
		logging.Dt, _ = dt.GetText()
	}
	if tm, err := ui.GetEntry("lg_time"); err == nil {
		logging.Tm, _ = tm.GetText()
	}

	if callsign, err := ui.GetEntry("lg_callsign"); err == nil {
		logging.Callsign, _ = callsign.GetText()
	}
	if freq, err := ui.GetEntry("lg_frequency"); err == nil {
		logging.Frequency, _ = freq.GetText()
	}
	if city, err := ui.GetEntry("lg_city"); err == nil {
		logging.City, _ = city.GetText()
	}
	if prov, err := ui.GetEntry("lg_province"); err == nil {
		logging.Prov, _ = prov.GetText()
	}
	if cnty, err := ui.GetEntry("lg_country"); err == nil {
		logging.Cnty, _ = cnty.GetText()
	}

	if sig, err := ui.GetTextBuffer("lg_signal_buffer"); err == nil {
		s, e := sig.GetBounds()
		logging.Signal, _ = sig.GetText(s, e, false)
	}
	if prg, err := ui.GetTextBuffer("lg_programming_buffer"); err == nil {
		s, e := prg.GetBounds()
		logging.Programming, _ = prg.GetText(s, e, false)
	}

	if rcvr, err := ui.GetComboBox("lg_receiver"); err != nil {
		logging.Rcvr = rcvr.GetActive()
	}
	if ant, err := ui.GetComboBox("lg_antenna"); err == nil {
		logging.Ant = ant.GetActive()
	}
	if id != 0 {
		logging.ID = int(id)
		db.UpdateLogging(logging)
	} else {
		db.AddLogging(logging)
	}
	clearLogEntry()
}

func loadForm(id uint) {

	l, err := db.GetLoggingByID(id)
	if err != nil {
		log.Println(err.Error())
	}
	if dt, err := ui.GetEntry("lg_date"); err == nil {
		dt.SetText(l.Dt)
	}
	if tm, err := ui.GetEntry("lg_time"); err == nil {
		tm.SetText(l.Tm)
	}

	if callsign, err := ui.GetEntry("lg_callsign"); err == nil {
		callsign.SetText(l.Callsign)
	}
	if freq, err := ui.GetEntry("lg_frequency"); err == nil {
		freq.SetText(l.Frequency)
	}
	if loc, err := ui.GetEntry("lg_city"); err == nil {
		loc.SetText(l.City)
	}
	if loc, err := ui.GetEntry("lg_province"); err == nil {
		loc.SetText(l.Prov)
	}
	if loc, err := ui.GetEntry("lg_country"); err == nil {
		loc.SetText(l.Cnty)
	}

	if sig, err := ui.GetTextBuffer("lg_signal_buffer"); err == nil {
		sig.SetText(l.Signal)
	}
	if prg, err := ui.GetTextBuffer("lg_programming_buffer"); err == nil {
		prg.SetText(l.Programming)
	}

	if rcvr, err := ui.GetComboBox("lg_receiver"); err == nil {
		rcvr.SetActive(l.Rcvr)
	}
	if ant, err := ui.GetComboBox("lg_antenna"); err == nil {
		ant.SetActive(l.Ant)
	}
}
