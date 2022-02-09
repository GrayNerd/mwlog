package main

import (
	"fmt"
	"strconv"
	"strings"

	"log"
	"mwlog/db"
	"mwlog/ui"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

type mwListTab struct {
}


func (lt *mwListTab) showMWListTab() {
	tsDetails := ui.GetTreeStore("mwlist_store")
	tsDetails.Clear()

	// var sortFreq gtk.TreeIterCompareFunc = func(model *gtk.TreeModel, iter1, iter2 *gtk.TreeIter) int {
	// 	cv1, _ := model.GetValue(iter1, 0)
	// 	cv2, _ := model.GetValue(iter2, 0)
	// 	gv1,_  := cv1.GoValue()
	// 	gv2, _ := cv2.GoValue()
	// 	v1,_ := strconv.Atoi(gv1.(string))
	// 	v2, _ := strconv.Atoi(gv2.(string))
	// 	if v1 < v2 {
	// 		return -1
	// 	}
	// 	if v1 > v2 {
	// 		return 1
	// 	}
	// 	return 0
	// }
	// var sortPower gtk.TreeIterCompareFunc = func(model *gtk.TreeModel, iter1, iter2 *gtk.TreeIter) int {
	// 	cv1, _ := model.GetValue(iter1, 5)
	// 	cv2, _ := model.GetValue(iter2, 5)
	// 	gv1, _ := cv1.GoValue()
	// 	gv2, _ := cv2.GoValue()
	// 	v1, _ := strconv.ParseFloat(gv1.(string), 32)
	// 	v2, _ := strconv.ParseFloat(gv2.(string), 32)
	// 	if v1 < v2 {
	// 		return -1
	// 	}
	// 	if v1 > v2 {
	// 		return 1
	// 	}
	// 	return 0
	// }
	// var sortDistance gtk.TreeIterCompareFunc = func(model *gtk.TreeModel, iter1, iter2 *gtk.TreeIter) int {
	// 	cv1, _ := model.GetValue(iter1, 6)
	// 	cv2, _ := model.GetValue(iter2, 6)
	// 	gv1, _ := cv1.GoValue()
	// 	gv2, _ := cv2.GoValue()
	// 	s1 := gv1.(string)
	// 	s2 := gv2.(string)
	// 	fmt.Printf("%s - %s", s1, s2)
	// 	v1, _ := strconv.Atoi(strings.Trim(s1, " "))
	// 	v2, _ := strconv.Atoi(strings.Trim(s2, " "))
	// 	if s1 == "" || s2 == "" {
	// 		return 0
	// 	}
	// 	if v1 < v2 {
	// 		return -1
	// 	}
	// 	if v1 > v2 {
	// 		return 1
	// 	}
	// 	return 0
	// }
	// tsDetails.SetSortFunc(5, sortPower)
	// tsDetails.SetSortFunc(6, sortDistance)

	rows := db.GetAllMWList()
	defer rows.Close()
	var station, frequency, city, state, country, day, night, power string
	var distance, bearing, fPower float64
	var iFreq int16
	oldfreq := ""
	var parent *gtk.TreeIter
	for rows.Next() {
		rows.Scan(&frequency, &station, &city, &state, &country, &day, &night, &distance, &bearing)
		f, _ := strconv.Atoi(frequency)
		if f%10 != 0 {
			continue
		}
		if frequency != oldfreq {
			oldfreq = frequency
			parent = tsDetails.Append(nil)
			if err := tsDetails.SetValue(parent, 0, frequency); err != nil {
				log.Println(err.Error())
			}
			fmt.Sscanf(frequency, "%d", &iFreq)
			freqSort := fmt.Sprintf("%04d", iFreq)
			if err := tsDetails.SetValue(parent, 8, freqSort); err != nil {
				log.Println(err.Error())
			}
		}

		if night != "" && strings.Compare(night, day) != 0 {
			power = fmt.Sprintf("%v/%v", day, night)
		} else {
			power = day
		}

		fmt.Sscanf(power, "%f", &fPower)
		if fPower <= 0.01 {
			continue
		}
		powerSort := fmt.Sprintf("%010.5f", fPower)
		distanceSort := fmt.Sprintf("%015.5f", distance)
		bearingSort := fmt.Sprintf("%04.0f", bearing)

		col := []int{1, 2, 3, 4, 5, 6, 7, 9, 10, 11}
		var val []interface{}
		val = append(val, station, city, state, country, power, fmt.Sprintf("%5.0f", distance), fmt.Sprintf("% 3.0f", bearing),
			powerSort, distanceSort, bearingSort)
		if err := tsDetails.InsertWithValues(nil, parent, -1, col, val); err != nil {
			log.Println(err.Error())
		}
	}
}

func onMWListTreeButtonPressEvent(_ *gtk.TreeView, e *gdk.Event) {
	// Both of these work for a double click...which one's better?
	// gdk.EVENT_2BUTTON_PRESS
	// gdk.EVENT_DOUBLE_BUTTON_PRESS
	var l logging
	eb := gdk.EventButtonNewFromEvent(e)
	switch eb.Type() {
	case gdk.EVENT_DOUBLE_BUTTON_PRESS: // Double click
		l.open(-1)

	}

}

// func (x *gtk.TreeIterCompareFunc)  sortFunc(model *gtk.TreeModel, iter1, iter2 *gtk.TreeIter) int {
// 	// col = model.GetSortColumnId()
// 	val1, err := model.GetValue(iter1, 0)
// 	if err != nil {
// 		return 0
// 	}
// 	x,err  := val1.GoValue()
// 	y := x.(float64)
// 	_=y
// 	return 1
// }
