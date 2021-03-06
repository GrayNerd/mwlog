package main

import (
	// "fmt"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"
	"unsafe"

	"mwlog/db"
	"mwlog/ui"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"

	"github.com/araddon/dateparse"
	"github.com/keep94/sunrise"
)

//var oneAndDone bool = false
var loggingWindow *gtk.Window = nil

func openLogging(id uint) {
	if loggingWindow == nil {
		loggingWindow = ui.GetWindow("logging_window")
		loggingWindow.HideOnDelete()
	}
	loggingWindow.ShowAll()

	btn := ui.GetButton("logging_save_button")

	var hdl glib.SignalHandle
	if id == 0 {
		loggingWindow.SetTitle("Add Logging")
		btn.SetLabel("Add")
		hdl, _ = btn.Connect("clicked", func() {
			saveLogEntry(loggingWindow, 0)
			btn.HandlerDisconnect(hdl)
		})
		clearLogEntry()
		prefillLogEntry()
	} else {
		loggingWindow.SetTitle("Edit Logging")
		btn.SetLabel("Update")
		i := id
		hdl, _ = btn.Connect("clicked", func(btn *gtk.Button) {
			saveLogEntry(loggingWindow, int(i))
			btn.HandlerDisconnect(hdl)
		})
		loadForm(id)
	}
	loggingWindow.ShowAll()
}

func prefillLogEntry() {
	dt := ui.GetEntry("logging_date")
	tm := ui.GetEntry("logging_time")
	tm.SetInputPurpose(gtk.INPUT_PURPOSE_DIGITS)

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
	ui.GetEntry("logging_state").SetText("")
	ui.GetEntry("logging_country").SetText("")
	ui.GetEntry("logging_signal").SetText("")
	ui.GetEntry("logging_format").SetText("")
	ui.GetTextBuffer("logging_remarks_buffer").SetText("")
	ui.GetComboBox("logging_receiver").SetActive(0)
	ui.GetComboBox("logging_antenna").SetActive(0)

	ui.GetLabel("logging_latitude").SetText("")
	ui.GetLabel("logging_longitude").SetText("")

	ui.GetLabel("logging_distance").SetText("")
	ui.GetLabel("logging_bearing").SetText("")

	ui.GetEntry("logging_station").GrabFocus()
}

func validateDate(c *gtk.Entry) bool {
	dt, err := c.GetText()
	if err != nil {
		log.Println(err.Error())
	}

	if len(dt) > 0 {
		d, err := dateparse.ParseLocal(dt)
		if err != nil {
			dlg := gtk.MessageDialogNew(loggingWindow, gtk.DIALOG_DESTROY_WITH_PARENT,
				gtk.MESSAGE_ERROR, gtk.BUTTONS_OK, err.Error())
			dlg.Run()
			dlg.Destroy()
			if _, err := glib.IdleAdd(func() { c.GrabFocus() }); err != nil {
				log.Println("Can't add validateDate IdleAdd")
			}
			return gdk.GDK_EVENT_PROPAGATE
		}
		c.SetText(fmt.Sprintf("%s", d.Format("2006-01-02")))
	}
	return gdk.GDK_EVENT_PROPAGATE
}

func validateTime(c *gtk.Entry) bool {
	tm, err := c.GetText()
	if err != nil {
		log.Println(err.Error())
	}

	hours, err := strconv.Atoi(tm[:2])
	if err == nil {
		mins, err := strconv.Atoi(tm[2:])
		if err == nil {
			if hours < 24 && mins < 60 {
				return gdk.GDK_EVENT_PROPAGATE
			}
		}
	}
	
	dlg := gtk.MessageDialogNew(loggingWindow, gtk.DIALOG_DESTROY_WITH_PARENT,
		gtk.MESSAGE_ERROR, gtk.BUTTONS_OK, "Invalid time, must be between 0000 and 2359")
	dlg.Run()
	dlg.Destroy()
	if _, err := glib.IdleAdd(func() { c.GrabFocus() }); err != nil {
		log.Println("Can't add validateDate IdleAdd")
	}
	return gdk.GDK_EVENT_PROPAGATE
}

func validateCall(c *gtk.Entry) bool {

	// log.Printf("%v", ev)
	station, _ := c.GetText()
	if len(station) == 0 {
		return gdk.GDK_EVENT_PROPAGATE
	}

	if err := loadMWListData(station); err == nil {
		rise, set := calcSunTimes()
		ui.GetLabel("logging_sunrise").SetLabel(rise)
		ui.GetLabel("logging_sunset").SetLabel(set)
		return gdk.GDK_EVENT_PROPAGATE
	}

	d := gtk.MessageDialogNew(loggingWindow, gtk.DIALOG_DESTROY_WITH_PARENT,
		gtk.MESSAGE_ERROR, gtk.BUTTONS_OK, "Station not found in MWList database")
	d.Run()
	d.Destroy()

	ui.GetEntry("logging_station").SetText("")
	ui.GetEntry("logging_frequency").SetText("")
	ui.GetEntry("logging_city").SetText("")
	ui.GetEntry("logging_state").SetText("")
	ui.GetEntry("logging_country").SetText("")
	ui.GetLabel("logging_latitude").SetLabel("")
	ui.GetLabel("logging_longitude").SetLabel("")
	ui.GetLabel("logging_distance").SetLabel("")
	ui.GetLabel("logging_bearing").SetLabel("")
	ui.GetLabel("logging_sunrise").SetLabel("")
	ui.GetLabel("logging_sunset").SetLabel("")

	if _, err := glib.IdleAdd(func() { c.GrabFocus() }); err != nil {
		println("Can't add idleadd")
	}

	return gdk.GDK_EVENT_PROPAGATE
}

