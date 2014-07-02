package gorph

import (
	"image"
	"image/color"
	"testing"
)

func BenchmarkImageSet(b *testing.B) {
	n := b.N
	test := image.NewNRGBA64(image.Rectangle{image.Point{0, 0}, image.Point{n, n}})
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			test.Set(i, j, color.NRGBA64{0, 0, 0, 1<<16 - 1})
		}
	}
}

func BenchmarkImagePix(b *testing.B) {
	n := b.N
	test := image.NewNRGBA64(image.Rectangle{image.Point{0, 0}, image.Point{n, n}})
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			temp := color.NRGBA64{0, 0, 0, 1<<16 - 1}
			offset := test.PixOffset(i, j)
			test.Pix[offset] = uint8(temp.R >> 8)
			test.Pix[offset+1] = uint8(temp.R)
			test.Pix[offset+2] = uint8(temp.G >> 8)
			test.Pix[offset+3] = uint8(temp.G)
			test.Pix[offset+4] = uint8(temp.B >> 8)
			test.Pix[offset+5] = uint8(temp.B)
			test.Pix[offset+6] = uint8(temp.A >> 8)
			test.Pix[offset+7] = uint8(temp.A)
		}
	}
}
