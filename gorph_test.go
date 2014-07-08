package gorph

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"
)

func AssertEqualsUint32(t *testing.T, val1, val2 uint32, message ...string) {
	if val1 != val2 {
		t.Fail()
		if message != nil && testing.Verbose() {
			t.Log(message, val1, val2)
		} else if testing.Verbose() {
			t.Log(val1, val2)
		}
	}
}

func BenchmarkImageSet(b *testing.B) {
	n := b.N
	test := image.NewRGBA64(image.Rectangle{image.Point{0, 0}, image.Point{n, n}})
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			test.Set(i, j, color.RGBA64{0, 0, 0, 1<<16 - 1})
		}
	}
}

func BenchmarkImagePix(b *testing.B) {
	n := b.N
	test := image.NewRGBA64(image.Rectangle{image.Point{0, 0}, image.Point{n, n}})
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			temp := color.RGBA64{0, 0, 0, 1<<16 - 1}
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

func TestMergePixelsInLine(t *testing.T) {
	n := 16
	test := image.NewRGBA64(image.Rectangle{image.Point{0, 0}, image.Point{n, 1}})
	test.Set(0, 0, color.RGBA64{0, 0xffff, 0, 0xffff})
	for i := 1; i < n; i++ {
		test.Set(i, 0, color.RGBA64{0xffff, 0, 0, 0xffff})
	}
	testFile, err := os.Create("test.png")
	if err != nil {
		t.Fatal(err.Error())
	}
	err = png.Encode(testFile, test)
	if err != nil {
		t.Fatal(err.Error())
	}
	testFile.Close()
	testTwo := image.NewRGBA64(image.Rectangle{image.Point{0, 0}, image.Point{n, 1}})
	testTwo.Set(0, 0, color.RGBA64{0, 0, 0, 0xffff})
	for i := 1; i < n; i++ {
		testTwo.Set(i, 0, color.RGBA64{0, 0, 0xffff, 0xffff})
	}
	testFile, err = os.Create("testTwo.png")
	if err != nil {
		t.Fatal(err.Error())
	}
	err = png.Encode(testFile, testTwo)
	if err != nil {
		t.Fatal(err.Error())
	}
	testFile.Close()
	mergePixelsInLine(true, 0, true, true, 0.8, 3.5, 2.8, 5.2, test, testTwo)
	testFile, err = os.Create("testThree.png")
	if err != nil {
		t.Fatal(err.Error())
	}
	err = png.Encode(testFile, testTwo)
	if err != nil {
		t.Fatal(err.Error())
	}
	testFile.Close()
}

func TestMergePixelsInHorizontalLine(t *testing.T) {
	n := 4
	test := image.NewRGBA64(image.Rectangle{image.Point{0, 0}, image.Point{n, 1}})
	test.Set(0, 0, color.RGBA64{0, 0xffff, 0, 0xffff})
	for i := 1; i < n; i++ {
		test.Set(i, 0, color.RGBA64{0xffff, 0, 0, 0xffff})
	}
	testTwo := image.NewRGBA64(image.Rectangle{image.Point{0, 0}, image.Point{n, 1}})
	mergePixelsInLine(true, 0, false, false, 0.5, 1.0, 2.0, 3.0, test, testTwo)
	r, g, b, a := testTwo.At(2, 0).RGBA()
	AssertEqualsUint32(t, r, 0)
	AssertEqualsUint32(t, g, 0xffff)
	AssertEqualsUint32(t, b, 0)
	AssertEqualsUint32(t, a, 0xffff)
}

func TestStretchHorizontally(t *testing.T) {
	width := 4
	height := 4
	test := image.NewRGBA64(image.Rectangle{image.Point{0, 0}, image.Point{width, height}})
	test.Set(0, 0, color.RGBA64{0, 0xffff, 0, 0xffff})
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			test.Set(i, j, color.RGBA64{0xffff * uint16(i), 0xffff * uint16(height - j - 1), 0, 0xffff})	
		}
	}
	testFile, err := os.Create("PreStretch.png")
	if err != nil {
		t.Fatal(err.Error())
	}
	err = png.Encode(testFile, test)
	if err != nil {
		t.Fatal(err.Error())
	}
	testFile.Close()
	testTwo := image.NewRGBA64(image.Rectangle{image.Point{0, 0}, image.Point{width, height}})
	mGrid := NewMorphGrid()
	mGrid.AddPoints(0, 0, image.Point{0, 0}, image.Point{0, 0})
	mGrid.AddPoints(0, 2, image.Point{width, 0}, image.Point{width, 0})
	mGrid.AddPoints(2, 0, image.Point{0, height}, image.Point{0, height})
	mGrid.AddPoints(2, 2, image.Point{width, height}, image.Point{width, height})
	mGrid.AddPoints(1, 0, image.Point{0, 2}, image.Point{0, 2})
	mGrid.AddPoints(1, 1, image.Point{3, 2}, image.Point{3, 2})
	mGrid.AddPoints(1, 2, image.Point{width, 2}, image.Point{width, 2})
	mGrid.AddPoints(0, 1, image.Point{3, 0}, image.Point{3, 0})
	mGrid.AddPoints(2, 1, image.Point{3, height}, image.Point{3, height})
	start, end, nSplines, err := mGrid.allCubicCatmullRomSplines(true, 0.5, 5)
	AssertEqualsInt(t, nSplines, 3)
	err = stretchPixelsHorizontally(0, 2, start, end, test, testTwo)
	if err != nil {
		t.Fatal(err.Error())
	}
	testFile, err = os.Create("PostStretch.png")
	if err != nil {
		t.Fatal(err.Error())
	}
	err = png.Encode(testFile, test)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestCreateSameColor(t *testing.T) {
	colorOne := color.RGBA64{0, 0x1000, 0x2000, 0x1000}
	r, g, b, a := colorOne.RGBA()
	AssertEqualsUint32(t, r, 0)
	AssertEqualsUint32(t, g, 0x1000)
	AssertEqualsUint32(t, b, 0x2000)
	AssertEqualsUint32(t, a, 0x1000)
}

