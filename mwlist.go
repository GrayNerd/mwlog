package main

import (
	"fmt"
	"strconv"
	"strings"

	// "github.com/gotk3/gotk3/gtk"
	"log"
	"mwlog/db"
	"mwlog/ui"

	"github.com/gotk3/gotk3/gtk"
)

type mwListTab struct {
}

func (lt *mwListTab) showMWListTab() {
	ls := ui.GetListStore("mwlist_store")
	ls.Clear()

	var sortFreq gtk.TreeIterCompareFunc = func(model *gtk.TreeModel, iter1, iter2 *gtk.TreeIter) int {
		cv1, _ := model.GetValue(iter1, 0)
		cv2, _ := model.GetValue(iter2, 0)
		gv1,_  := cv1.GoValue()
		gv2, _ := cv2.GoValue()
		v1,_ := strconv.Atoi(gv1.(string))
		v2, _ := strconv.Atoi(gv2.(string))
		if v1 < v2 {
			return -1
		}
		if v1 > v2 {
			return 1
		}
		return 0
	}

	var sortDistance gtk.TreeIterCompareFunc = func(model *gtk.TreeModel, iter1, iter2 *gtk.TreeIter) int {
		cv1, _ := model.GetValue(iter1, 6)
		cv2, _ := model.GetValue(iter2, 6)
		gv1,_  := cv1.GoValue()
		gv2, _ := cv2.GoValue()
		s1 ,_:= gv1.(string)
		s2 ,_:= gv2.(string)
		v1,_ := strconv.Atoi(strings.Trim(s1, " "))
		v2, _ := strconv.Atoi(strings.Trim(s2, " "))
		if v1 < v2 {
			return -1
		}
		if v1 > v2 {
			return 1
		}
		return 0
	}

	ls.SetSortFunc(0, sortFreq)
	ls.SetSortFunc(6, sortDistance)

	rows := db.GetAllMWList()
	defer rows.Close()
	var station, frequency, city, state, country, day, night, power string
	var distance, bearing float64
	for rows.Next() {
		rows.Scan(&frequency, &station, &city, &state, &country, &day, &night, &distance, &bearing)
		if night != "" {
			power = fmt.Sprintf("%v/%v", day, night)
		} else {
			power = day
		}
		col := []int{0, 1, 2, 3, 4, 5, 6, 7}
		var val []interface{}
		val = append(val, frequency, station, city, state, country, power, fmt.Sprintf("%5.0f", distance), fmt.Sprintf("% 3.0f", bearing))
		if err := ls.InsertWithValues(nil, -1, col, val); err != nil {
			log.Println(err.Error())
		}
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
