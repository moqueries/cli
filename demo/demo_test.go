package demo_test

import (
	"errors"
	"math"
	"reflect"
	"strings"
	"testing"

	"github.com/myshkin5/moqueries/demo"
	"github.com/myshkin5/moqueries/pkg/moq"
)

func TestOnlyWriteFavoriteNumbers(t *testing.T) {
	scene := moq.NewScene(t)
	isFavMock := newMockIsFavorite(scene, nil)
	writerMock := newMockWriter(scene, nil)

	isFavMock.onCall(1).returnResults(false)
	isFavMock.onCall(2).returnResults(false)
	isFavMock.onCall(3).returnResults(true)

	writerMock.onCall().Write([]byte("3")).
		returnResults(1, nil)

	d := demo.FavWriter{
		IsFav: isFavMock.mock(),
		W:     writerMock.mock(),
	}

	err := d.WriteFavorites([]int{1, 2, 3})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	writerMock.AssertExpectationsMet()
}

func TestWriteError(t *testing.T) {
	scene := moq.NewScene(t)
	isFavMock := newMockIsFavorite(
		scene, &moq.MockConfig{Expectation: moq.Nice})
	writerMock := newMockWriter(scene, nil)

	isFavMock.onCall(3).returnResults(true)

	writerMock.onCall().Write([]byte("3")).
		returnResults(0, errors.New("couldn't write"))

	d := demo.FavWriter{
		IsFav: isFavMock.mock(),
		W:     writerMock.mock(),
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
	isFavMock := newMockIsFavorite(scene, nil)
	writerMock := newMockWriter(scene, nil)

	isFavMock.onCall(7).
		returnResults(false).times(5).
		returnResults(true)
	isFavMock.onCall(3).
		returnResults(false)

	writerMock.onCall().Write([]byte("7")).
		returnResults(1, nil)

	d := demo.FavWriter{
		IsFav: isFavMock.mock(),
		W:     writerMock.mock(),
	}

	err := d.WriteFavorites([]int{7, 7, 7, 7, 7, 7, 3})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	scene.AssertExpectationsMet()
}

func TestChangedMyMindImNotSure(t *testing.T) {
	scene := moq.NewScene(t)
	config := moq.MockConfig{Sequence: moq.SeqDefaultOn}
	isFavMock := newMockIsFavorite(scene, &config)
	writerMock := newMockWriter(scene, &config)

	isFavMock.onCall(7).
		returnResults(false).times(3)
	isFavMock.onCall(3).
		returnResults(false)

	isFavMock.onCall(7).
		returnResults(true)
	writerMock.onCall().Write([]byte("7")).
		returnResults(1, nil)

	isFavMock.onCall(7).
		returnResults(false).times(2)

	d := demo.FavWriter{
		IsFav: isFavMock.mock(),
		W:     writerMock.mock(),
	}

	err := d.WriteFavorites([]int{7, 7, 7, 3, 7, 7, 7})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	scene.AssertExpectationsMet()
}

func TestChangedMyMindIHateIt(t *testing.T) {
	scene := moq.NewScene(t)
	isFavMock := newMockIsFavorite(scene, nil)
	writerMock := newMockWriter(scene, nil)

	isFavMock.onCall(7).
		returnResults(true).times(2).
		returnResults(false).anyTimes()

	writerMock.onCall().Write([]byte("7")).
		returnResults(0, nil).times(2).
		returnResults(0, errors.New("I no longer like 7")).anyTimes()

	d := demo.FavWriter{
		IsFav: isFavMock.mock(),
		W:     writerMock.mock(),
	}

	err := d.WriteFavorites([]int{7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	scene.AssertExpectationsMet()
}

func TestNoGadgets(t *testing.T) {
	scene := moq.NewScene(t)
	writerMock := newMockWriter(scene, nil)
	storeMock := newMockStore(scene, nil)

	storeMock.onCall().AllWidgetsIds().
		returnResults([]int{42, 43}, nil)

	storeMock.onCall().GadgetsByWidgetId(0).anyWidgetId().
		returnResults(nil, nil).times(2)

	d := demo.FavWriter{
		W:     writerMock.mock(),
		Store: storeMock.mock(),
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
	writerMock := newMockWriter(scene, nil)
	storeMock := newMockStore(scene, nil)

	storeMock.onCall().AllWidgetsIds().
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
	storeMock.onCall().LightGadgetsByWidgetId(42, 0).anyMaxWeight().
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
	storeMock.onCall().LightGadgetsByWidgetId(0, 0).anyWidgetId().anyMaxWeight().
		returnResults([]demo.Gadget{g3, g4}, nil)

	d := demo.FavWriter{
		W:     writerMock.mock(),
		Store: storeMock.mock(),
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
	isFavMock := newMockIsFavorite(scene, nil)
	writerMock := newMockWriter(scene, nil)

	isFavMock.onCall(1).seq().returnResults(false)
	isFavMock.onCall(2).seq().returnResults(false)
	isFavMock.onCall(3).seq().returnResults(true)

	writerMock.onCall().Write([]byte("3")).
		returnResults(1, nil)

	d := demo.FavWriter{
		IsFav: isFavMock.mock(),
		W:     writerMock.mock(),
	}

	err := d.WriteFavorites([]int{1, 2, 3})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	writerMock.AssertExpectationsMet()
}

func TestOnlyWriteFavoriteNumbersSeqByNoSeq(t *testing.T) {
	scene := moq.NewScene(t)
	config := moq.MockConfig{Sequence: moq.SeqDefaultOn}
	isFavMock := newMockIsFavorite(scene, &config)
	writerMock := newMockWriter(scene, &config)

	isFavMock.onCall(1).returnResults(false)
	isFavMock.onCall(2).returnResults(false)
	isFavMock.onCall(3).returnResults(true)

	writerMock.onCall().Write([]byte("3")).noSeq().
		returnResults(1, nil)

	d := demo.FavWriter{
		IsFav: isFavMock.mock(),
		W:     writerMock.mock(),
	}

	err := d.WriteFavorites([]int{1, 2, 3})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	writerMock.AssertExpectationsMet()
}
