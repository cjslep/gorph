package gorph

import (
	"image"
	"math"
	"testing"
)

func LogVerbose(t *testing.T, args ...interface{}) {
	if testing.Verbose() {
		t.Log(args...)
	}
}

func LogfVerbose(t *testing.T, format string, args ...interface{}) {
	if testing.Verbose() {
		t.Logf(format, args...)
	}
}

func AssertEqualsInt(t *testing.T, val1, val2 int, message ...string) {
	if val1 != val2 {
		t.Fail()
		if message != nil && testing.Verbose() {
			t.Log(message, val1, val2)
		} else if testing.Verbose() {
			t.Log(val1, val2)
		}
	}
}

func AssertEqualsImagePoint(t *testing.T, p1, p2 image.Point, message ...string) {
	if !p1.Eq(p2) {
		t.Fail()
		if message != nil && testing.Verbose() {
			t.Log(message, p1, p2)
		} else if testing.Verbose() {
			t.Log(p1, p2)
		}
	}
}

func AssertEqualsFloat64Point(t *testing.T, p1, p2 Float64Point, message ...string) {
	if p1.X != p2.X || p1.Y != p2.Y {
		t.Fail()
		if message != nil && testing.Verbose() {
			t.Log(message, p1, p2)
		} else if testing.Verbose() {
			t.Log(p1, p2)
		}
	}
}

func AssertEqualsFloat64PointTolerance(t *testing.T, p1, p2 Float64Point, tolerance float64, message ...string) {
	if math.Abs(p1.X-p2.X) >= tolerance || math.Abs(p1.Y-p2.Y) >= tolerance {
		t.Fail()
		if message != nil && testing.Verbose() {
			t.Log(message, p1, p2)
		} else if testing.Verbose() {
			t.Log(p1, p2)
		}
	}
}

func TestMorphGridAddPoints(t *testing.T) {
	m := NewMorphGrid()
	m.AddPoints(3, 2, image.Point{1, 2}, image.Point{3, 4})
	p1, p2, err := m.Points(3, 2)
	if err != nil {
		LogVerbose(t, err.Error())
		t.Fail()
	}
	AssertEqualsImagePoint(t, p1, image.Point{1, 2})
	AssertEqualsImagePoint(t, p2, image.Point{3, 4})
}

func TestMorphGridLineCount(t *testing.T) {
	m := NewMorphGrid()
	m.AddPoints(3, 2, image.Point{1, 2}, image.Point{3, 4})
	m.AddPoints(3, 4, image.Point{2, 2}, image.Point{4, 6})
	length := m.HorizontalGridlineCount()
	AssertEqualsInt(t, length, 1, "HorizontalGridlineCount failed")
	length = m.VerticalGridlineCount()
	AssertEqualsInt(t, length, 2, "VerticalGridlineCount failed")
}

func TestMorphGridPoints(t *testing.T) {
	m := NewMorphGrid()
	m.AddPoints(3, 2, image.Point{1, 2}, image.Point{3, 4})
	m.AddPoints(3, 4, image.Point{2, 2}, image.Point{4, 6})
	source, dest, err := m.Points(3, 2)
	if err != nil {
		LogVerbose(t, err.Error())
		t.Fail()
	}
	AssertEqualsImagePoint(t, source, image.Point{1, 2})
	AssertEqualsImagePoint(t, dest, image.Point{3, 4})
}

func TestMorphGridHorizontalLine(t *testing.T) {
	m := NewMorphGrid()
	m.AddPoints(3, 2, image.Point{1, 2}, image.Point{3, 4})
	m.AddPoints(3, 4, image.Point{2, 2}, image.Point{4, 6})
	source, dest := m.HorizontalLine(3)
	AssertEqualsInt(t, len(source), 2, "Source length incorrect")
	AssertEqualsInt(t, len(dest), 2, "Destination length incorrect")
	AssertEqualsImagePoint(t, source[0], image.Point{1, 2})
	AssertEqualsImagePoint(t, source[1], image.Point{2, 2})
	AssertEqualsImagePoint(t, dest[0], image.Point{3, 4})
	AssertEqualsImagePoint(t, dest[1], image.Point{4, 6})
}

func TestMorphGridVerticalLine(t *testing.T) {
	m := NewMorphGrid()
	m.AddPoints(3, 2, image.Point{1, 2}, image.Point{3, 4})
	m.AddPoints(3, 4, image.Point{2, 2}, image.Point{4, 6})
	source, dest := m.VerticalLine(2)
	AssertEqualsInt(t, len(source), 1, "Source length incorrect")
	AssertEqualsInt(t, len(dest), 1, "Destination length incorrect")
	AssertEqualsImagePoint(t, source[0], image.Point{1, 2})
	AssertEqualsImagePoint(t, dest[0], image.Point{3, 4})
}

func TestMorphGridRemovePoints(t *testing.T) {
	m := NewMorphGrid()
	m.AddPoints(3, 2, image.Point{1, 2}, image.Point{3, 4})
	m.AddPoints(3, 4, image.Point{2, 2}, image.Point{4, 6})
	m.AddPoints(4, 4, image.Point{5, 7}, image.Point{6, 10})
	source, dest := m.HorizontalLine(3)
	AssertEqualsInt(t, len(source), 2, "Before removal source length incorrect")
	AssertEqualsInt(t, len(dest), 2, "Before removal dest length incorrect")
	err := m.RemovePoints(3, 2)
	if err != nil {
		LogVerbose(t, err.Error())
		t.Fail()
	}
	source, dest = m.HorizontalLine(3)
	AssertEqualsInt(t, len(source), 1, "After removal source length incorrect")
	AssertEqualsInt(t, len(dest), 1, "After removal dest length incorrect")
}

