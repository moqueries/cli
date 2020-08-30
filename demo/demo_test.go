package demo_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/myshkin5/moqueries/demo"
	"github.com/myshkin5/moqueries/pkg/hash"
)

var _ = Describe("Demo", func() {
	It("only writes my favorite numbers", func() {
		isFavMock := newMockIsFavorite()
		writerMock := newMockWriter()

		isFavMock.onCall(1).ret(false)
		isFavMock.onCall(2).ret(false)
		isFavMock.onCall(3).ret(true)

		writerMock.onCall().Write([]byte("3"))

		d := demo.FavWriter{
			IsFav: isFavMock.mock(),
			W:     writerMock.mock(),
		}

		err := d.WriteFavorites([]int{1, 2, 3})

		Expect(err).NotTo(HaveOccurred())
		Expect(writerMock.params_Write).To(Receive(&mockWriter_Write_params{p: hash.DeepHash([]byte("3"))}))
	})

	It("returns an error when writing", func() {
		isFavMock := newMockIsFavorite()
		writerMock := newMockWriter()

		isFavMock.onCall(3).ret(true)

		writerMock.onCall().Write([]byte("3")).ret(0, errors.New("couldn't write"))

		d := demo.FavWriter{
			IsFav: isFavMock.mock(),
			W:     writerMock.mock(),
		}

		err := d.WriteFavorites([]int{1, 2, 3})

		Expect(err).To(MatchError("that pesky writer says that it 'couldn't write'"))
	})
})
