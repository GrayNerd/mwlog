package ui

import (
	"fmt"
	"log"

	"github.com/gotk3/gotk3/gtk"
)

var builder *gtk.Builder

// LoadBuilder gets the GtkBuilder UI definition in the glade file.
func LoadBuilder(f string) {
	b, err := gtk.BuilderNewFromFile(f)
	if err != nil {
		log.Fatalf("unable to open builder file: %v", f)
	}
	builder = b
}
// ConnectSignals just calls the builder.ConnectSignals
func ConnectSignals(signals map[string]interface{}) {
	builder.ConnectSignals(signals)
}

// GetWindow returns a pointer to the named gtk.window
func GetWindow(item string) *gtk.Window {
	if obj, err := builder.GetObject(item); err == nil {
		if win, ok := obj.(*gtk.Window); ok {
			return win
		}
	}
	s := fmt.Sprintf("not a *gtk.Window: %v", item)
	log.Fatalln(s)
	return nil
}

// GetDialog returns a pointer to the named gtk.dialog
func GetDialog(item string) *gtk.Dialog {
	if obj, err := builder.GetObject(item); err == nil {
		if dlg, ok := obj.(*gtk.Dialog); ok {
			return dlg
		}
	}
	s := fmt.Sprintf("not a *gtk.Dialog: %v", item)
	log.Fatalln(s)
	return nil
}

//GetListStore returns a pointer to the named gtk.liststore
func GetListStore(item string) *gtk.ListStore {
	if obj, err := builder.GetObject(item); err == nil {
		if ls, ok := obj.(*gtk.ListStore); ok {
			return ls
		}
	}
	s := fmt.Sprintf("not a *gtk.ListStore: %v", item)
	log.Fatalln(s)
	return nil
}

// GetTreeView returns a ponter to the named gtk.treeview
func GetTreeView(item string) *gtk.TreeView {
	if obj, err := builder.GetObject(item); err == nil {
		if tv, ok := obj.(*gtk.TreeView); ok {
			return tv
		}
	}
	s := fmt.Sprintf("not a *gtk.TreeView: %v", item)
	log.Fatalln(s)
	return nil
}

// GetCheckButton returns a pointer to the named gtk.checkbutton
func GetCheckButton(item string) *gtk.CheckButton {
	if obj, err := builder.GetObject(item); err == nil {
		if cb, ok := obj.(*gtk.CheckButton); ok {
			return cb
		}
	}
	s := fmt.Sprintf("not a *gtk.GetCheckButton: %v", item)
	log.Fatalln(s)
	return nil
}

// GetEntry returns a pointer to the named gtk.entry
func GetEntry(item string) *gtk.Entry {
	if obj, err := builder.GetObject(item); err == nil {
		if e, ok := obj.(*gtk.Entry); ok {
			return e
		}
	}
	s := fmt.Sprintf("not a *gtk.Entry: %v", item)
	log.Fatalln(s)
	return nil
}

// GetTreeStore returns a pointer to the named gtk.TreeStore
func GetTreeStore(item string) *gtk.TreeStore {
	if obj, err := builder.GetObject(item); err == nil {
		if e, ok := obj.(*gtk.TreeStore); ok {
			return e
		}
	}
	s := fmt.Sprintf("not a *gtk.TreeStore: %v", item)
	log.Fatalln(s)
	return nil
}

// GetTreeViewColumn returns a pointer to the named gtk.TreeViewColumn
func GetTreeViewColumn(item string) *gtk.TreeViewColumn {
	if obj, err := builder.GetObject(item); err == nil {
		if e, ok := obj.(*gtk.TreeViewColumn); ok {
			return e
		}
	}
	s := fmt.Sprintf("not a *gtk.TreeViewColumn: %v", item)
	log.Fatalln(s)
	return nil
}

// GetCellRenderer returns a pointer to the named gtk.GetCellRenderer
func GetCellRenderer(item string) *gtk.CellRenderer {
	if obj, err := builder.GetObject(item); err == nil {
		if e, ok := obj.(*gtk.CellRenderer); ok {
			return e
		}
	}
	s := fmt.Sprintf("not a *gtk.CellRenderer: %v", item)
	log.Fatalln(s)
	return nil
}

// GetCellRendererText returns a pointer to the named gtk.GetCellRendererText
func GetCellRendererText(item string) *gtk.CellRendererText {
	if obj, err := builder.GetObject(item); err == nil {
		if e, ok := obj.(*gtk.CellRendererText); ok {
			return e
		}
	}
	s := fmt.Sprintf("not a *gtk.CellRendererText: %v", item)
	log.Fatalln(s)
	return nil
}

// GetTreeSelection returns a pointer to the named gtk.TreeSelection
func GetTreeSelection(item string) *gtk.TreeSelection {
	if obj, err := builder.GetObject(item); err == nil {
		if e, ok := obj.(*gtk.TreeSelection); ok {
			return e
		}
	}
	s := fmt.Sprintf("not a *gtk.TreeSelection: %v", item)
	log.Fatalln(s)
	return nil
}

// GetNotebook returns a pointer to the named gtk.Notebook
func GetNotebook(item string) *gtk.Notebook {
	if obj, err := builder.GetObject(item); err == nil {
		if e, ok := obj.(*gtk.Notebook); ok {
			return e
		}
	}
	s := fmt.Sprintf("not a *gtk.Notebook: %v", item)
	log.Fatalln(s)
	return nil
}

// GetComboBox returns a pointer to the named gtk.ComboBox
func GetComboBox(item string) *gtk.ComboBox {
	if obj, err := builder.GetObject(item); err == nil {
		if e, ok := obj.(*gtk.ComboBox); ok {
			return e
		}
	}
	s := fmt.Sprintf("not a *gtk.ComboBox: %v", item)
	log.Fatalln(s)
	return nil
}

// GetButton returns a pointer to the named gtk.ComboBox
func GetButton(item string) *gtk.Button {
	if obj, err := builder.GetObject(item); err == nil {
		if e, ok := obj.(*gtk.Button); ok {
			return e
		}
	}
	s := fmt.Sprintf("not a *gtk.Button: %v", item)
	log.Fatalln(s)
	return nil
}
// GetTextView returns a pointer to the named gtk.TextView
func GetTextView(item string) *gtk.TextView{
	if obj, err := builder.GetObject(item); err == nil {
		if e, ok := obj.(*gtk.TextView); ok {
			return e
		}
	}
	s := fmt.Sprintf("not a *gtk.TextView: %v", item)
	log.Fatalln(s)
	return nil
}

// GetTextBuffer returns a pointer to the named gtk.TextBuffer
func GetTextBuffer(item string) *gtk.TextBuffer {
	if obj, err := builder.GetObject(item); err == nil {
		if e, ok := obj.(*gtk.TextBuffer); ok {
			return e
		}
	}
	s := fmt.Sprintf("not a *gtk.TextBuffer: %v", item)
	log.Fatalln(s)
	return nil
}
// GetLabel returns a pointer to the named gtk.Label
func GetLabel(item string) *gtk.Label {
	if obj, err := builder.GetObject(item); err == nil {
		if e, ok := obj.(*gtk.Label); ok {
			return e
		}
	}
	s := fmt.Sprintf("not a *gtk.Label: %v", item)
	log.Fatalln(s)
	return nil
}
