package main

import (
	"fmt"
	"log"
	"mwlog/db"
	"mwlog/ui"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

func logbookLoad() {
	ls:= ui.GetListStore("logbook_store")
	ls.Clear()

	rows, err := db.GetLogBookStore()
	if err != nil {
		log.Println(err.Error())
	}
	defer rows.Close()
	var id uint
	var dt, tm, station, frequency, city, province, country, signal, remarks string
	for rows.Next() {
		rows.Scan(&id, &dt, &tm, &station, &frequency, &city, &province, &country, &signal, &remarks)
		iter := ls.Append()
		col := []int{0,1,2,3,4,5,6,7,8}
		var val []interface{}
		val = append(val, id, dt, tm, station, frequency, fmt.Sprintf("%s, %s %s",city, province, country), signal, remarks)
		if err = ls.InsertWithValues(iter, 0, col, val); err != nil {
			log.Println(err.Error())
		}
	}
}

func onLogbookTreeKeyPressEvent(tv *gtk.TreeView, e *gdk.Event) bool {
	ek := gdk.EventKey{e}
	if ek.KeyVal() == gdk.KEY_Delete {
		s, err := tv.GetSelection()
		if err != nil {
			log.Println(err.Error())
		}
		model, iter, ok := s.GetSelected()
		if !ok {
			log.Println("Unable to GetSelected in logbookDeleteSelected")
			return false
		}

		v, _ := model.(*gtk.TreeModel).GetValue(iter, 0)
		id, err := v.GoValue()
		if err != nil {
			log.Println(err.Error())
		}
		db.DeleteLogging(id.(uint))

		ui.GetListStore("logbook_store").Remove(iter)
		return false
	}
	return true
}

func onLogbookTreeRowActivated(tv *gtk.TreeView) bool {
	s, err := tv.GetSelection()
	if err != nil {
		log.Println(err.Error())
	}
	model, iter, ok := s.GetSelected()
	if !ok {
		log.Println("Unable to GetSelected in logbookEditSelected")
		return false
	}

	v, _ := model.(*gtk.TreeModel).GetValue(iter, 0)
	id, err := v.GoValue()
	if err != nil {
		log.Println(err.Error())
	}

	openLogging(id.(uint))
	return false
}

func logbookUpdateRow(id int, logging db.LogEntry) {
	tv := ui.GetTreeView("logbook_tree")
	ls := ui.GetListStore("logbook_store")

	s, err := tv.GetSelection()
	if err != nil {
		log.Println(err.Error())
	}
	_, iter, ok := s.GetSelected()
	if !ok {
		// nothing selected...must be new entry
		iter = ls.Append()
	}

	if err = ls.SetValue(iter, 1, logging.Dt); err != nil {
		log.Println(err.Error())
	}
	if err = ls.SetValue(iter, 2, logging.Tm); err != nil {
		log.Println(err.Error())
	}
	if err = ls.SetValue(iter, 3, logging.Station); err != nil {
		log.Println(err.Error())
	}
	if err = ls.SetValue(iter, 4, logging.Frequency); err != nil {
		log.Println(err.Error())
	}
	if err = ls.SetValue(iter, 5, logging.City+", "+logging.Prov+" "+logging.Cnty); err != nil {
		log.Println(err.Error())
	}
	if err = ls.SetValue(iter, 6, logging.Signal); err != nil {
		log.Println(err.Error())
	}
	if err = ls.SetValue(iter, 7, logging.Remarks); err != nil {
		log.Println(err.Error())
	}
}
