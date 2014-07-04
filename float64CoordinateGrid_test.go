package gorph

import (
	"testing"
)

func TestFloat64CoordinateGridAddPoints(t *testing.T) {
	m := newFloat64CoordinateGrid()
	m.addPoint(3, 2, Float64Point{1.1, 2.2})
	p1, err := m.point(3, 2)
	if err != nil {
		LogVerbose(t, err.Error())
		t.Fail()
	}
	AssertEqualsFloat64Point(t, p1, Float64Point{1.1, 2.2})
}

func TestFloat64CoordinateGridLineCount(t *testing.T) {
	m := newFloat64CoordinateGrid()
	m.addPoint(3, 2, Float64Point{1.1, 24.55})
	m.addPoint(3, 4, Float64Point{2.2456, 2.66})
	length := m.horizontalGridlineCount()
	AssertEqualsInt(t, length, 1, "HorizontalGridlineCount failed")
	length = m.verticalGridlineCount()
	AssertEqualsInt(t, length, 2, "VerticalGridlineCount failed")
}

func TestFloat64CoordinateGridPoint(t *testing.T) {
	m := newFloat64CoordinateGrid()
	m.addPoint(3, 2, Float64Point{1.1, 24.55})
	m.addPoint(3, 4, Float64Point{2.2456, 2.66})
	pt, err := m.point(3, 2)
	if err != nil {
		LogVerbose(t, err.Error())
		t.Fail()
	}
	AssertEqualsFloat64Point(t, pt, Float64Point{1.1, 24.55})
}

func TestFloat64CoordinateGridHorizontalLine(t *testing.T) {
	m := newFloat64CoordinateGrid()
	m.addPoint(3, 2, Float64Point{1.1, 24.55})
	m.addPoint(3, 4, Float64Point{2.2456, 2.66})
	source := m.horizontalLine(3)
	AssertEqualsInt(t, len(source), 2, "Length incorrect")
	AssertEqualsFloat64Point(t, source[0], Float64Point{1.1, 24.55})
	AssertEqualsFloat64Point(t, source[1], Float64Point{2.2456, 2.66})
}

func TestFloat64CoordinateGridVerticalLine(t *testing.T) {
	m := newFloat64CoordinateGrid()
	m.addPoint(3, 2, Float64Point{1.1, 24.55})
	m.addPoint(3, 4, Float64Point{2.2456, 2.66})
	source := m.verticalLine(2)
	AssertEqualsInt(t, len(source), 1, "Length incorrect")
	AssertEqualsFloat64Point(t, source[0], Float64Point{1.1, 24.55})
}

func TestFloat64CoordinateGridRemovePoints(t *testing.T) {
	m := newFloat64CoordinateGrid()
	m.addPoint(3, 2, Float64Point{1.1, 24.55})
	m.addPoint(3, 4, Float64Point{2.2456, 2.66})
	m.addPoint(4, 4, Float64Point{5.55, 10.101})
	source := m.horizontalLine(3)
	AssertEqualsInt(t, len(source), 2, "Before removal length incorrect")
	err := m.removePoint(3, 2)
	if err != nil {
		LogVerbose(t, err.Error())
		t.Fail()
	}
	source = m.horizontalLine(3)
	AssertEqualsInt(t, len(source), 1, "After removal length incorrect")
}

func TestFloat64CoordinateGridRemovePointsVertical(t *testing.T) {
	m := newFloat64CoordinateGrid()
	m.addPoint(3, 2, Float64Point{1.1, 24.55})
	m.addPoint(3, 4, Float64Point{2.2456, 2.66})
	m.addPoint(4, 4, Float64Point{5.55, 10.101})
	source := m.verticalLine(4)
	AssertEqualsInt(t, len(source), 2, "Before removal length incorrect")
	err := m.removePoint(3, 4)
	if err != nil {
		LogVerbose(t, err.Error())
		t.Fail()
	}
	source = m.verticalLine(4)
	AssertEqualsInt(t, len(source), 1, "After removal length incorrect")
}
