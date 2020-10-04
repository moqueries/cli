package demo_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/myshkin5/moqueries/demo"
	"github.com/myshkin5/moqueries/pkg/config"
	"github.com/myshkin5/moqueries/pkg/hash"
)

func TestOnlyWriteFavoriteNumbers(t *testing.T) {
	isFavMock := newMockIsFavorite(t, nil)
	writerMock := newMockWriter(t, nil)

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
	if len(writerMock.params_Write) != 1 {
		t.Errorf("expected Write to be called once, called %d",
			len(writerMock.params_Write))
		return
	}
	params := <-writerMock.params_Write
	if params.p != hash.DeepHash([]byte("3")) {
		t.Errorf("unexpected parameters in call to Write")
	}
}

func TestWriteError(t *testing.T) {
	isFavMock := newMockIsFavorite(
		t, &config.MockConfig{Expectation: config.Nice})
	writerMock := newMockWriter(t, nil)

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
}

func TestChangedMyMindILikeIt(t *testing.T) {
	isFavMock := newMockIsFavorite(t, nil)
	writerMock := newMockWriter(t, nil)

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
	if len(writerMock.params_Write) != 1 {
		t.Errorf("expected Write to be called once, called %d",
			len(writerMock.params_Write))
		return
	}
	params := <-writerMock.params_Write
	if params.p != hash.DeepHash([]byte("7")) {
		t.Errorf("unexpected parameters in call to Write")
	}
}

func TestChangedMyMindIHateIt(t *testing.T) {
	isFavMock := newMockIsFavorite(t, nil)
	writerMock := newMockWriter(t, nil)

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
	if len(writerMock.params_Write) != 2 {
		t.Errorf("expected Write to be called once, called %d",
			len(writerMock.params_Write))
		return
	}
	params := <-writerMock.params_Write
	if params.p != hash.DeepHash([]byte("7")) {
		t.Errorf("unexpected parameters in call to Write")
	}
	params = <-writerMock.params_Write
	if params.p != hash.DeepHash([]byte("7")) {
		t.Errorf("unexpected parameters in call to Write")
	}
}
