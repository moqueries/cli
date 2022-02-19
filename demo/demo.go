package demo

import (
	"fmt"
	"io"
	"math"
	"strconv"
)

//go:generate moqueries --import io Writer

//go:generate moqueries IsFavorite

type IsFavorite func(n int) bool

type FavWriter struct {
	IsFav IsFavorite
	W     io.Writer
	Store Store
}

func (d *FavWriter) WriteFavorites(nums []int) error {
	for _, num := range nums {
		if d.IsFav(num) {
			_, err := d.W.Write([]byte(strconv.Itoa(num)))
			if err != nil {
				return fmt.Errorf("that pesky writer says that it '%w'",
					err)
			}
		}
	}
	return nil
}

type Widget struct {
	Id             int
	GadgetsByColor map[string]Gadget
}

type Gadget struct {
	Id       int
	WidgetId int
	Color    string
	Weight   uint32
}

//go:generate moqueries Store

type Store interface {
	AllWidgetsIds() ([]int, error)
	GadgetsByWidgetId(widgetId int) ([]Gadget, error)
	LightGadgetsByWidgetId(widgetId int, maxWeight uint32) ([]Gadget, error)
}

func (d *FavWriter) Load(maxWeight uint32) (map[int]Widget, error) {
	wIds, err := d.Store.AllWidgetsIds()
	if err != nil {
		return nil, fmt.Errorf("failed to load widget ids: %w", err)
	}

	widgets := make(map[int]Widget)
	for _, wId := range wIds {
		w, ok := widgets[wId]
		if !ok {
			w = Widget{Id: wId, GadgetsByColor: make(map[string]Gadget)}
			widgets[wId] = w
		}

		var gadgets []Gadget
		if maxWeight == math.MaxUint32 {
			gadgets, err = d.Store.GadgetsByWidgetId(wId)
		} else {
			gadgets, err = d.Store.LightGadgetsByWidgetId(wId, maxWeight)
		}
		if err != nil {
			return nil, fmt.Errorf("failed to load gadgets for widget %d ids: %w", wId, err)
		}

		for _, g := range gadgets {
			w.GadgetsByColor[g.Color] = g
		}
	}

	return widgets, nil
}
