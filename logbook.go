package main

import (
	"fmt"
	"log"
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

	// var id int
	var id, dt, tm, station, frequency, city, state, country, signal, remarks string
	var iter *gtk.TreeIter
	for rows.Next() {
		rows.Scan(&id, &dt, &tm, &station, &frequency, &city, &state, &country, &signal, &remarks)
		// iter = ls.Append()
		col := []int{0, 1, 2, 3, 4, 5, 6, 7}
		var val []interface{}
		val = append(val, id, dt, tm, frequency, station, fmt.Sprintf("%s, %s %s", city, state, country), signal, remarks)
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
			return false
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
			return false
		}

		dialog := gtk.MessageDialogNew(nil, gtk.DIALOG_DESTROY_WITH_PARENT, gtk.MESSAGE_QUESTION, gtk.BUTTONS_OK_CANCEL, "%s", "\nDelete the logging?")
		// dialog.SetTitle("Delete Logging?")
		// dialog.SetSizeRequest(300, 200)
		res := dialog.Run()
		if res == gtk.RESPONSE_OK {
			db.DeleteLogging(id.(int))
			ui.GetListStore("logbook_store").Remove(iter)
		}
		dialog.Close()
		return false
	}
	return true
}

func logbookUpdateRow(new bool, logging *db.LogRecord) {
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
	if err := ls.SetValue(iter, 0, logging.ID); err != nil {
		log.Println(err.Error())
	}
	if err := ls.SetValue(iter, 1, logging.Dt); err != nil {
		log.Println(err.Error())
	}
	if err := ls.SetValue(iter, 2, logging.Tm); err != nil {
		log.Println(err.Error())
	}
	if err := ls.SetValue(iter, 3, logging.Station); err != nil {
		log.Println(err.Error())
	}
	if err := ls.SetValue(iter, 4, logging.Frequency); err != nil {
		log.Println(err.Error())
	}
	if err := ls.SetValue(iter, 5, logging.City+", "+logging.State+" "+logging.Cnty); err != nil {
		log.Println(err.Error())
	}
	if err := ls.SetValue(iter, 6, logging.Signal); err != nil {
		log.Println(err.Error())
	}
	if err := ls.SetValue(iter, 7, logging.Format); err != nil {
		log.Println(err.Error())
	}
	if err := ls.SetValue(iter, 8, logging.Remarks); err != nil {
		log.Println(err.Error())
	}
	ui.GetTreeSelection("lb_tree_selection").SelectIter(iter)
}
