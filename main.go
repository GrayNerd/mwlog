package main

import (
	// "fmt"
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
	log.Printf("Application Starting")

	gtk.Init(nil)
	// Create a new application.
	application, err := gtk.ApplicationNew(appID, glib.APPLICATION_FLAGS_NONE)
	if err != nil {
		log.Panic(err)
	}

	// if set, err := gtk.SettingsGetDefault(); err == nil {
	// 	if set.SetProperty("gtk-theme-name", "Mint-Y-Darker-Grey") != nil {
	// 		log.Println(err.Error())
	// 	}
	// 	if set.SetProperty("gtk-application-prefer-dark-theme", true) != nil {
	// 				log.Println(err.Error())
	// 	}
	// } else {
	// 			log.Println(err.Error())
	// }

	application.Connect("startup", func() {
		var ch chanTab
		var mt mapsTab
		var lt mwListTab
		var le logging

		notebookSwitcher := func(pn int) {
			switch pn {
			case 0: // logbook
				logbookTBSetup()
			case 1: // channels
				ui.GetToolButton("tb_edit").SetSensitive(false)
			case 2: // mw list
				ui.GetToolButton("tb_edit").SetSensitive(false)
				lt.showMWListTab()
			case 3: // maps
				ui.GetToolButton("tb_edit").SetSensitive(false)
				mt.showMapsTab()
			}
		}

		db.OpenDB()
		application.Connect("activate", func() {}) // log.Println("application activate") })
		// application.Connect("shutdown", func() { log.Println("application shutdown") })

		// Get the GtkBuilder UI definition in the glade file.
		ui.LoadBuilder(bFile)

		win := ui.GetWindow("main_window")

		ch.buildFreqList()

		// Map the handlers to callback functions, and connect the signals to the Builder.
		var signals = map[string]interface{}{
			"on_import_mwlist_activate":  func() { db.ImportMWList() },
			"on_display_mwlist_activate": func() { lshow.LoadLS() },

			// *** Logging Window ***
			"on_logging_date_focus_out_event":    func(e *gtk.Entry) { le.validateDate(e) },
			"on_logging_time_focus_out_event":    func(e *gtk.Entry) { le.validateTime(e) },
			"on_logging_station_focus_out_event": func(e *gtk.Entry) { le.validateCall(e) },
			"on_logging_cancel_button_clicked":   func(_ *gtk.Button) { le.window.Close() },
			// "on_receiver_completion_match_selected": func(m *gtk.ListStore, i *gtk.TreeIter) { ms(m, i) },

			"on_notebook_switch_page": func(_ *gtk.Notebook, _ *gtk.Widget, pn int) { notebookSwitcher(pn) },

			// *** Menu Items ***
			"on_menu_logbook_clicked":      func() { showLogbook() },
			"on_menu_add_logging_clicked":  func() { le.open(0) },
			"on_menu_edit_logging_clicked": func() { le.edit() },

			// *** Channels Tab ***
			"on_chan_freq_sel_changed": func(ts *gtk.TreeSelection) { ch.loadChannels(ts) },
			"on_chan_save_clicked":     func() { ch.saveChannels() },

			// *** Logbook Tab ***
			"on_lb_tree_selection_changed":       func(ts *gtk.TreeSelection) { displayRow(ts) },
			"on_logbook_tree_key_press_event":    func(tv *gtk.TreeView, e *gdk.Event) { onLogbookTreeKeyPressEvent(tv, e) },
			"on_logbook_tree_button_press_event": func(tv *gtk.TreeView, e *gdk.Event) { onLogbookTreeButtonPressEvent(tv, e, le) },

			// *** Logbook popup menu ***
			"on_pu_edit_activate":   func() { le.edit() },
			"on_pu_delete_activate": func() { onLogbookDelete() },

			// *** MWList Tab ***
			"on_mwlist_tv_button_press_event": func(tv *gtk.TreeView, e *gdk.Event) { onMWListTreeButtonPressEvent(tv, e) },

			// *** Maps Tab ***
			"on_maps_viewport_size_allocate":        func() { mt.mapResize() },
			"on_maps_viewport_scroll_event":         func(_ *glib.Object, e *gdk.Event) { mt.zoom(e) },
			"on_maps_viewport_button_release_event": func(_ *glib.Object, e *gdk.Event) { mt.click(e) },

			// "click":     func() { log.Println("I was clicked") },

			// "changed": func() bool {
			// 	log.Println("changed")
			// 	return false
			// },
		}
		ui.ConnectSignals(signals)

		loadCSS()
		loadSelections()
		loadLogbook()

		win.ShowAll()

		application.AddWindow(win)
	})

	// Launch the application
	os.Exit(application.Run(os.Args))
}

// func ms(m *gtk.ListStore, i *gtk.TreeIter) {
// 	rcvr, err := m.GetValue(i, 1)
// 	if err != nil {
// 		log.Println(err.Error())
// 		return
// 	}
// 	str, err := rcvr.GetString()
// 	if err != nil {
// 		log.Println(err.Error())
// 		return
// 	}
// 	ui.GetEntry("logging_receiver_entry").SetText(str)
// }

func loadCSS() {
	var cssProv *gtk.CssProvider
	var screen *gdk.Screen
	var err error

	if cssProv, err = gtk.CssProviderNew(); err != nil {
		log.Panic(err)
	}
	if err = cssProv.LoadFromPath("mwlog.css"); err != nil {
		log.Panic(err)
	}
	if screen, err = gdk.ScreenGetDefault(); err != nil {
		log.Panic(err)
	}
	gtk.AddProviderForScreen(screen, cssProv, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
}

func loadSelections() {
	ls := ui.GetListStore("format_ls")
	ls.Clear()

	rows := db.GetAllFormats()
	defer rows.Close()

	var id int
	var value string
	for rows.Next() {
		err := rows.Scan(&id, &value)
		if err != nil {
			log.Println(err.Error())
		}
		var iter *gtk.TreeIter
		col := []int{0, 1}
		var val []interface{}
		val = append(val, id, value)
		if err := ls.InsertWithValues(iter, 0, col, val); err != nil {
			log.Println(err.Error())
		}
	}
}

func displayRow(ts *gtk.TreeSelection) {
	model, iter, ok := ts.GetSelected()
	if ok {
		path, err := model.(*gtk.TreeModel).GetPath(iter)
		if err != nil {
			log.Println(err.Error())
		}
		tv := ui.GetTreeView("logbook_tree")
		tv.ScrollToCell(path, nil, false, 0, 0)
	}
}
