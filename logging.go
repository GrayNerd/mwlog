package main

import (
	// "fmt"
	"fmt"
	"log"

	"strconv"
	"time"

	"mwlog/db"
	"mwlog/ui"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	"github.com/araddon/dateparse"
	"github.com/keep94/sunrise"
)

//var oneAndDone bool = false
type logging struct {
	loggingWindow *gtk.Window
	cancel        bool
}

func (l *logging) open(id uint) {
	if l.loggingWindow == nil {
		l.loggingWindow = ui.GetWindow("logging_window")
		l.loggingWindow.HideOnDelete()
	}
	l.loggingWindow.ShowAll()

	btn := ui.GetButton("logging_save_button")

	var hdl glib.SignalHandle
	if id == 0 {
		l.loggingWindow.SetTitle("Add Logging")
		btn.SetLabel("Add")
		hdl, _ = btn.Connect("clicked", func() {
			l.save(l.loggingWindow, 0)
			btn.HandlerDisconnect(hdl)
		})
		l.clear()
		l.prefill()
	} else {
		l.loggingWindow.SetTitle("Edit Logging")
		btn.SetLabel("Update")
		i := id
		hdl, _ = btn.Connect("clicked", func(btn *gtk.Button) {
			l.save(l.loggingWindow, int(i))
			btn.HandlerDisconnect(hdl)
		})
		l.load(id)
	}
	l.loggingWindow.ShowAll()
}
func (l *logging) edit() bool {
	tv := ui.GetTreeView("logbook_tree")
	s, err := tv.GetSelection()
	if err != nil {
		log.Println(err.Error())
	}
	model, iter, ok := s.GetSelected()
	if !ok {
		log.Println("Unable to GetSelected in onLogbookTreeRowActivated")
		return false
	}

	v, _ := model.(*gtk.TreeModel).GetValue(iter, 0)
	id, err := v.GoValue()
	if err != nil {
		log.Println(err.Error())
	}

	l.open(id.(uint))
	return false
}

func (l *logging) prefill() {
	dt := ui.GetEntry("logging_date")
	tm := ui.GetEntry("logging_time")
	tm.SetInputPurpose(gtk.INPUT_PURPOSE_DIGITS)

	currentTime := time.Now()
	dt.SetText(currentTime.Format("2006-01-02"))
	tm.SetText(currentTime.Format("1504"))

}