func calcSunTimes() (string, string) {
	var s sunrise.Sunrise
	location, err := time.LoadLocation("Local")
	if err != nil {
		log.Println(err.Error())
	}
	dt, _ := ui.GetEntry("logging_date").GetText()

	var day, month, year int
	if _, err := fmt.Sscanf(dt, "%d-%d-%d", &year, &month, &day); err != nil {
		log.Println(err.Error())
	}
	startTime := time.Date(year, time.Month(month), day, 0, 0, 0, 0, location)

	l := ui.GetLabel("logging_latitude").GetLabel()
	lat, _ := strconv.ParseFloat(l, 64)
	l = ui.GetLabel("logging_longitude").GetLabel()
	long, _ := strconv.ParseFloat(l, 64)

	s.Around(lat, long, startTime)
	rise := s.Sunrise().Format("15:04")
	set := s.Sunset().Format("15:04")
	return rise, set

}

func saveLogEntry(win *gtk.Window, id int) {

	var logging db.LogEntry

	f := func(w interface{}, msg string) {
		dlg := gtk.MessageDialogNew(loggingWindow, gtk.DIALOG_DESTROY_WITH_PARENT, gtk.MESSAGE_WARNING, gtk.BUTTONS_OK, msg)
		dlg.Run()
		dlg.Destroy()

		x:= (*reflect.Type)(unsafe.Pointer(&w))
		if _, err := glib.IdleAdd(func() { w.(*gtk.Widget).GrabFocus() }); err != nil {
			log.Println("Can't add saveLogEntry IdleAdd")
		}
	}

	logging.Dt, _ = ui.GetEntry("logging_date").GetText()
	if len(logging.Dt) < 1 {
		f(ui.GetEntry("logging_date"), "Date field cannot be blank")

	}

	logging.Tm, _ = ui.GetEntry("logging_time").GetText()
	if len(logging.Tm) < 1 {
		f(ui.GetEntry("logging_time"), "Time field cannot be blank")
	}

	logging.Station, _ = ui.GetEntry("logging_station").GetText()
	if len(logging.Station) < 1 {
		f(ui.GetEntry("logging_station"), "Station field cannot be blank")
	}

	logging.Frequency, _ = ui.GetEntry("logging_frequency").GetText()
	logging.City, _ = ui.GetEntry("logging_city").GetText()
	logging.State, _ = ui.GetEntry("logging_state").GetText()
	logging.Cnty, _ = ui.GetEntry("logging_country").GetText()

	logging.Signal, _ = ui.GetEntry("logging_signal").GetText()
	if len(logging.Signal) < 1 {
		f(ui.GetEntry("logging_signal"), "Signal field cannot be blank")
	}

	logging.Format, _ = ui.GetEntry("logging_format").GetText()

	lrb := ui.GetTextBuffer("logging_remarks_buffer")
	s, e := lrb.GetBounds()
	logging.Remarks, _ = lrb.GetText(s, e, false)
	if len(logging.Remarks) < 1 {
		f(ui.GetTextView("logging_remarks"), "Remarks field cannot be blank")
	}

	logging.Rcvr = ui.GetComboBox("logging_receiver").GetActive()
	if logging.Rcvr == -1 {
		f(ui.GetComboBox("logging_receiver"), "Receiver field cannot be blank")
	}

	logging.Ant = ui.GetComboBox("logging_antenna").GetActive()
	if logging.Ant == -1 {
		f(ui.GetComboBox("logging_antenna"), "Antenna field cannot be blank")
	}

	logging.Latitude, _ = strconv.ParseFloat(ui.GetLabel("logging_latitude").GetLabel(), 64)
	logging.Longitude, _ = strconv.ParseFloat(ui.GetLabel("logging_longitude").GetLabel(), 64)
	logging.Distance, _ = strconv.ParseFloat(ui.GetLabel("logging_distance").GetLabel(), 64)
	logging.Bearing, _ = strconv.ParseFloat(ui.GetLabel("logging_bearing").GetLabel(), 64)
	logging.Sunrise = ui.GetLabel("logging_sunrise").GetLabel()
	logging.Sunset = ui.GetLabel("logging_sunset").GetLabel()

	var isNew bool
	if id != 0 {
		logging.ID = int(id)
		db.UpdateLogging(&logging)
		isNew = false
	} else {
		db.AddLogging(&logging)
		isNew = true
	}

	logbookUpdateRow(isNew, &logging)
	win.Hide()
	ui.GetTreeView("logbook_tree").ScrollToPoint(0, 0)
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
	ui.GetEntry("logging_state").SetText(l.State)
	ui.GetEntry("logging_country").SetText(l.Cnty)
	ui.GetEntry("logging_signal").SetText(l.Signal)
	ui.GetEntry("logging_format").SetText(l.Format)
	ui.GetTextBuffer("logging_remarks_buffer").SetText(l.Remarks)
	ui.GetComboBox("logging_receiver").SetActive(l.Rcvr)
	ui.GetComboBox("logging_antenna").SetActive(l.Ant)
	ui.GetLabel("logging_distance").SetLabel(fmt.Sprintf("%.0f", l.Distance))
	ui.GetLabel("logging_bearing").SetLabel(fmt.Sprintf("%.0f", l.Bearing))
	ui.GetLabel("logging_latitude").SetLabel(fmt.Sprintf("%.0f", l.Latitude))
	ui.GetLabel("logging_longitude").SetLabel(fmt.Sprintf("%.0f", l.Longitude))
	ui.GetLabel("logging_sunrise").SetLabel(l.Sunrise)
	ui.GetLabel("logging_sunset").SetLabel(l.Sunset)

	ui.GetEntry("logging_date").GrabFocus()
}
