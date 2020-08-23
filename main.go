package main

import (
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	"mwlog/lshow"
)

const appID = "com.github.graynerd.mwlog"

func main() {
	// Create a new application.
	application, err := gtk.ApplicationNew(appID, glib.APPLICATION_FLAGS_NONE)
	if err != nil {
		log.Panic(err)
	}
	
	// Connect function to application startup event, this is not required.
	application.Connect("startup", func() {
		// os.Remove("mw.log.db")
		openDB()
	})
	
	// Connect function to application activate event
	application.Connect("activate", func() {
		log.Println("application activate")

		// Get the GtkBuilder UI definition in the glade file.
		builder, err := gtk.BuilderNewFromFile("main.ui")
		if err != nil {
			log.Panic(err)
		}

		ls, err := getListStore(builder, "liststore")
		if err != nil {
			log.Panic(err)
		}
		
		win, err := getWindow(builder, "main_window")
		if err != nil {
			log.Panic(err)
		}

		// Map the handlers to callback functions, and connect the signals to the Builder.
		 var signals = map[string]interface{}{
			"import_fcc":   func() { importFCC(win) },
			"load_ls": 		func() { lshow.LoadLS(sqldb, ls)},
		}
		builder.ConnectSignals(signals)

		// Show the Window and all of its components.
		win.ShowAll()
		application.AddWindow(win)
	})

	// Connect function to application shutdown event, this is not required.
	application.Connect("shutdown", func() {
		log.Println("application shutdown")
	})

	// Launch the application
	os.Exit(application.Run(os.Args))
}
