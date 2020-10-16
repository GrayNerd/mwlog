package main

import (
	// "fmt"
	"fmt"
	"log"
	"time"
	"strconv"

	"mwlog/db"
	"mwlog/ui"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"

	"github.com/keep94/sunrise"
)

//var oneAndDone bool = false
var loggingWindow *gtk.Window = nil

func notebookSwitcher(pn int) {
	switch pn {
	case 1:
		openLogging(0)
	}
}
func openLogging(id uint) {
	if loggingWindow == nil {
		loggingWindow = ui.GetWindow("logging_window")
		loggingWindow.HideOnDelete()
	} else {
		loggingWindow.ShowAll()
	}
	btn := ui.GetButton("logging_ok_button")

	if id == 0 {
		loggingWindow.SetTitle("Add Logging")
		btn.SetLabel("Add")
		btn.Connect("clicked", func() { saveLogEntry(loggingWindow, 0) })
		clearLogEntry()
		prefillLogEntry()
	} else {
		loggingWindow.SetTitle("Edit Logging")
		btn.SetLabel("Update")
		btn.Connect("clicked", func(b *gtk.Button) { saveLogEntry(loggingWindow, int(id)) })
		loadForm(id)
	}
	loggingWindow.ShowAll()
}

func prefillLogEntry() {
	dt := ui.GetEntry("logging_date")
	tm := ui.GetEntry("logging_time")

	currentTime := time.Now()
	dt.SetText(currentTime.Format("2006-01-02"))
	tm.SetText(currentTime.Format("1504"))

}

func clearLogEntry() {
	ui.GetEntry("logging_date").SetText("")
	ui.GetEntry("logging_time").SetText("")
	ui.GetEntry("logging_station").SetText("")
	ui.GetEntry("logging_frequency").SetText("")
	ui.GetEntry("logging_city").SetText("")
	ui.GetEntry("logging_province").SetText("")
	ui.GetEntry("logging_country").SetText("")
	ui.GetEntry("logging_signal").SetText("")
	ui.GetTextBuffer("logging_remarks_buffer").SetText("")
	ui.GetComboBox("logging_receiver").SetActive(0)
	ui.GetComboBox("logging_antenna").SetActive(0)

	ui.GetLabel("logging_latitude").SetText("")
	ui.GetLabel("logging_longitude").SetText("")

	ui.GetLabel("logging_distance").SetText("")
	ui.GetLabel("logging_bearing").SetText("")
}

func validateDate() {
	ldt := ui.GetEntry("logging_date")
	dt, err := ldt.GetText()
	if err != nil {
		log.Println(err)
	}

	loc, err := time.LoadLocation("Local")
	if err != nil {
		log.Println(err)
	}

	_, err = time.ParseInLocation("2006-01-02", dt, loc)
}

func validateCall(c *gtk.Entry, ev *gdk.Event) {
	log.Println(ev)
	station, _ := ui.GetEntry("logging_station").GetText()
	if len(station) > 0 {
		if err := loadFCCData(station); err != nil {
			log.Println(err)
			c.GrabFocus()
		} else {
			rise, set := calcSunrise()
			ui.GetLabel("logging_sunrise").SetLabel(rise)
			ui.GetLabel("logging_sunset").SetLabel(set)
		}
	}
}

func calcSunrise() (string, string) {
	var s sunrise.Sunrise
	location, err := time.LoadLocation("Local")
	if err != nil {
		log.Println(err)
	}
	dt, _ := ui.GetEntry("logging_date").GetText()

	var day, month, year int
	if _, err := fmt.Sscanf(dt, "%d-%d-%d", &year, &month, &day); err != nil {
		log.Println(err)
	}
	startTime := time.Date(year, time.Month(month), day, 0, 0, 0, 0, location)

	l := ui.GetLabel("logging_latitude").GetLabel()
	lat,_ := strconv.ParseFloat(l, 64)
	l = ui.GetLabel("logging_longitude").GetLabel()
	long,_ := strconv.ParseFloat(l, 64)

	s.Around(lat, long, startTime)
	rise := s.Sunrise().Format("15:04")
	set := s.Sunset().Format("15:04")
	return rise, set

}
func saveLogEntry(win *gtk.Window, id int) {
	var logging db.LogEntry
	logging.Dt, _ = ui.GetEntry("logging_date").GetText()
	logging.Tm, _ = ui.GetEntry("logging_time").GetText()
	logging.Station, _ = ui.GetEntry("logging_station").GetText()
	logging.Frequency, _ = ui.GetEntry("logging_frequency").GetText()
	logging.City, _ = ui.GetEntry("logging_city").GetText()
	logging.Prov, _ = ui.GetEntry("logging_province").GetText()
	logging.Cnty, _ = ui.GetEntry("logging_country").GetText()
	logging.Signal, _ = ui.GetEntry("logging_signal").GetText()

	lrb := ui.GetTextBuffer("logging_remarks_buffer")
	s, e := lrb.GetBounds()
	logging.Remarks, _ = lrb.GetText(s, e, false)

	logging.Rcvr = ui.GetComboBox("logging_receiver").GetActive()
	logging.Ant = ui.GetComboBox("logging_antenna").GetActive()
	logging.Latitude,_ = strconv.ParseFloat(ui.GetLabel("logging_latitude").GetLabel(), 64)
	logging.Longitude,_ = strconv.ParseFloat(ui.GetLabel("logging_longitude").GetLabel(), 64)
	logging.Distance, _ = strconv.ParseFloat(ui.GetLabel("logging_distance").GetLabel(), 64)
	logging.Bearing, _ = strconv.ParseFloat(ui.GetLabel("logging_bearing").GetLabel(), 64)
	logging.Sunrise = ui.GetLabel("logging_sunrise").GetLabel()
	logging.Sunset = ui.GetLabel("logging_sunset").GetLabel()

	if id != 0 {
		logging.ID = int(id)
		db.UpdateLogging(logging)
	} else {
		id = db.AddLogging(logging)
	}

	logbookUpdateRow(id, logging)
	win.Hide()
}

func loadForm(id uint) {

	l, err := db.GetLoggingByID(id)
	if err != nil {
		log.Println(err.Error())
	}
	ui.GetEntry("logging_date").SetText(l.Dt)
	ui.GetEntry("logging_time").SetText(l.Tm)
	ui.GetEntry("logging_station").SetText(l.Station)
	ui.GetEntry("logging_frequency").SetText(l.Frequency)
	ui.GetEntry("logging_city").SetText(l.City)
	ui.GetEntry("logging_province").SetText(l.Prov)
	ui.GetEntry("logging_country").SetText(l.Cnty)
	ui.GetEntry("logging_signal").SetText(l.Signal)
	ui.GetTextBuffer("logging_remarks_buffer").SetText(l.Remarks)
	ui.GetComboBox("logging_receiver").SetActive(l.Rcvr)
	ui.GetComboBox("logging_antenna").SetActive(l.Ant)
	ui.GetLabel("logging_distance").SetLabel(fmt.Sprintf("%.0f", l.Distance))
	ui.GetLabel("logging_bearing").SetLabel(fmt.Sprintf("%.0f", l.Bearing))
	ui.GetLabel("logging_latitude").SetLabel(fmt.Sprintf("%.0f", l.Latitude))
	ui.GetLabel("logging_longitude").SetLabel(fmt.Sprintf("%.0f", l.Longitude))
	ui.GetLabel("logging_sunrise").SetLabel(l.Sunrise)
	ui.GetLabel("logging_sunset").SetLabel(l.Sunset)
}
