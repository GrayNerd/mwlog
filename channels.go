package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"mwlog/db"
	"mwlog/ui"

	"github.com/gotk3/gotk3/gtk"
)

type chanTab struct {
}

func (ct *chanTab) showChannelsTab() {
	// buildFreqList()
	nb := ui.GetNotebook("notebook")
	nb.SetCurrentPage(1)
	nb.ShowAll()
}

func (ct *chanTab) buildFreqList() {
	ts := ui.GetTreeStore("chan_freq_list_store")
	for n := 530; n <= 1700; n += 10 {
		f := fmt.Sprintf("%4d", n)
		i := ts.Append(nil)
		if err := ts.SetValue(i, 0, f); err != nil {
			log.Fatalln(fmt.Errorf(err.Error()), "Unable to set frequency")
		}
	}
}

func (ct *chanTab) loadChannels(ts *gtk.TreeSelection) {
	model, iter, ok := ts.GetSelected()
	if !ok {
		log.Println("Unable to GetSelected in logbookEditSelected")
	}

	v, _ := model.(*gtk.TreeModel).GetValue(iter, 0)
	freq, err := v.GoValue()
	if err != nil {
		log.Println(err.Error())
	}
	f := strings.TrimSpace(freq.(string))
	ch, err := db.GetChannels(f)
	if err != nil {
		ui.GetLabel("chan_id").SetText("0")
		ui.GetLabel("chan_freq").SetText(f)
		ui.GetTextBuffer("chan_daytime_buffer").SetText("")
		ui.GetTextBuffer("chan_nighttime_buffer").SetText("")

	} else {
		ui.GetLabel("chan_id").SetText(strconv.Itoa(ch.ID))
		ui.GetLabel("chan_freq").SetText(ch.Frequency)
		// ui.GetEntry("chan_class").SetText(ch.Class)
		ui.GetTextBuffer("chan_daytime_buffer").SetText(ch.Daytime)
		ui.GetTextBuffer("chan_nighttime_buffer").SetText(ch.Nighttime)
	}
	loadChannelsLoggings(f)
}

func (ct *chanTab) saveChannels() {
	var ch db.Channel

	id, _ := ui.GetLabel("chan_id").GetText()
	ch.ID, _ = strconv.Atoi(id)

	ch.Frequency, _ = ui.GetLabel("chan_freq").GetText()

	day := ui.GetTextBuffer("chan_daytime_buffer")
	s, e := day.GetBounds()
	ch.Daytime, _ = day.GetText(s, e, false)

	night := ui.GetTextBuffer("chan_nighttime_buffer")
	s, e = night.GetBounds()
	ch.Nighttime, _ = night.GetText(s, e, false)

	if err := db.SaveChannels(&ch); err != nil {
		log.Println(err.Error())
	}
}

func loadChannelsLoggings(freq string) {
	ls := ui.GetTreeStore("chan_log_store")
	ls.Clear()

	rows, err := db.GetLoggingForFreq(freq)
	if err != nil {
		log.Println(err.Error())
	}
	defer rows.Close()

	var id int
	var station, city, state, country, firstHeard, lastHeard, timesHeard string
	var powerDay, powerNight, Distance, Bearing string
	for rows.Next() {
		err := rows.Scan(&id, &station, &city, &state, &country, &powerDay, &powerNight, 
										 &Distance, &Bearing, &firstHeard, &lastHeard, &timesHeard)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		var iter *gtk.TreeIter
		//ls.Append(iter)
		col := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
		var val []interface{}
		val = append(val, id, station, fmt.Sprintf("%s, %s %s", city, state, country),
			powerDay, Distance, Bearing, firstHeard, lastHeard, timesHeard)
		if err = ls.InsertWithValues(iter, nil, 0, col, val); err != nil {
			log.Println(err.Error())
		}
	}
}

func channelsTBSetup() {
	b := ui.GetToolButton("tb_edit")
	tv := ui.GetTreeView("logbook_tree")
	b.SetSensitive(false)
	if s, _ := tv.GetSelection(); s == nil {
		// set edit icon off
	}
	// turn edit icon on
	b.SetSensitive(true)
}
