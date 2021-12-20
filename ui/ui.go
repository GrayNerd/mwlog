package ui

import (
	"log"

	// "github.com/gotk3/gotk3/glib"
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
	obj, err := builder.GetObject(item)
	if err == nil {
		if win, ok := obj.(*gtk.Window); !ok {
			log.Fatalf("%v is not a gtk.Window", item)
		} else {
			return win
		}
	}
	log.Fatalf("%v not found in builder file: %v", item, err.Error())
	return nil
}

// GetDialog returns a pointer to the named gtk.dialog
func GetDialog(item string) *gtk.Dialog {
	if obj, err := builder.GetObject(item); err == nil {
		if dlg, ok := obj.(*gtk.Dialog); ok {
			return dlg
		}
		log.Fatalf("%v is not a *gtk.Dialog", item)
	} else {
		log.Fatalf("%v not found in builder file: %v", item, err.Error())
	}
	return nil
}

//GetListStore returns a pointer to the named gtk.ListStore
func GetListStore(item string) *gtk.ListStore {
	if obj, err := builder.GetObject(item); err == nil {
		if ls, ok := obj.(*gtk.ListStore); ok {
			return ls
		}
		log.Fatalf("%v is not a *gtk.ListStore", item)
	} else {
		log.Fatalf("%v not found in builder file: %v", item, err.Error())
	}
	return nil
}

// GetTreeView returns a pointer to the named gtk.TreeView
func GetTreeView(item string) *gtk.TreeView {
	if obj, err := builder.GetObject(item); err == nil {
		if tv, ok := obj.(*gtk.TreeView); ok {
			return tv
		}
		log.Fatalf("%v is not a *gtk.TreeView", item)
	} else {
		log.Fatalf("%v not found in builder file: %v", item, err.Error())
	}
	return nil
}

// GetCheckButton returns a pointer to the named gtk.checkbutton
func GetCheckButton(item string) *gtk.CheckButton {
	if obj, err := builder.GetObject(item); err == nil {
		if cb, ok := obj.(*gtk.CheckButton); ok {
			return cb
		}
		log.Fatalf("%v is not a *gtk.CheckButton", item)
	} else {
		log.Fatalf("%v not found in builder file: %v", item, err.Error())
	}
	return nil
}

// GetEntry returns a pointer to the named gtk.entry
func GetEntry(item string) *gtk.Entry {
	if obj, err := builder.GetObject(item); err == nil {
		if e, ok := obj.(*gtk.Entry); ok {
			return e
		}
		log.Fatalf("%v is not a *gtk.Entry", item)
	} else {
		log.Fatalf("%v not found in builder file: %v", item, err.Error())
	}
	return nil
}

// GetTreeStore returns a pointer to the named gtk.TreeStore
func GetTreeStore(item string) *gtk.TreeStore {
	if obj, err := builder.GetObject(item); err == nil {
		if e, ok := obj.(*gtk.TreeStore); ok {
			return e
		}
		log.Fatalf("%v is not a *gtk.TreeStore", item)
	} else {
		log.Fatalf("%v not found in builder file: %v", item, err.Error())
	}
	return nil
}

// GetTreeViewColumn returns a pointer to the named gtk.TreeViewColumn
func GetTreeViewColumn(item string) *gtk.TreeViewColumn {
	if obj, err := builder.GetObject(item); err == nil {
		if e, ok := obj.(*gtk.TreeViewColumn); ok {
			return e
		}
		log.Fatalf("not a *gtk.TreeViewColumn: %v", item)
	} else {
		log.Fatalf("%v not found in builder file: %v", item, err.Error())
	}
	return nil
}

// GetCellRenderer returns a pointer to the named gtk.CellRenderer
func GetCellRenderer(item string) *gtk.CellRenderer {
	if obj, err := builder.GetObject(item); err == nil {
		if e, ok := obj.(*gtk.CellRenderer); ok {
			return e
		}
		log.Fatalf("%v is not a *gtk.CellRenderer", item)
	} else {
		log.Fatalf("%v not found in builder file: %v", item, err.Error())
	}
	return nil
}

// GetCellRendererText returns a pointer to the named gtk.CellRendererText
func GetCellRendererText(item string) *gtk.CellRendererText {
	if obj, err := builder.GetObject(item); err == nil {
		if e, ok := obj.(*gtk.CellRendererText); ok {
			return e
		}
		log.Fatalf("%v is not a *gtk.CellRendererText", item)
	} else {
		log.Fatalf("%v not found in builder file: %v", item, err.Error())
	}
	return nil
}

// GetTreeSelection returns a pointer to the named gtk.TreeSelection
func GetTreeSelection(item string) *gtk.TreeSelection {
	if obj, err := builder.GetObject(item); err == nil {
		if e, ok := obj.(*gtk.TreeSelection); ok {
			return e
		}
		log.Fatalf("%v is not a *gtk.TreeSelection", item)
	} else {
		log.Fatalf("%v not found in builder file: %v", item, err.Error())
	}
	return nil
}

