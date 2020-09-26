package ui

import (
	"fmt"
	"log"

	"github.com/gotk3/gotk3/gtk"
)

var builder *gtk.Builder

// LoadBuilder gets the GtkBuilder UI definition in the glade file.
func LoadBuilder(f string) error {
	b, err := gtk.BuilderNewFromFile(f)
	if err != nil {
		return fmt.Errorf("unable to open %v", f)
	}
	builder = b
	return nil
}
// ConnectSignals just calls the builder.ConnectSignals
func ConnectSignals(signals map[string]interface{}) {
	builder.ConnectSignals(signals)
}

// GetWindow returns a pointer to the named gtk.window
func GetWindow(item string) (*gtk.Window, error) {
	if obj, err := builder.GetObject(item); err == nil {
		if win, ok := obj.(*gtk.Window); ok {
			return win, nil
		}
	}
	s := fmt.Sprintf("not a *gtk.Window: %v", item)
	log.Println(s)
	return nil, fmt.Errorf(s)
}

// GetDialog returns a pointer to the named gtk.dialog
func GetDialog(item string) (*gtk.Dialog, error) {
	if obj, err := builder.GetObject(item); err == nil {
		if dlg, ok := obj.(*gtk.Dialog); ok {
			return dlg, nil
		}
	}
	s := fmt.Sprintf("not a *gtk.Dialog: %v", item)
	log.Println(s)
	return nil, fmt.Errorf(s)
}

//GetListStore returns a pointer to the named gtk.liststore
func GetListStore(item string) (*gtk.ListStore, error) {
	if obj, err := builder.GetObject(item); err == nil {
		if ls, ok := obj.(*gtk.ListStore); ok {
			return ls, nil
		}
	}
	s := fmt.Sprintf("not a *gtk.ListStore: %v", item)
	log.Println(s)
	return nil, fmt.Errorf(s)
}

// GetTreeView returns a ponter to the named gtk.treeview
func GetTreeView(item string) (*gtk.TreeView, error) {
	if obj, err := builder.GetObject(item); err == nil {
		if tv, ok := obj.(*gtk.TreeView); ok {
			return tv, nil
		}
	}
	s := fmt.Sprintf("not a *gtk.TreeView: %v", item)
	log.Println(s)
	return nil, fmt.Errorf(s)
}

// GetCheckButton returns a pointer to the named gtk.checkbutton
func GetCheckButton(item string) (*gtk.CheckButton, error) {
	if obj, err := builder.GetObject(item); err == nil {
		if cb, ok := obj.(*gtk.CheckButton); ok {
			return cb, nil
		}
	}
	s := fmt.Sprintf("not a *gtk.GetCheckButton: %v", item)
	log.Println(s)
	return nil, fmt.Errorf(s)
}

// GetEntry returns a pointer to the named gtk.entry
func GetEntry(item string) (*gtk.Entry, error) {
	if obj, err := builder.GetObject(item); err == nil {
		if e, ok := obj.(*gtk.Entry); ok {
			return e, nil
		}
	}
	s := fmt.Sprintf("not a *gtk.Entry: %v", item)
	log.Println(s)
	return nil, fmt.Errorf(s)
}

// GetTreeStore returns a pointer to the named gtk.TreeStore
func GetTreeStore(item string) (*gtk.TreeStore, error) {
	if obj, err := builder.GetObject(item); err == nil {
		if e, ok := obj.(*gtk.TreeStore); ok {
			return e, nil
		}
	}
	s := fmt.Sprintf("not a *gtk.TreeStore: %v", item)
	log.Println(s)
	return nil, fmt.Errorf(s)
}

// GetTreeViewColumn returns a pointer to the named gtk.TreeViewColumn
func GetTreeViewColumn(item string) (*gtk.TreeViewColumn, error) {
	if obj, err := builder.GetObject(item); err == nil {
		if e, ok := obj.(*gtk.TreeViewColumn); ok {
			return e, nil
		}
	}
	s := fmt.Sprintf("not a *gtk.TreeViewColumn: %v", item)
	log.Println(s)
	return nil, fmt.Errorf(s)
}

// GetCellRenderer returns a pointer to the named gtk.GetCellRenderer
func GetCellRenderer(item string) (*gtk.CellRenderer, error) {
	if obj, err := builder.GetObject(item); err == nil {
		if e, ok := obj.(*gtk.CellRenderer); ok {
			return e, nil
		}
	}
	s := fmt.Sprintf("not a *gtk.CellRenderer: %v", item)
	log.Println(s)
	return nil, fmt.Errorf(s)
}

// GetCellRendererText returns a pointer to the named gtk.GetCellRendererText
func GetCellRendererText(item string) (*gtk.CellRendererText, error) {
	if obj, err := builder.GetObject(item); err == nil {
		if e, ok := obj.(*gtk.CellRendererText); ok {
			return e, nil
		}
	}
	s := fmt.Sprintf("not a *gtk.CellRendererText: %v", item)
	log.Println(s)
	return nil, fmt.Errorf(s)
}

// GetTreeSelection returns a pointer to the named gtk.TreeSelection
func GetTreeSelection(item string) (*gtk.TreeSelection, error) {
	if obj, err := builder.GetObject(item); err == nil {
		if e, ok := obj.(*gtk.TreeSelection); ok {
			return e, nil
		}
	}
	s := fmt.Sprintf("not a *gtk.TreeSelection: %v", item)
	log.Println(s)
	return nil, fmt.Errorf(s)
}

// GetNotebook returns a pointer to the named gtk.Notebook
func GetNotebook(item string) (*gtk.Notebook, error) {
	if obj, err := builder.GetObject(item); err == nil {
		if e, ok := obj.(*gtk.Notebook); ok {
			return e, nil
		}
	}
	s := fmt.Sprintf("not a *gtk.Notebook: %v", item)
	log.Println(s)
	return nil, fmt.Errorf(s)
}

// GetComboBox returns a pointer to the named gtk.ComboBox
func GetComboBox(item string) (*gtk.ComboBox, error) {
	if obj, err := builder.GetObject(item); err == nil {
		if e, ok := obj.(*gtk.ComboBox); ok {
			return e, nil
		}
	}
	s := fmt.Sprintf("not a *gtk.ComboBox: %v", item)
	log.Println(s)
	return nil, fmt.Errorf(s)
}

// GetButton returns a pointer to the named gtk.ComboBox
func GetButton(item string) (*gtk.Button, error) {
	if obj, err := builder.GetObject(item); err == nil {
		if e, ok := obj.(*gtk.Button); ok {
			return e, nil
		}
	}
	s := fmt.Sprintf("not a *gtk.Button: %v", item)
	log.Println(s)
	return nil, fmt.Errorf(s)
}
// GetTextView returns a pointer to the named gtk.TextView
func GetTextView(item string) (*gtk.TextView, error) {
	if obj, err := builder.GetObject(item); err == nil {
		if e, ok := obj.(*gtk.TextView); ok {
			return e, nil
		}
	}
	s := fmt.Sprintf("not a *gtk.TextView: %v", item)
	log.Println(s)
	return nil, fmt.Errorf(s)
}

// GetTextBuffer returns a pointer to the named gtk.TextBuffer
func GetTextBuffer(item string) (*gtk.TextBuffer, error) {
	if obj, err := builder.GetObject(item); err == nil {
		if e, ok := obj.(*gtk.TextBuffer); ok {
			return e, nil
		}
	}
	s := fmt.Sprintf("not a *gtk.TextBuffer: %v", item)
	log.Println(s)
	return nil, fmt.Errorf(s)
}
