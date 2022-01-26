package main

import (
	"testing"

	"github.com/gotk3/gotk3/gdk"
)

func Test_mapsTab_showMapsTab(t *testing.T) {
	type fields struct {
		origBuf *gdk.Pixbuf
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mt := &mapsTab{
				origBuf: tt.fields.origBuf,
			}
			_=t
			mt.showMapsTab()
		})
	}
}