// GetTreeModelSort returns a pointer to the named gtk.TreeModelSort
func GetTreeModelSort(item string) *gtk.TreeModelSort {
	if obj, err := builder.GetObject(item); err == nil {
		if e, ok := obj.(*gtk.TreeModelSort); ok {
			return e
		}
		log.Fatalf("%v is not a *gtk.TreeModelSort", item)
	} else {
		log.Fatalf("%v not found in builder file: %v", item, err.Error())
	}
	return nil
}

// GetNotebook returns a pointer to the named gtk.Notebook
func GetNotebook(item string) *gtk.Notebook {
	if obj, err := builder.GetObject(item); err == nil {
		if e, ok := obj.(*gtk.Notebook); ok {
			return e
		}
		log.Fatalf("%v is not a *gtk.Notebook", item)
	} else {
		log.Fatalf("%v not found in builder file: %v", item, err.Error())
	}
	return nil
}

// GetComboBox returns a pointer to the named gtk.ComboBox
func GetComboBox(item string) *gtk.ComboBox {
	if obj, err := builder.GetObject(item); err == nil {
		if e, ok := obj.(*gtk.ComboBox); ok {
			return e
		}
		log.Fatalf("%v is not a *gtk.ComboBox", item)
	} else {
		log.Fatalf("%v not found in builder file: %v", item, err.Error())
	}
	return nil
}

// GetButton returns a pointer to the named gtk.Button
func GetButton(item string) *gtk.Button {
	if obj, err := builder.GetObject(item); err == nil {
		if e, ok := obj.(*gtk.Button); ok {
			return e
		}
		log.Fatalf(" %v is not a *gtk.Button", item)
	} else {
		log.Fatalf("%v not found in builder file: %v", item, err.Error())
	}
	return nil
}

// GetTextView returns a pointer to the named gtk.TextView
func GetTextView(item string) *gtk.TextView {
	if obj, err := builder.GetObject(item); err == nil {
		if e, ok := obj.(*gtk.TextView); ok {
			return e
		}
		log.Fatalf("%v is not a *gtk.TextView", item)
	} else {
		log.Fatalf("%v not found in builder file: %v", item, err.Error())
	}
	return nil
}

// GetTextBuffer returns a pointer to the named gtk.TextBuffer
func GetTextBuffer(item string) *gtk.TextBuffer {
	if obj, err := builder.GetObject(item); err == nil {
		if e, ok := obj.(*gtk.TextBuffer); ok {
			return e
		}
		log.Fatalf("%v is not a *gtk.TextBuffer", item)
	} else {
		log.Fatalf("%v not found in builder file: %v", item, err.Error())
	}
	return nil
}

// GetLabel returns a pointer to the named gtk.Label
func GetLabel(item string) *gtk.Label {
	if obj, err := builder.GetObject(item); err == nil {
		if e, ok := obj.(*gtk.Label); ok {
			return e
		}
		log.Fatalf("%v is not a gtk.Label", item)
	} else {
		log.Fatalf("%v not found in builder file: %v", item, err.Error())
	}
	return nil
}

// GetListBox returns a pointer to the named gtk.ListBox
func GetListBox(item string) *gtk.ListBox {
	if obj, err := builder.GetObject(item); err == nil {
		if e, ok := obj.(*gtk.ListBox); ok {
			return e
		}
		log.Fatalf("%v is not a gtk.ListBox", item)
	} else {
		log.Fatalf("%v not found in builder file: %v", item, err.Error())
	}
	return nil
}

// GetImage returns a pointer to the named gtk.Image
func GetImage(item string) *gtk.Image {
	if obj, err := builder.GetObject(item); err == nil {
		if e, ok := obj.(*gtk.Image); ok {
			return e
		}
		log.Fatalf("%v is not a gtk.Image", item)
	} else {
		log.Fatalf("%v not found in builder file: %v", item, err.Error())
	}
	return nil
}

// GetViewport returns a pointer to the named gtk.Viewport
func GetViewport(item string) *gtk.Viewport {
	if obj, err := builder.GetObject(item); err == nil {
		if e, ok := obj.(*gtk.Viewport); ok {
			return e
		}
		log.Fatalf("%v is not a gtk.Viewport", item)
	} else {
		log.Fatalf("%v not found in builder file: %v", item, err.Error())
	}
	return nil
}

// GetMenu returns a pointer to the named gtk.Menu
func GetMenu(item string) *gtk.Menu {
	if obj, err := builder.GetObject(item); err == nil {
		if e, ok := obj.(*gtk.Menu); ok {
			return e
		}
		log.Fatalf("%v is not a gtk.Menu", item)
	} else {
		log.Fatalf("%v not found in builder file: %v", item, err.Error())
	}
	return nil
}
