package demo_test

import (
	"errors"
	"math"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/myshkin5/moqueries/demo"
	"github.com/myshkin5/moqueries/moq"
)

func TestOnlyWriteFavoriteNumbers(t *testing.T) {
	scene := moq.NewScene(t)
	isFavMoq := newMoqIsFavorite(scene, nil)
	writerMoq := newMoqWriter(scene, nil)

	isFavMoq.onCall(1).returnResults(false)
	isFavMoq.onCall(2).returnResults(false)
	isFavMoq.onCall(3).returnResults(true)

	writerMoq.onCall().Write([]byte("3")).
		returnResults(1, nil)

	d := demo.FavWriter{
		IsFav: isFavMoq.mock(),
		W:     writerMoq.mock(),
	}

	err := d.WriteFavorites([]int{1, 2, 3})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	writerMoq.AssertExpectationsMet()
}

func TestWriteError(t *testing.T) {
	scene := moq.NewScene(t)
	isFavMoq := newMoqIsFavorite(
		scene, &moq.Config{Expectation: moq.Nice})
	writerMoq := newMoqWriter(scene, nil)

	isFavMoq.onCall(3).returnResults(true)

	writerMoq.onCall().Write([]byte("3")).
		returnResults(0, errors.New("couldn't write"))

	d := demo.FavWriter{
		IsFav: isFavMoq.mock(),
		W:     writerMoq.mock(),
	}

	err := d.WriteFavorites([]int{1, 2, 3})
	if err == nil {
		t.Errorf("expected error")
		return
	}
	if !strings.Contains(err.Error(),
		"that pesky writer says that it 'couldn't write'") {
		t.Errorf("unexpected message in error: %v", err)
	}

	scene.AssertExpectationsMet()
}

func TestChangedMyMindILikeIt(t *testing.T) {
	scene := moq.NewScene(t)
	isFavMoq := newMoqIsFavorite(scene, nil)
	writerMoq := newMoqWriter(scene, nil)

	isFavMoq.onCall(7).
		returnResults(false).repeat(moq.Times(5)).
		returnResults(true)
	isFavMoq.onCall(3).
		returnResults(false)

	writerMoq.onCall().Write([]byte("7")).
		returnResults(1, nil)

	d := demo.FavWriter{
		IsFav: isFavMoq.mock(),
		W:     writerMoq.mock(),
	}

	err := d.WriteFavorites([]int{7, 7, 7, 7, 7, 7, 3})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	scene.AssertExpectationsMet()
}

func TestChangedMyMindImNotSure(t *testing.T) {
	scene := moq.NewScene(t)
	config := moq.Config{Sequence: moq.SeqDefaultOn}
	isFavMoq := newMoqIsFavorite(scene, &config)
	writerMoq := newMoqWriter(scene, &config)

	isFavMoq.onCall(7).
		returnResults(false).repeat(moq.Times(3))
	isFavMoq.onCall(3).
		returnResults(false)

	isFavMoq.onCall(7).
		returnResults(true)
	writerMoq.onCall().Write([]byte("7")).
		returnResults(1, nil)

	isFavMoq.onCall(7).
		returnResults(false).repeat(moq.Times(2))

	d := demo.FavWriter{
		IsFav: isFavMoq.mock(),
		W:     writerMoq.mock(),
	}

	err := d.WriteFavorites([]int{7, 7, 7, 3, 7, 7, 7})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	scene.AssertExpectationsMet()
}

