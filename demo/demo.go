// Package demo provides code examples used in the ../README.md
package demo

import (
	"fmt"
	"io"
	"math"
	"strconv"
)

//go:generate moqueries --import io Writer

//go:generate moqueries IsFavorite

// IsFavorite is the function type for a function that determines if a number
// is a favorite
type IsFavorite func(n int) bool

// FavWriter is used to filter and write only favorite number
type FavWriter struct {
	IsFav IsFavorite
	W     io.Writer
	Store Store
}

// WriteFavorites iterates through a given list of numbers and writes any
// favorite number to the given Writer
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

// Widget holds the state of a widget
type Widget struct {
	Id             int
	GadgetsByColor map[string]Gadget
}

// Gadget holds the state of a gadget
type Gadget struct {
	Id       int
	WidgetId int
	Color    string
	Weight   uint32
}

//go:generate moqueries Store

// Store is the abstraction for Widget and Gadget types
type Store interface {
	AllWidgetsIds() ([]int, error)
	GadgetsByWidgetId(widgetId int) ([]Gadget, error)
	LightGadgetsByWidgetId(widgetId int, maxWeight uint32) ([]Gadget, error)
}

// Load retrieves multiple Widget objects
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
