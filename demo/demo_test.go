package demo_test

import (
	"errors"
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