func TestChangedMyMindIHateIt(t *testing.T) {
	scene := moq.NewScene(t)
	isFavMoq := newMoqIsFavorite(scene, nil)
	writerMoq := newMoqWriter(scene, nil)

	isFavMoq.onCall(7).
		returnResults(true).repeat(moq.Times(2)).
		returnResults(false).repeat(moq.AnyTimes())

	err := errors.New("I no longer like 7")
	writerMoq.onCall().Write([]byte("7")).
		returnResults(0, nil).repeat(moq.Times(2)).
		returnResults(0, err).repeat(moq.AnyTimes())

	d := demo.FavWriter{
		IsFav: isFavMoq.mock(),
		W:     writerMoq.mock(),
	}

	err = d.WriteFavorites([]int{7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	scene.AssertExpectationsMet()
}

func TestNoGadgets(t *testing.T) {
	scene := moq.NewScene(t)
	writerMoq := newMoqWriter(scene, nil)
	storeMoq := newMoqStore(scene, nil)

	storeMoq.onCall().AllWidgetsIds().
		returnResults([]int{42, 43}, nil)

	storeMoq.onCall().GadgetsByWidgetId(0).any().widgetId().
		returnResults(nil, nil).repeat(moq.Times(2))

	d := demo.FavWriter{
		W:     writerMoq.mock(),
		Store: storeMoq.mock(),
	}

	expected := map[int]demo.Widget{
		42: {Id: 42, GadgetsByColor: map[string]demo.Gadget{}},
		43: {Id: 43, GadgetsByColor: map[string]demo.Gadget{}},
	}

	widgets, err := d.Load(math.MaxUint32)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(widgets, expected) {
		t.Errorf("unexpected difference in loaded widgets: %#v", widgets)
	}

	scene.AssertExpectationsMet()
}

func TestLightGadgets(t *testing.T) {
	scene := moq.NewScene(t)
	writerMoq := newMoqWriter(scene, nil)
	storeMoq := newMoqStore(scene, nil)

	storeMoq.onCall().AllWidgetsIds().
		returnResults([]int{42, 43}, nil)

	g1 := demo.Gadget{
		Id:       4201,
		WidgetId: 42,
		Color:    "red",
		Weight:   200,
	}
	g2 := demo.Gadget{
		Id:       4202,
		WidgetId: 42,
		Color:    "blue",
		Weight:   201,
	}
	storeMoq.onCall().LightGadgetsByWidgetId(42, 0).any().maxWeight().
		returnResults([]demo.Gadget{g1, g2}, nil)
	g3 := demo.Gadget{
		Id:       4301,
		WidgetId: 43,
		Color:    "grey",
		Weight:   100,
	}
	g4 := demo.Gadget{
		Id:       4302,
		WidgetId: 43,
		Color:    "heliotrope",
		Weight:   101,
	}
	storeMoq.onCall().LightGadgetsByWidgetId(0, 0).
		any().widgetId().any().maxWeight().
		returnResults([]demo.Gadget{g3, g4}, nil)

	d := demo.FavWriter{
		W:     writerMoq.mock(),
		Store: storeMoq.mock(),
	}

	expected := map[int]demo.Widget{
		42: {
			Id: 42,
			GadgetsByColor: map[string]demo.Gadget{
				"red":  g1,
				"blue": g2,
			},
		},
		43: {
			Id: 43,
			GadgetsByColor: map[string]demo.Gadget{
				"grey":       g3,
				"heliotrope": g4,
			},
		},
	}

	widgets, err := d.Load(8382)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(widgets, expected) {
		t.Errorf("unexpected difference in loaded widgets: %#v", widgets)
	}

	scene.AssertExpectationsMet()
}

func TestOnlyWriteFavoriteNumbersSeqBySeq(t *testing.T) {
	scene := moq.NewScene(t)
	isFavMoq := newMoqIsFavorite(scene, nil)
	writerMoq := newMoqWriter(scene, nil)

	isFavMoq.onCall(1).seq().returnResults(false)
	isFavMoq.onCall(2).seq().returnResults(false)
	isFavMoq.onCall(3).seq().returnResults(true)

	writerMoq.onCall().Write([]byte("3")).
		returnResults(1, nil)

	d := demo.FavWriter{
		IsFav: isFavMoq.mock(),
		W:     writerMoq.mock(),
	}

	err := d.WriteFavorites([]int{1, 2, 3})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	writerMoq.AssertExpectationsMet()
}

func TestOnlyWriteFavoriteNumbersSeqByNoSeq(t *testing.T) {
	scene := moq.NewScene(t)
	config := moq.Config{Sequence: moq.SeqDefaultOn}
	isFavMoq := newMoqIsFavorite(scene, &config)
	writerMoq := newMoqWriter(scene, &config)

	isFavMoq.onCall(1).returnResults(false)
	isFavMoq.onCall(2).returnResults(false)
	isFavMoq.onCall(3).returnResults(true)

	writerMoq.onCall().Write([]byte("3")).noSeq().
		returnResults(1, nil)

	d := demo.FavWriter{
		IsFav: isFavMoq.mock(),
		W:     writerMoq.mock(),
	}

	err := d.WriteFavorites([]int{1, 2, 3})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	writerMoq.AssertExpectationsMet()
}

func TestOnlyWriteFavoriteNumbersWithDoFn(t *testing.T) {
	scene := moq.NewScene(t)
	isFavMoq := newMoqIsFavorite(scene, nil)
	writerMoq := newMoqWriter(scene, nil)

	sum := 0
	sumFn := func(n int) {
		sum += n
	}

	isFavMoq.onCall(1).returnResults(false).andDo(sumFn)
	isFavMoq.onCall(2).returnResults(false).andDo(sumFn)
	isFavMoq.onCall(3).returnResults(true).andDo(sumFn)

	writerMoq.onCall().Write([]byte("3")).
		returnResults(1, nil)

	d := demo.FavWriter{
		IsFav: isFavMoq.mock(),
		W:     writerMoq.mock(),
	}

	err := d.WriteFavorites([]int{1, 2, 3})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if sum != 6 {
		t.Errorf("unexpected sum: %d", sum)
	}

	writerMoq.AssertExpectationsMet()
}

func TestOnlyWriteFavoriteNumbersWithDoReturn(t *testing.T) {
	scene := moq.NewScene(t)
	isFavMoq := newMoqIsFavorite(scene, nil)
	writerMoq := newMoqWriter(scene, nil)

	isFavFn := func(n int) bool {
		return n%2 == 0
	}

	isFavMoq.onCall(0).any().n().
		doReturnResults(isFavFn).repeat(moq.AnyTimes())

	bytesWritten := 0
	var capturedFavs []int
	writeFn := func(p []byte) (n int, err error) {
		bytesWritten += len(p)
		fav, err := strconv.Atoi(string(p))
		if err != nil {
			t.Errorf("error parsing fav number: %v", p)
		}
		capturedFavs = append(capturedFavs, fav)
		return 0, nil
	}

	writerMoq.onCall().Write(nil).any().p().
		doReturnResults(writeFn).repeat(moq.AnyTimes())

	d := demo.FavWriter{
		IsFav: isFavMoq.mock(),
		W:     writerMoq.mock(),
	}

	err := d.WriteFavorites([]int{1, 2, 3, 4, 5, 6})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if bytesWritten != 3 {
		t.Errorf("unexpected bytes written: %d", bytesWritten)
	}

	if !reflect.DeepEqual(capturedFavs, []int{2, 4, 6}) {
		t.Errorf("unexpected captured favorites: %v", capturedFavs)
	}

	writerMoq.AssertExpectationsMet()
}
