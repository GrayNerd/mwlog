package main

import (
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	"mwlog/db"
	"mwlog/lshow"
	"mwlog/ui"
)

const appID = "com.github.graynerd.mwlog"
const bFile = "main.ui"

func main() {
	// Create a new application.
	application, err := gtk.ApplicationNew(appID, glib.APPLICATION_FLAGS_NONE)
	if err != nil {
		log.Panic(err)
	}

	application.Connect("startup", func() {
		db.OpenDB()
		application.Connect("activate", func() { log.Println("application activate") })
		// application.Connect("shutdown", func() { log.Println("application shutdown") })

		// Get the GtkBuilder UI definition in the glade file.
		if err := ui.LoadBuilder(bFile); err != nil {
			log.Fatalf("Unable to load %v", bFile)
		}

		win, err := ui.GetWindow("main_window")
		if err != nil {
			log.Fatalln(err)
		}

		// Map the handlers to callback functions, and connect the signals to the Builder.
		var signals = map[string]interface{}{
			"on_import_fcc_activate":          func() { db.ImportFCC() },
			"on_display_fcc_activate":         func() { lshow.LoadLS() },
			"on_lg_date_focus_out_event":      func() { validateDate() },
			"on_lg_callsign_focus_out_event":  func() { validateCall() },
			"on_notebook_switch_page":         func(n *gtk.Notebook, p *gtk.Widget, pn int) { notebookSwitcher(pn) },
			"on_cx1_clicked":                  func(tv *gtk.TreeView, s *gtk.TreePath) { sidebarSelected(s) },
			"on_cancel_btn_clicked":           func() { clearLogEntry(); prefillLogEntry() },
			"on_ok_btn_clicked":               func() { saveLogEntry(0) },
			"on_logbook_tree_key_press_event": func(tv *gtk.TreeView, e *gdk.Event) { logbookDeleteSelected(tv, e) },
			"on_logbook_tree_row_activated":   func(tv *gtk.TreeView) { logbookEditSelected(tv)},
		}
		ui.ConnectSignals(signals)

		buildSidebar()
		initLogBook()
		win.ShowAll()
		application.AddWindow(win)
	})

	// Launch the application
	os.Exit(application.Run(os.Args))
}