func TestWeightColorNone(t *testing.T) {
	colorOne := color.RGBA64{0, 0x1000, 0x2000, 0x1000}
	result := weightColor(colorOne, 0.0)
	r, g, b, a := result.RGBA()
	AssertEqualsUint32(t, r, 0)
	AssertEqualsUint32(t, g, 0)
	AssertEqualsUint32(t, b, 0)
	AssertEqualsUint32(t, a, 0)
}

func TestWeightColorOne(t *testing.T) {
	colorOne := color.RGBA64{0, 0x1000, 0x2000, 0x1000}
	result := weightColor(colorOne, 1.0)
	r, g, b, a := result.RGBA()
	AssertEqualsUint32(t, r, 0)
	AssertEqualsUint32(t, g, 0x1000)
	AssertEqualsUint32(t, b, 0x2000)
	AssertEqualsUint32(t, a, 0x1000)
}

func TestWeightColorOverflow(t *testing.T) {
	colorOne := color.RGBA64{0, 0x1000, 0x2000, 0x1000}
	result := weightColor(colorOne, 1<<4)
	r, g, b, a := result.RGBA()
	AssertEqualsUint32(t, r, 0)
	AssertEqualsUint32(t, g, 1<<16-1)
	AssertEqualsUint32(t, b, 1<<16-1)
	AssertEqualsUint32(t, a, 1<<16-1)
}

func TestAddColorsNone(t *testing.T) {
	colorOne := color.RGBA64{0, 0x1000, 0x2000, 0x1000}
	colorTwo := color.RGBA64{0, 0, 0, 0}
	result := addColors(colorOne, colorTwo)
	r, g, b, a := result.RGBA()
	AssertEqualsUint32(t, r, 0)
	AssertEqualsUint32(t, g, 0x1000)
	AssertEqualsUint32(t, b, 0x2000)
	AssertEqualsUint32(t, a, 0x1000)
}

func TestAddColors(t *testing.T) {
	colorOne := color.RGBA64{0, 0x1000, 0x2000, 0x1000}
	colorTwo := color.RGBA64{0x1254, 0x3333, 0x2222, 0x0909}
	result := addColors(colorOne, colorTwo)
	r, g, b, a := result.RGBA()
	AssertEqualsUint32(t, r, 0x1254)
	AssertEqualsUint32(t, g, 0x4333)
	AssertEqualsUint32(t, b, 0x4222)
	AssertEqualsUint32(t, a, 0x1909)
}

func TestAddColorsOverflow(t *testing.T) {
	colorOne := color.RGBA64{0, 0xfffe, 0x2000, 0x1000}
	colorTwo := color.RGBA64{0x1254, 0x0002, 0x2222, 0x0909}
	result := addColors(colorOne, colorTwo)
	r, g, b, a := result.RGBA()
	AssertEqualsUint32(t, r, 0x1254)
	AssertEqualsUint32(t, g, 0xffff) // Overflowed
	AssertEqualsUint32(t, b, 0x4222)
	AssertEqualsUint32(t, a, 0x1909)
}

func TestInterpolateColorsHalf(t *testing.T) {
	colorOne := color.RGBA64{0, 0x1000, 0x2000, 0x1000}
	colorTwo := color.RGBA64{0, 0, 0, 0}
	result := interpolateColors(colorOne, colorTwo, 0.5)
	r, g, b, a := result.RGBA()
	AssertEqualsUint32(t, r, 0)
	AssertEqualsUint32(t, g, 0x0800)
	AssertEqualsUint32(t, b, 0x1000)
	AssertEqualsUint32(t, a, 0x0800)
}

func TestInterpolateColorsOne(t *testing.T) {
	colorOne := color.RGBA64{0, 0x1000, 0x2000, 0x1000}
	colorTwo := color.RGBA64{0, 0, 0, 0}
	result := interpolateColors(colorOne, colorTwo, 1)
	r, g, b, a := result.RGBA()
	AssertEqualsUint32(t, r, 0)
	AssertEqualsUint32(t, g, 0x1000)
	AssertEqualsUint32(t, b, 0x2000)
	AssertEqualsUint32(t, a, 0x1000)
}

func TestInterpolateColorsEqual(t *testing.T) {
	colorOne := color.RGBA64{0, 0x1000, 0x2000, 0x1000}
	colorTwo := color.RGBA64{0, 0x1000, 0x2000, 0x1000}
	result := interpolateColors(colorOne, colorTwo, 0.33333333333)
	r, g, b, a := result.RGBA()
	AssertEqualsUint32(t, r, 0)
	AssertEqualsUint32(t, g, 0x1000)
	AssertEqualsUint32(t, b, 0x2000)
	AssertEqualsUint32(t, a, 0x1000)
}