func TestMorphGridRemovePointsVertical(t *testing.T) {
	m := NewMorphGrid()
	m.AddPoints(3, 2, image.Point{1, 2}, image.Point{3, 4})
	m.AddPoints(3, 4, image.Point{2, 2}, image.Point{4, 6})
	m.AddPoints(4, 4, image.Point{5, 7}, image.Point{6, 10})
	source, dest := m.VerticalLine(4)
	AssertEqualsInt(t, len(source), 2, "Before removal source length incorrect")
	AssertEqualsInt(t, len(dest), 2, "Before removal dest length incorrect")
	err := m.RemovePoints(3, 4)
	if err != nil {
		LogVerbose(t, err.Error())
		t.Fail()
	}
	source, dest = m.VerticalLine(4)
	AssertEqualsInt(t, len(source), 1, "After removal source length incorrect")
	AssertEqualsInt(t, len(dest), 1, "After removal dest length incorrect")
}

func TestMorphGridSquare(t *testing.T) {
	width := 4
	height := 4
	mGrid := NewMorphGrid()
	mGrid.AddPoints(0, 0, image.Point{0, 0}, image.Point{0, 0})
	mGrid.AddPoints(0, 2, image.Point{width, 0}, image.Point{width, 0})
	mGrid.AddPoints(2, 0, image.Point{0, height}, image.Point{0, height})
	mGrid.AddPoints(2, 2, image.Point{width, height}, image.Point{width, height})
	mGrid.AddPoints(1, 0, image.Point{0, 2}, image.Point{0, 2})
	mGrid.AddPoints(1, 1, image.Point{2, 2}, image.Point{3, 2})
	mGrid.AddPoints(1, 2, image.Point{width, 2}, image.Point{width, 2})
	mGrid.AddPoints(0, 1, image.Point{2, 0}, image.Point{3, 0})
	mGrid.AddPoints(2, 1, image.Point{2, height}, image.Point{3, height})
	source, dest := mGrid.VerticalLine(1)
	AssertEqualsInt(t, len(source), 3)
	AssertEqualsInt(t, source[0].X, 2, "source[0].X failed")
	AssertEqualsInt(t, source[0].Y, 0, "source[0].Y failed")
	AssertEqualsInt(t, source[1].X, 2, "source[1].X failed")
	AssertEqualsInt(t, source[1].Y, 2, "source[1].Y failed")
	AssertEqualsInt(t, source[2].X, 2, "source[2].X failed")
	AssertEqualsInt(t, source[2].Y, height, "source[2].Y failed")
	AssertEqualsInt(t, len(dest), 3)
}

func TestMorphGridInterpolatedGrid(t *testing.T) {
	m := NewMorphGrid()
	m.AddPoints(3, 2, image.Point{1, 2}, image.Point{3, 4})
	m.AddPoints(3, 4, image.Point{2, 2}, image.Point{4, 6})
	m.AddPoints(4, 4, image.Point{5, 7}, image.Point{6, 10})
	g := m.interpolatedGrid(LinearInterpolationImagePoints, 0.5)
	AssertEqualsInt(t, g.verticalGridlineCount(), m.VerticalGridlineCount(), "Interpolated/Morph grid mismatch")
	AssertEqualsInt(t, g.horizontalGridlineCount(), m.HorizontalGridlineCount(), "Interpolated/Morph grid mismatch")
	p1, err := g.point(3, 2)
	if err != nil {
		LogVerbose(t, err.Error())
		t.Fail()
	}
	p2, err := g.point(3, 4)
	if err != nil {
		LogVerbose(t, err.Error())
		t.Fail()
	}
	p3, err := g.point(4, 4)
	if err != nil {
		LogVerbose(t, err.Error())
		t.Fail()
	}
	AssertEqualsFloat64PointTolerance(t, p1, Float64Point{2.0, 3.0}, 0.000001, "Interpolation point incorrect")
	AssertEqualsFloat64PointTolerance(t, p2, Float64Point{3.0, 4.0}, 0.000001, "Interpolation point incorrect")
	AssertEqualsFloat64PointTolerance(t, p3, Float64Point{5.5, 8.5}, 0.000001, "Interpolation point incorrect")
	g = m.interpolatedGrid(LinearInterpolationImagePoints, 0.25)
	AssertEqualsInt(t, g.verticalGridlineCount(), m.VerticalGridlineCount(), "Interpolated/Morph grid mismatch")
	AssertEqualsInt(t, g.horizontalGridlineCount(), m.HorizontalGridlineCount(), "Interpolated/Morph grid mismatch")
	p1, err = g.point(3, 2)
	if err != nil {
		LogVerbose(t, err.Error())
		t.Fail()
	}
	p2, err = g.point(3, 4)
	if err != nil {
		LogVerbose(t, err.Error())
		t.Fail()
	}
	p3, err = g.point(4, 4)
	if err != nil {
		LogVerbose(t, err.Error())
		t.Fail()
	}
	AssertEqualsFloat64PointTolerance(t, p1, Float64Point{1.5, 2.5}, 0.000001, "Interpolation point incorrect")
	AssertEqualsFloat64PointTolerance(t, p2, Float64Point{2.5, 3}, 0.000001, "Interpolation point incorrect")
	AssertEqualsFloat64PointTolerance(t, p3, Float64Point{5.25, 7.75}, 0.000001, "Interpolation point incorrect")
}
