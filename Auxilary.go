package gorph

type Float64Point struct {
	X float64
	Y float64
}

func MaxInt(x, y int) int {
	if x > y {
		return x
	}
	return y
}
