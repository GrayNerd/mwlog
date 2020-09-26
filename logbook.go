package main

import (
	"log"
	"mwlog/db"
	"mwlog/ui"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

func initLogBook() {
	db.FillLogBookStore()

}

func logbookDeleteSelected(tv *gtk.TreeView, e *gdk.Event) bool {
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

		ls, err := ui.GetListStore("logbook_store")
		if err != nil {
			log.Println(err.Error())
		}
		ls.Remove(iter)
		return false
	}
	return true
}

func logbookEditSelected(tv *gtk.TreeView) bool {
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
	nb, _ := ui.GetNotebook("notebook")
	nb.SetCurrentPage(1)

	b, _ := ui.GetButton("lg_ok_button")
	b.SetLabel("Update")
	b.Connect("clicked", func() { saveLogEntry(id.(uint)) })

	loadForm(id.(uint))
	return false
}
