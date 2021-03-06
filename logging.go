package main

import (
	// "fmt"
	"fmt"
	"log"

	"math"
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

type logging struct {
	window *gtk.Window
	cancel bool
	rec    db.LogRecord
}

func (l *logging) open(id int) {
	if l.window == nil {
		l.window = ui.GetWindow("logging_window")
		l.window.HideOnDelete()
	}
	l.loadCombos()
	l.window.ShowAll()

	btn := ui.GetButton("logging_save_button")

	var hdl glib.SignalHandle
	switch id {
	case 0:
		fallthrough
	case -1:
		l.window.SetTitle("Add Logging")
		btn.SetLabel("Add")
		hdl = btn.Connect("clicked", func() {
			l.save(l.window, 0)
			btn.HandlerDisconnect(hdl)
		})
		l.clear()
		l.prefill()
		if id == -1 {
			tv := ui.GetTreeView("mwlist_tv")
			s, err := tv.GetSelection()
			if err != nil {
				log.Println(err.Error())
			}
			model, iter, ok := s.GetSelected()
			if !ok {
				log.Println("Unable to GetSelected in onLogbookTreeRowActivated")
			}
			v, _ := model.(*gtk.TreeModel).GetValue(iter, 1)
			call, err := v.GoValue()
			if err != nil {
				log.Println(err.Error())
			}
			e := ui.GetEntry("logging_station")
			e.SetText(call.(string))
			l.validateCall(e)
		}
		break
	default:
		l.window.SetTitle("Edit Logging")
		btn.SetLabel("Update")
		hdl = btn.Connect("clicked", func(btn *gtk.Button) {
			l.save(l.window, id)
			btn.HandlerDisconnect(hdl)
		})
		l.load(id)
		break
	}
	l.window.ShowAll()
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
	i, _ := strconv.Atoi(id.(string))
	l.open(i)
	return false
}

// func (l *logging) formatTabControl (e *gdk.Event) bool{
// 	ev := gdk.EventKeyNewFromEvent(e)
// 	kv := ev.KeyVal()
// 	log.Printf("key val is %d", kv)
// 	log.Printf("tab = %d", kv)
// 	if kv == gdk.KEY_Tab {
// 		// glib.IdleAdd(func() { ui.GetTextView("logging_remarks").GrabFocus() })
// 		// ui.GetTextView("logging_remarks").GrabFocus()
// 	}
// 	return gdk.GDK_EVENT_STOP
// }
// func (l *logging) create() bool {
// 	tv := ui.GetTreeView("mwlist_tv")
// 	s, err := tv.GetSelection()
// 	if err != nil {
// 		log.Println(err.Error())
// 	}
// 	model, iter, ok := s.GetSelected()
// 	if !ok {
// 		log.Println("Unable to GetSelected in onLogbookTreeRowActivated")
// 		return false
// 	}
// 	v, _ := model.(*gtk.TreeModel).GetValue(iter, 1)
// 	id, err := v.GoValue()
// 	if err != nil {
// 		log.Println(err.Error())
// 	}
// 	e := ui.GetEntry("logging_station")
// 	e.SetText(id.(string))
// 	l.validateCall(e)
// 	return true
// }

func (l *logging) prefill() {
	dt := ui.GetEntry("logging_date")
	tm := ui.GetEntry("logging_time")
	ss := ui.GetEntry("logging_sunstatus")
	tm.SetInputPurpose(gtk.INPUT_PURPOSE_DIGITS)

	currentTime := time.Now()
	dt.SetText(currentTime.Format("2006-01-02"))
	tm.SetText(currentTime.Format("1504"))
	ss.SetText(l.calcSunStatus(currentTime))
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
	ui.GetTextBuffer("logging_remarks_buffer").SetText("")
	// ToDo: setup configuration for default format, receiver and antenna
	// ui.GetComboBox("logging_format").SetActive(0)
	ui.GetComboBox("logging_receiver").SetActive(0)
	ui.GetComboBox("logging_antenna").SetActive(0)

	ui.GetEntry("logging_latitude").SetText("")
	ui.GetEntry("logging_longitude").SetText("")

	ui.GetEntry("logging_distance").SetText("")
	ui.GetEntry("logging_bearing").SetText("")

	ui.GetEntry("logging_station").GrabFocus()
	l.rec = db.LogRecord{}
}

func (l *logging) calcSunStatus(tm time.Time) string {
	var s sunrise.Sunrise

	//TODO: use config file instead
	lat := 50.5
	long := -105.5

	s.Around(lat, long, tm)
	rise := s.Sunrise()
	set := s.Sunset()

	dif := rise.Sub(tm).Hours()
	if math.Abs(dif) <= 2.0 {
		return "Sunrise"
	}
	dif = set.Sub(tm).Hours()
	if math.Abs(dif) <= 2.0 {
		return "Sunset"
	}
	time1 := tm.Sub(rise).Hours()
	time2 := tm.Sub(set).Hours()
	if time1 >= 0 && time2 <= 0 {
		return "Daytime"
	}
	return "Nighttime"
}

func (l *logging) save(win *gtk.Window, id int) {

	f := func(w interface{}, msg string) {
		dlg := gtk.MessageDialogNew(l.window, gtk.DIALOG_DESTROY_WITH_PARENT, gtk.MESSAGE_WARNING, gtk.BUTTONS_OK, msg)
		dlg.Run()
		dlg.Destroy()
		var err error

		if e, ok := w.(*gtk.Entry); ok {
			glib.IdleAdd(func() { e.GrabFocus() })
		} else if e, ok := w.(*gtk.TextView); ok {
			glib.IdleAdd(func() { e.GrabFocus() })
		} else if e, ok := w.(*gtk.ComboBox); ok {
			glib.IdleAdd(func() { e.GrabFocus() })
		} else {
			log.Println("Unconfigured widget type in logging.save()")
		}
		if err != nil {
			log.Println("Can't add save IdleAdd")
		}
	}

	l.rec.Dt, _ = ui.GetEntry("logging_date").GetText()
	if len(l.rec.Dt) < 1 {
		f(ui.GetEntry("logging_date"), "Date field cannot be blank")
		return
	}

	l.rec.Tm, _ = ui.GetEntry("logging_time").GetText()
	if len(l.rec.Tm) < 1 {
		f(ui.GetEntry("logging_time"), "Time field cannot be blank")
		return
	}

	l.rec.Station, _ = ui.GetEntry("logging_station").GetText()
	if len(l.rec.Station) < 1 {
		f(ui.GetEntry("logging_station"), "Station field cannot be blank")
		return
	}

	l.rec.Frequency, _ = ui.GetEntry("logging_frequency").GetText()
	l.rec.City, _ = ui.GetEntry("logging_city").GetText()
	l.rec.State, _ = ui.GetEntry("logging_state").GetText()
	l.rec.Country, _ = ui.GetEntry("logging_country").GetText()

	l.rec.Signal, _ = ui.GetEntry("logging_signal").GetText()
	if len(l.rec.Signal) < 1 {
		f(ui.GetEntry("logging_signal"), "Signal field cannot be blank")
		return
	}

	// // l.rec.Format,_ = strconv.Atoi(ui.GetComboBox("logging_format").GetActiveID())
	// fmt, _ := ui.GetEntry("logging_format_entry").GetText()
	// l.rec.Format = db.GetFormatIDByName(fmt)
	// log.Printf("Format saved as %d", l.rec.Format)
	// if l.rec.Format == -1 {
	// 	if len(fmt) < 1 {
	// 		f(ui.GetComboBox("logging_format"), "Format cannot be blank")
	// 	} else {
	// 		f(ui.GetComboBox("logging_format"), "Invalid format")
	// 	}
	// 	return
	// }

	lrb := ui.GetTextBuffer("logging_remarks_buffer")
	s, e := lrb.GetBounds()
	l.rec.Remarks, _ = lrb.GetText(s, e, false)
	if len(l.rec.Remarks) < 1 {
		f(ui.GetTextView("logging_remarks"), "Remarks field cannot be blank")
		return
	}

	rcvr, _ := ui.GetEntry("logging_receiver_entry").GetText()
	l.rec.Receiver = db.GetReceiverIDByName(rcvr)
	log.Printf("receiver saved as %d", l.rec.Receiver)
	if l.rec.Receiver == -1 {
		if len(rcvr) < 1 {
			f(ui.GetComboBox("logging_receiver"), "Receiver cannot be blank")
		} else {
			f(ui.GetComboBox("logging_receiver"), "Invalid receiver")
		}
		return
	}

	ant, _ := ui.GetEntry("logging_antenna_entry").GetText()
	l.rec.Antenna = db.GetAntennaIDByName(ant)
	log.Printf("antenna saved as %d", l.rec.Receiver)
	if l.rec.Antenna == -1 {
		if len(ant) < 1 {
			f(ui.GetComboBox("logging_antenna"), "Antenna cannot be blank")
		} else {
			f(ui.GetComboBox("logging_antenna"), "Invalid antenna")
		}
		return
	}

	t, _ := ui.GetEntry("logging_latitude").GetText()
	l.rec.Latitude, _ = strconv.ParseFloat(t, 64)

	t, _ = ui.GetEntry("logging_longitude").GetText()
	l.rec.Longitude, _ = strconv.ParseFloat(t, 64)

	t, _ = ui.GetEntry("logging_distance").GetText()
	l.rec.Distance, _ = strconv.ParseFloat(t, 64)

	t, _ = ui.GetEntry("logging_bearing").GetText()
	l.rec.Bearing, _ = strconv.ParseFloat(t, 64)

	l.rec.Sunstatus, _ = ui.GetEntry("logging_sunstatus").GetText()

	isNew := true
	if id != 0 {
		l.rec.ID = id
		db.UpdateLogging(&l.rec)
		isNew = false
	} else if id, err := db.AddLogging(l.rec); err != nil {
		log.Println(err.Error())
		win.Hide()
		return
	} else {
		l.rec.ID = id
	}

	logbookUpdateRow(isNew, l.rec)
	win.Hide()
	ui.GetTreeView("logbook_tree").ScrollToPoint(0, 0)
}

func (l *logging) load(id int) {
	rec, err := db.GetLoggingByID(id)
	if err != nil {
		log.Println(err.Error())
	}
	l.rec = rec
	ui.GetEntry("logging_date").SetText(rec.Dt)
	ui.GetEntry("logging_time").SetText(rec.Tm)
	ui.GetEntry("logging_station").SetText(rec.Station)
	ui.GetEntry("logging_frequency").SetText(rec.Frequency)
	ui.GetEntry("logging_city").SetText(rec.City)
	ui.GetEntry("logging_state").SetText(rec.State)
	ui.GetEntry("logging_country").SetText(rec.Country)
	ui.GetEntry("logging_signal").SetText(rec.Signal)
	// ui.GetComboBox("logging_format").SetActive(getComboIndex("format_ls", rec.Format))
	ui.GetTextBuffer("logging_remarks_buffer").SetText(rec.Remarks)
	ui.GetComboBox("logging_receiver").SetActive(getComboIndex("receiver_ls", rec.Receiver))
	// log.Printf("receiver loaded as %d", rec.Receiver)
	ui.GetComboBox("logging_antenna").SetActive(getComboIndex("antenna_ls", rec.Antenna))
	// log.Printf("antenna loaded as %d", rec.Antenna)
	ui.GetEntry("logging_distance").SetText(fmt.Sprintf("%.0f", rec.Distance))
	ui.GetEntry("logging_bearing").SetText(fmt.Sprintf("%.0f", rec.Bearing))
	ui.GetEntry("logging_latitude").SetText(fmt.Sprintf("%.2f", rec.Latitude))
	ui.GetEntry("logging_longitude").SetText(fmt.Sprintf("%.2f", rec.Longitude))
	ui.GetEntry("logging_sunstatus").SetText(rec.Sunstatus)

	ui.GetEntry("logging_date").GrabFocus()
}

// GetComboIndex returns a liststore row for given id
func getComboIndex(ln string, id int) int {
	ls := ui.GetListStore(ln)

	iter, ok := ls.ToTreeModel().GetIterFirst()
	if !ok {
		log.Println("No match found")
	}
	for index := 0; iter != nil; index++ {
		v, err := ls.GetValue(iter, 0)
		if err != nil {
			log.Printf("getComboIndex - unable to GetValue for combo box for %s", ln)
		}
		s, err := v.GoValue()
		if err != nil {
			log.Printf("getComboIndex - unable to GoValue for combo box for %s", ln)
		}
		i, err := strconv.Atoi(s.(string))
		if err != nil {
			log.Printf("getComboIndex - unable to convert string to int for %s", s.(string))
		}
		if i == id {
			return index
		}
		if ls.IterNext(iter) == false {
			break
		}
	}
	log.Printf("getComboIndex combo type not found in index of %s", ln)
	return -1
}

func (l *logging) validateDate(c *gtk.Entry) bool {
	if l.cancel {
		return gdk.GDK_EVENT_PROPAGATE
	}
	dt, err := c.GetText()
	if err != nil {
		log.Println(err.Error())
	}
	if dt == l.rec.Dt {
		return gdk.GDK_EVENT_PROPAGATE
	}

	if len(dt) > 0 {
		d, err := dateparse.ParseLocal(dt)
		if err != nil {
			dlg := gtk.MessageDialogNew(l.window, gtk.DIALOG_DESTROY_WITH_PARENT,
				gtk.MESSAGE_ERROR, gtk.BUTTONS_OK, err.Error())
			dlg.Run()
			dlg.Destroy()
			_ = glib.IdleAdd(func() { c.GrabFocus() })
			return gdk.GDK_EVENT_PROPAGATE
		}
		c.SetText(fmt.Sprintf("%s", d.Format("2006-01-02")))
	}
	return gdk.GDK_EVENT_PROPAGATE
}

func (l *logging) validateTime(c *gtk.Entry) bool {
	if l.cancel {
		return gdk.GDK_EVENT_PROPAGATE
	}

	tm, err := c.GetText()
	if err != nil {
		log.Println(err.Error())
	}
	if tm == l.rec.Tm {
		return gdk.GDK_EVENT_PROPAGATE
	}

	var hours, mins int
	hours, err = strconv.Atoi(tm[:2])
	if err == nil {
		mins, err = strconv.Atoi(tm[2:])
		if err == nil {
			if hours < 24 && mins < 60 {
				de, _ := ui.GetEntry("logging_date").GetText()
				dt, _ := dateparse.ParseLocal(de)
				t := time.Date(dt.Year(), dt.Month(), dt.Day(), hours, mins, 0, 0, dt.Location())

				status := l.calcSunStatus(t)
				ui.GetEntry("logging_sunstatus").SetText(status)

				return gdk.GDK_EVENT_PROPAGATE
			}
		}
	}

	dlg := gtk.MessageDialogNew(l.window, gtk.DIALOG_DESTROY_WITH_PARENT,
		gtk.MESSAGE_ERROR, gtk.BUTTONS_OK, "Invalid time, must be between 0000 and 2359")
	dlg.Run()
	dlg.Destroy()
	_ = glib.IdleAdd(func() { c.GrabFocus() })
	return gdk.GDK_EVENT_PROPAGATE
}

func (l *logging) validateCall(c *gtk.Entry) bool {
	if l.cancel {
		return gdk.GDK_EVENT_PROPAGATE
	}

	station, _ := c.GetText()

	if station == l.rec.Station {
		return gdk.GDK_EVENT_PROPAGATE
	}
	if len(station) > 0 {
		if err := loadMWListData(station); err == nil {
			return gdk.GDK_EVENT_PROPAGATE
		}
		d := gtk.MessageDialogNew(l.window, gtk.DIALOG_DESTROY_WITH_PARENT,
			gtk.MESSAGE_ERROR, gtk.BUTTONS_OK, "Station not found in MWList database")
		d.Run()
		d.Destroy()
	}

	ui.GetEntry("logging_station").SetText("")
	ui.GetEntry("logging_frequency").SetText("")
	ui.GetEntry("logging_city").SetText("")
	ui.GetEntry("logging_state").SetText("")
	ui.GetEntry("logging_country").SetText("")
	ui.GetEntry("logging_latitude").SetText("")
	ui.GetEntry("logging_longitude").SetText("")
	ui.GetEntry("logging_distance").SetText("")
	ui.GetEntry("logging_bearing").SetText("")
	ui.GetEntry("logging_sunstatus").SetText("")

	_ = glib.IdleAdd(func() { c.GrabFocus() })
	return gdk.GDK_EVENT_PROPAGATE
}

func (l *logging) loadCombos() {
	l.loadReceivers()
	l.loadAntennas()
	// l.loadFormats()
}

func (l *logging) loadReceivers() {
	ls := ui.GetListStore("receiver_ls")
	ls.Clear()

	rows := db.GetAllReceivers()
	defer rows.Close()

	var id int
	var name string
	var iter *gtk.TreeIter
	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Println(err.Error())
		}
		col := []int{0, 1}
		var val []interface{}
		val = append(val, strconv.Itoa(id), name)

		if err = ls.InsertWithValues(iter, 0, col, val); err != nil {
			log.Println(err.Error())
		}
	}
}

func (l *logging) loadAntennas() {
	ls := ui.GetListStore("antenna_ls")
	ls.Clear()

	rows := db.GetAllAntennas()
	defer rows.Close()

	var id int
	var name string
	var iter *gtk.TreeIter
	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Println(err.Error())
		}
		col := []int{0, 1}
		var val []interface{}
		val = append(val, id, name)
		
		if err = ls.InsertWithValues(iter, 0, col, val); err != nil {
			log.Println(err.Error())
		}
	}
}

// func (l *logging) loadFormats() {
// 	ls := ui.GetListStore("format_ls")
// 	ls.Clear()
// 	rows := db.GetAllFormats()
// 	defer rows.Close()

// 	var id int
// 	var name string
// 	var iter *gtk.TreeIter
// 	for rows.Next() {
// 		err := rows.Scan(&id, &name)
// 		if err != nil {
// 			log.Println(err.Error())
// 		}
// 		// iter = ls.Append()
// 		col := []int{0, 1}
// 		var val []interface{}
// 		val = append(val, strconv.Itoa(id), name)
// 		// log.Println(id, name)
// 		if err = ls.InsertWithValues(iter, 0, col, val); err != nil {
// 			log.Println(err.Error())
// 		}
// 	}
// }
