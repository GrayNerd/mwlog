package main

import (
	"errors"
	"fmt"

	"github.com/gotk3/gotk3/gtk"
)


func getWindow(builder *gtk.Builder, item string) (*gtk.Window, error) {
	var err error
	if obj, err := builder.GetObject(item); err == nil {
		if win, ok := obj.(*gtk.Window); ok {
			return win, nil
		}
	}
	return nil, fmt.Errorf("not a *gtk.Window: %s", err)
}

func getDialog(builder *gtk.Builder, item string) (*gtk.Dialog, error) {
	if obj, err := builder.GetObject(item); err == nil {
		if dlg, ok := obj.(*gtk.Dialog); ok {
			return dlg, nil
		}
	}
	return nil, errors.New("not a *gtk.Dialog")
}

func getListStore(builder *gtk.Builder, item string) (*gtk.ListStore, error) {
	if obj, err := builder.GetObject(item); err == nil {
		if ls, ok := obj.(*gtk.ListStore); ok {
			return ls, nil
		}
	}
	return nil, errors.New("not a *gtk.ListStore")
}

func getTreeView(builder *gtk.Builder, item string) (*gtk.TreeView, error) {
	if obj, err := builder.GetObject(item); err == nil {
		if tv, ok := obj.(*gtk.TreeView); ok {
			return tv, nil
		}
	}
	return nil, errors.New("not a *gtk.TreeView")
}

// func getCombo(builder *gtk.Builder, item string) (*gtk.CellRendererCombo, error) {
// 	if obj, err := builder.GetObject(item); err == nil {
// 		if combo, ok := obj.(*gtk.CellRendererCombo); ok {
// 			return combo, nil
// 		}
// 	}
// 	return nil, errors.New("not a *gtk.CellRendererCombo")
// }

func getCheckButton(builder *gtk.Builder, item string) (*gtk.CheckButton, error) {
	if obj, err := builder.GetObject(item); err == nil {
		if cb, ok := obj.(*gtk.CheckButton); ok {
			return cb, nil
		}
	}
	return nil, errors.New("not a *gtk.CheckButton")
}


