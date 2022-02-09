package main

import (
	"fmt"
	"log"
	"strconv"

	"mwlog/db"
	"mwlog/ui"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

func showLogbook() {
	nb := ui.GetNotebook("notebook")
	nb.SetCurrentPage(0) // logbook
	loadLogbook()
	nb.ShowAll()
}

func loadLogbook() {
	ls := ui.GetListStore("logbook_store")
	ls.Clear()

	rows, err := db.GetLogBookStore()
	if err != nil {
		log.Println(err.Error())
	}
	defer rows.Close()

	var id, dt, tm, station, frequency, city, state, country, signal, remarks string
	var iter *gtk.TreeIter
	for rows.Next() {
		err := rows.Scan(&id, &dt, &tm, &station, &frequency, &city, &state, &country, &signal, &remarks)
		if err != nil {
			log.Println(err.Error())
		}
		// iter = ls.Append()
		col := []int{0, 1, 2, 3, 4, 5, 6, 7}
		var val []interface{}
		val = append(val, id, dt, tm, frequency, station, fmt.Sprintf("%s, %s %s", city, state, country), signal, remarks)
		if err = ls.InsertWithValues(iter, 0, col, val); err != nil {
			log.Println(err.Error())
		}
	}
}

func logbookUpdateRow(new bool, l db.LogRecord) {
	tv := ui.GetTreeView("logbook_tree")
	ls := ui.GetListStore("logbook_store")

	var iter *gtk.TreeIter
	var ok bool
	if new {
		iter = ls.Prepend()
	} else {
		s, err := tv.GetSelection()
		if err != nil {
			log.Println(err.Error())
		}
		_, iter, ok = s.GetSelected()
		if !ok {
			log.Println("logbookUpdateRow: unable to get selected")
		}
	}
	if err := ls.SetValue(iter, 0, l.ID); err != nil {
		log.Println(err.Error())
	}
	if err := ls.SetValue(iter, 1, l.Dt); err != nil {
		log.Println(err.Error())
	}
	if err := ls.SetValue(iter, 2, l.Tm); err != nil {
		log.Println(err.Error())
	}
	if err := ls.SetValue(iter, 4, l.Station); err != nil {
		log.Println(err.Error())
	}
	if err := ls.SetValue(iter, 3, l.Frequency); err != nil {
		log.Println(err.Error())
	}
	if err := ls.SetValue(iter, 5, l.City+", "+l.State+" "+l.Country); err != nil {
		log.Println(err.Error())
	}
	if err := ls.SetValue(iter, 6, l.Signal); err != nil {
		log.Println(err.Error())
	}
	if err := ls.SetValue(iter, 7, l.Remarks); err != nil {
		log.Println(err.Error())
	}
	ui.GetTreeSelection("lb_tree_selection").SelectIter(iter)
}

func onLogbookTreeKeyPressEvent(tv *gtk.TreeView, e *gdk.Event) bool {
	ek := gdk.EventKeyNewFromEvent(e)
	if ek.KeyVal() == gdk.KEY_Delete {
		// TODO: add confirmation option
		doDelete(tv)
	}
	return true
}

func onLogbookTreeButtonPressEvent(_ *gtk.TreeView, e *gdk.Event, l logging) {
	// Both of these work for a double click...which one's better?
	// gdk.EVENT_2BUTTON_PRESS
	// gdk.EVENT_DOUBLE_BUTTON_PRESS

	eb := gdk.EventButtonNewFromEvent(e)
	switch eb.Type() {
	case gdk.EVENT_DOUBLE_BUTTON_PRESS: // Double click
		l.edit()
	}

	if eb.Button() == gdk.BUTTON_SECONDARY { // Right click
		//x := eb.X()
		//y := eb.Y()
		//log.Printf("%f -- %f", x, y)
		ui.GetMenu("logbook_popup").PopupAtPointer(e)
	}
}

func onLogbookDelete() {
	tv := ui.GetTreeView("logbook_tree")
	doDelete(tv)
}

func doDelete(tv *gtk.TreeView) {
	s, err := tv.GetSelection()
	if err != nil {
		log.Println(err.Error())
		return
	}
	model, iter, ok := s.GetSelected()
	if !ok {
		log.Println("Unable to GetSelected in logbookDeleteSelected")
		return
	}

	v, _ := model.(*gtk.TreeModel).GetValue(iter, 0)
	id, err := v.GoValue()
	if err != nil {
		log.Println(err.Error())
		return
	}

	dialog := gtk.MessageDialogNew(nil, gtk.DIALOG_DESTROY_WITH_PARENT, gtk.MESSAGE_QUESTION, gtk.BUTTONS_OK_CANCEL, "%s", "\nDelete the logging?")
	res := dialog.Run()
	if res == gtk.RESPONSE_OK {
		i, _ := strconv.Atoi(id.(string))
		db.DeleteLogging(i)
		ui.GetListStore("logbook_store").Remove(iter)
	}
	dialog.Close()
}
