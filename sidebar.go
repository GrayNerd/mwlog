package main

import (
	"fmt"
	"log"

	"mwlog/ui"

	"github.com/gotk3/gotk3/gtk"
)

func buildSidebar() {
	treeStore:= ui.GetTreeStore("sidebar_ls")

	sidebarAppend(treeStore, nil, "Log Book", "0")
	sidebarAppend(treeStore, nil, "Add Logging", "-1")
	sidebarAppend(treeStore, nil, "Statistics", "1")
	sidebarAppend(treeStore, nil, "FCC Lookup", "3")

	iter :=sidebarAppend(treeStore, nil, "Settings", "4")
	sidebarAppend(treeStore, iter, "Languages", "4:0")
	sidebarAppend(treeStore, iter, "Equipment", "4:1")
}


func sidebarAppend (ts *gtk.TreeStore,iter *gtk.TreeIter, d string, v string) *gtk.TreeIter {
	i := ts.Append(iter)
	if err := ts.SetValue(i, 0, d); err != nil {
		log.Fatalln(fmt.Errorf("Unable to set sidebar %v, %v", d, err))
	}
	if err := ts.SetValue(i, 1, v); err != nil {
		log.Fatalln(fmt.Errorf("Unable to set sidebar %v, %v", v, err))
	}
	return i
}

func onSidebarMenuClicked(tpath *gtk.TreePath) {
	nb := ui.GetNotebook("notebook")

	t := tpath.String()
	switch t {
	case "0":
		nb.SetCurrentPage(0)  // logbook
		logbookLoad()
		nb.ShowAll()
	case "1":
		// nb.SetCurrentPage(0)  // log entry
		openLogging(0)
	case "2":
		nb.SetCurrentPage(3)  // statistics		
	case "3":
		nb.SetCurrentPage(2)  // fcc
	case "4":
		nb.SetCurrentPage(4)  // settings
	}

}
