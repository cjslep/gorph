package gorph

// Float64Point represents a double precision point. It is capable of representing
// fractional pixels in the same intention as an image.Point. A floating point of
// (1.25, 2.5) represents the location a quarter of the way horizontally into and
// halfway down the pixel given by image.Point(1, 2).
type Float64Point struct {
	X float64
	Y float64
}

// MaxInt is a convenience function that returns the larger of two integers.
func MaxInt(x, y int) int {
	if x > y {
		return x
	}
	return y
}
