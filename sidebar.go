package main

import (
	"fmt"
	"log"

	"mwlog/ui"

	"github.com/gotk3/gotk3/gtk"
)

func buildSidebar() {
	treeStore, err := ui.GetTreeStore("sidebar_ls")
	if err != nil {
		log.Println(err.Error())
	}

	_, err = sidebarAppend(treeStore, nil, "Log Book", "0")
	if err != nil {
		log.Println(err.Error())
	}

	_, err = sidebarAppend(treeStore, nil, "Add Logging", "1")
	if err != nil {
		log.Println(err.Error())
	}

	_, err = sidebarAppend(treeStore, nil, "Statistics", "2")
	if err != nil {
		log.Println(err.Error())
	}

	_, err = sidebarAppend(treeStore, nil, "FCC Lookup", "3")
	if err != nil {
		log.Println(err.Error())
	}

	iter, err := sidebarAppend(treeStore, nil, "Settings", "4")
	if err != nil {
		log.Println(err.Error())
	}
	_, err = sidebarAppend(treeStore, iter, "Languages", "4:0")
	if err != nil {
		log.Println(err.Error())
	}
	_, err = sidebarAppend(treeStore, iter, "Equipment", "4:1")
	if err != nil {
		log.Println(err.Error())
	}
}


func sidebarAppend (ts *gtk.TreeStore,iter *gtk.TreeIter, d string, v string) (*gtk.TreeIter, error) {
	i := ts.Append(iter)
	if err := ts.SetValue(i, 0, d); err != nil {
		return nil, fmt.Errorf("Unable to set sidebar %v, %v", d, err)
	}
	if err := ts.SetValue(i, 1, v); err != nil {
		return nil, fmt.Errorf("Unable to set sidebar %v, %v", v, err)
	}
	return i, nil
}

func sidebarSelected(tpath *gtk.TreePath) {
	nb, err := ui.GetNotebook("notebook")
	if err != nil {
		log.Println(err.Error())
	}

	t := tpath.String()
	switch t {
	case "0":
		nb.SetCurrentPage(0)  // logbook
		initLogBook()
		nb.ShowAll()
	case "1":
		nb.SetCurrentPage(1)  // log entry
	case "2":
		nb.SetCurrentPage(3)  // statistics		
	case "3":
		nb.SetCurrentPage(2)  // fcc
	case "4":
		nb.SetCurrentPage(4)  // settings
	}

}
