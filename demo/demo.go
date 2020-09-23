package demo

import (
	"fmt"
	"io"
	"strconv"
)

//go:generate moqueries --destination moq_writer_test.go --import io Writer

//go:generate moqueries --destination moq_isfavorite_test.go IsFavorite

type IsFavorite func(n int) bool

type FavWriter struct {
	IsFav IsFavorite
	W     io.Writer
}

func (d *FavWriter) WriteFavorites(nums []int) error {
	for _, num := range nums {
		if d.IsFav(num) {
			_, err := d.W.Write([]byte(strconv.Itoa(num)))
			if err != nil {
				return fmt.Errorf("that pesky writer says that it '%s'",
					err.Error())
			}
		}
	}
	return nil
}
