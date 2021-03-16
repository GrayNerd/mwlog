package main

import (
	"mwlog/ui"
	"mwlog/db"

	"image/color"

	"github.com/gotk3/gotk3/gdk"
	// "github.com/gotk3/gotk3/gtk"

	"github.com/flopp/go-staticmaps"
	"github.com/fogleman/gg"
	"github.com/golang/geo/s2"
)

type mapsTab struct {
	origBuf *gdk.Pixbuf
}

const mapHeight = 1920.0
const mapWidth = 1080.0
const mapRatio = mapHeight / mapWidth


func (mt *mapsTab) showMapsTab() {
	ctx := sm.NewContext()
	ctx.SetSize(int(mapHeight), int(mapWidth))

	rows := db.GetLoggingLocations()
	defer rows.Close()

	var station string
	var lat, long float64
	for rows.Next() {
		rows.Scan(&station, &lat, &long)
		mark := sm.NewMarker(s2.LatLngFromDegrees(lat, long), color.RGBA{0xff, 0, 0, 0xff}, 8.0)
		mark.Label = station + "     "
		mark.LabelColor = color.Black
		ctx.AddMarker(mark)
	}

	img, err := ctx.Render()
	if err != nil {
		panic(err)
	}

	if err := gg.SaveJPG("mwlog-locations.jpg", img, 90); err != nil {
		panic(err)
	}

	image := ui.GetImage("maps_locations")
	image.SetFromFile("mwlog-locations.jpg")
	mt.origBuf = image.GetPixbuf()
	mt.mapResize()
}

func (mt *mapsTab) mapResize() {
	view := ui.GetViewport("maps_viewport")
	w := float64(view.GetAllocatedWidth())
	h := float64(view.GetAllocatedHeight())
	if mapRatio > w/h {
		h = w / mapRatio
	} else {
		w = h * mapRatio
	}

	image := ui.GetImage("maps_locations")
	pixbuf, _ := mt.origBuf.ScaleSimple(int(w), int(h), gdk.INTERP_BILINEAR)
	image.SetFromPixbuf(pixbuf)
}