func (l *logging) clear() {
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

func (l *logging) validateDate(c *gtk.Entry) bool {
	dt, err := c.GetText()
	if err != nil {
		log.Println(err.Error())
	}

	if len(dt) > 0 {
		d, err := dateparse.ParseLocal(dt)
		if err != nil {
			dlg := gtk.MessageDialogNew(l.loggingWindow, gtk.DIALOG_DESTROY_WITH_PARENT,
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

func (l *logging) validateTime(c *gtk.Entry) bool {
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

	dlg := gtk.MessageDialogNew(l.loggingWindow, gtk.DIALOG_DESTROY_WITH_PARENT,
		gtk.MESSAGE_ERROR, gtk.BUTTONS_OK, "Invalid time, must be between 0000 and 2359")
	dlg.Run()
	dlg.Destroy()
	if _, err := glib.IdleAdd(func() { c.GrabFocus() }); err != nil {
		log.Println("Can't add validateDate IdleAdd")
	}
	return gdk.GDK_EVENT_PROPAGATE
}

func (l *logging) validateCall(c *gtk.Entry) bool {

	// log.Printf("%v", ev)
	station, _ := c.GetText()
	if len(station) == 0 {
		return gdk.GDK_EVENT_PROPAGATE
	}

	if err := loadMWListData(station); err == nil {
		rise, set := l.calcSunTimes()
		ui.GetLabel("logging_sunrise").SetLabel(rise)
		ui.GetLabel("logging_sunset").SetLabel(set)
		return gdk.GDK_EVENT_PROPAGATE
	}

	d := gtk.MessageDialogNew(l.loggingWindow, gtk.DIALOG_DESTROY_WITH_PARENT,
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

func (l *logging) calcSunTimes() (string, string) {
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

	x := ui.GetLabel("logging_latitude").GetLabel()
	lat, _ := strconv.ParseFloat(x, 64)
	x = ui.GetLabel("logging_longitude").GetLabel()
	long, _ := strconv.ParseFloat(x, 64)

	s.Around(lat, long, startTime)
	rise := s.Sunrise().Format("15:04")
	set := s.Sunset().Format("15:04")
	return rise, set

}

func (l *logging) save(win *gtk.Window, id int) {

	var rec db.LogRecord
	f := func(w interface{}, msg string) {
		_=w
		dlg := gtk.MessageDialogNew(l.loggingWindow, gtk.DIALOG_DESTROY_WITH_PARENT, gtk.MESSAGE_WARNING, gtk.BUTTONS_OK, msg)
		dlg.Run()
		dlg.Destroy()
		// if _, err := glib.IdleAdd(func() { w.(*gtk.Widget).GrabFocus() }); err != nil {
		// 	log.Println("Can't add save IdleAdd")
		// }
	}

	rec.Dt, _ = ui.GetEntry("logging_date").GetText()
	if len(rec.Dt) < 1 {
		f(ui.GetEntry("logging_date"), "Date field cannot be blank")
		return
	}

	rec.Tm, _ = ui.GetEntry("logging_time").GetText()
	if len(rec.Tm) < 1 {
		f(ui.GetEntry("logging_time"), "Time field cannot be blank")
		return
	}

	rec.Station, _ = ui.GetEntry("logging_station").GetText()
	if len(rec.Station) < 1 {
		f(ui.GetEntry("logging_station"), "Station field cannot be blank")
		return
	}

	rec.Frequency, _ = ui.GetEntry("logging_frequency").GetText()
	rec.City, _ = ui.GetEntry("logging_city").GetText()
	rec.State, _ = ui.GetEntry("logging_state").GetText()
	rec.Cnty, _ = ui.GetEntry("logging_country").GetText()

	rec.Signal, _ = ui.GetEntry("logging_signal").GetText()
	if len(rec.Signal) < 1 {
		f(ui.GetEntry("logging_signal"), "Signal field cannot be blank")
		return
	}

	rec.Format, _ = ui.GetEntry("logging_format").GetText()

	lrb := ui.GetTextBuffer("logging_remarks_buffer")
	s, e := lrb.GetBounds()
	rec.Remarks, _ = lrb.GetText(s, e, false)
	if len(rec.Remarks) < 1 {
		f(ui.GetTextView("logging_remarks"), "Remarks field cannot be blank")
		return
	}

	rec.Rcvr = ui.GetComboBox("logging_receiver").GetActive()
	if rec.Rcvr == -1 {
		f(ui.GetComboBox("logging_receiver"), "Receiver field cannot be blank")
		return
	}

	rec.Ant = ui.GetComboBox("logging_antenna").GetActive()
	if rec.Ant == -1 {
		f(ui.GetComboBox("logging_antenna"), "Antenna field cannot be blank")
		return
	}

	rec.Latitude, _ = strconv.ParseFloat(ui.GetLabel("logging_latitude").GetLabel(), 64)
	rec.Longitude, _ = strconv.ParseFloat(ui.GetLabel("logging_longitude").GetLabel(), 64)
	rec.Distance, _ = strconv.ParseFloat(ui.GetLabel("logging_distance").GetLabel(), 64)
	rec.Bearing, _ = strconv.ParseFloat(ui.GetLabel("logging_bearing").GetLabel(), 64)
	rec.Sunrise = ui.GetLabel("logging_sunrise").GetLabel()
	rec.Sunset = ui.GetLabel("logging_sunset").GetLabel()

	var isNew bool
	if id != 0 {
		rec.ID = int(id)
		db.UpdateLogging(&rec)
		isNew = false
	} else {
		db.AddLogging(&rec)
		isNew = true
	}

	logbookUpdateRow(isNew, &rec)
	win.Hide()
	ui.GetTreeView("logbook_tree").ScrollToPoint(0, 0)
}

func (l *logging) load(id uint) {

	rec, err := db.GetLoggingByID(id)
	if err != nil {
		log.Println(err.Error())
	}
	ui.GetEntry("logging_date").SetText(rec.Dt)
	ui.GetEntry("logging_time").SetText(rec.Tm)
	ui.GetEntry("logging_station").SetText(rec.Station)
	ui.GetEntry("logging_frequency").SetText(rec.Frequency)
	ui.GetEntry("logging_city").SetText(rec.City)
	ui.GetEntry("logging_state").SetText(rec.State)
	ui.GetEntry("logging_country").SetText(rec.Cnty)
	ui.GetEntry("logging_signal").SetText(rec.Signal)
	ui.GetEntry("logging_format").SetText(rec.Format)
	ui.GetTextBuffer("logging_remarks_buffer").SetText(rec.Remarks)
	ui.GetComboBox("logging_receiver").SetActive(rec.Rcvr)
	ui.GetComboBox("logging_antenna").SetActive(rec.Ant)
	ui.GetLabel("logging_distance").SetLabel(fmt.Sprintf("%.0f", rec.Distance))
	ui.GetLabel("logging_bearing").SetLabel(fmt.Sprintf("%.0f", rec.Bearing))
	ui.GetLabel("logging_latitude").SetLabel(fmt.Sprintf("%.0f", rec.Latitude))
	ui.GetLabel("logging_longitude").SetLabel(fmt.Sprintf("%.0f", rec.Longitude))
	ui.GetLabel("logging_sunrise").SetLabel(rec.Sunrise)
	ui.GetLabel("logging_sunset").SetLabel(rec.Sunset)

	ui.GetEntry("logging_date").GrabFocus()
}
