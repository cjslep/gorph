package gorph

import (
	"errors"
	"image"
)

type float64CoordinateGrid struct {
	xGridLines   []*sortedFloat64Line
	yIndexLines  []*sortedLine
	nXGridLines  int
	nYIndexLines int
}

func newFloat64CoordinateGrid() *float64CoordinateGrid {
	return &float64CoordinateGrid{nil, nil, 0, 0}
}

func (f *float64CoordinateGrid) VerticalGridlineCount() int {
	return f.nXGridLines
}

func (f *float64CoordinateGrid) HorizontalGridlineCount() int {
	return f.nYIndexLines
}

func (f *float64CoordinateGrid) verticalGridlineLen() int {
	return len(f.xGridLines)
}

func (f *float64CoordinateGrid) horizontalGridlineLen() int {
	return len(f.yIndexLines)
}

func (f *float64CoordinateGrid) AddPoint(horizLine, vertLine int, pt Float64Point) {
	for i := len(f.xGridLines); i <= vertLine; i++ {
		f.xGridLines = append(f.xGridLines, newSortedFloat64Line(true))
	}
	for i := len(f.yIndexLines); i <= horizLine; i++ {
		f.yIndexLines = append(f.yIndexLines, newSortedLine(true))
	}
	if !f.xGridLines[vertLine].HasPoints() {
		f.nXGridLines++
	}
	index := f.xGridLines[vertLine].AddPoint(pt)
	if !f.yIndexLines[horizLine].HasPoints() {
		f.nYIndexLines++
	}
	f.yIndexLines[horizLine].AddPoint(image.Point{vertLine, index})
}

func (f *float64CoordinateGrid) RemovePoint(horizLine, vertLine int) error {
	err := f.checkBounds(horizLine, vertLine)
	if err != nil {
		return err
	}
	yIndexPoint, err := f.yIndexLines[horizLine].PointWithXValue(vertLine)
	if err != nil {
		return err
	}
	err = f.xGridLines[vertLine].RemovePoint(yIndexPoint.Y)
	if err != nil {
		return err
	}
	if !f.xGridLines[vertLine].HasPoints() {
		f.nXGridLines--
	}
	err = f.yIndexLines[horizLine].RemovePointWithXValue(vertLine)
	if err != nil {
		return err
	}
	if !f.yIndexLines[horizLine].HasPoints() {
		f.nYIndexLines--
	}
	for i := 0; i < len(f.yIndexLines); i++ {
		pt, err := f.yIndexLines[i].PointWithXValue(vertLine)
		if err == nil && pt.Y > yIndexPoint.Y {
			err := f.yIndexLines[i].RemovePointWithXValue(vertLine)
			if err == nil {
				f.yIndexLines[i].AddPoint(image.Point{pt.X, pt.Y - 1})
			}
		}
	}
	return nil
}

func (f *float64CoordinateGrid) Point(horizLine, vertLine int) (Float64Point, error) {
	err := f.checkBounds(horizLine, vertLine)
	if err != nil {
		return Float64Point{0, 0}, err
	}
	yIndexPoint, err := f.yIndexLines[horizLine].PointWithXValue(vertLine)
	if err != nil {
		return Float64Point{0, 0}, err
	}
	return f.xGridLines[vertLine].Point(yIndexPoint.Y)
}

func (f *float64CoordinateGrid) checkBounds(horizLine, vertLine int) error {
	if horizLine < 0 || horizLine >= f.horizontalGridlineLen() {
		return errors.New("checkBounds: horizLine is out of bounds.")
	} else if vertLine < 0 || vertLine >= f.verticalGridlineLen() {
		return errors.New("checkBounds: vertLine is out of bounds.")
	}
	return nil
}

func (f *float64CoordinateGrid) HorizontalLine(index int) (source []Float64Point) {
	source = make([]Float64Point, 0, f.verticalGridlineLen())
	for vLine := 0; vLine < f.verticalGridlineLen(); vLine++ {
		temp, err := f.Point(index, vLine)
		if err == nil {
			source = append(source, temp)
		}
	}
	return
}

func (f *float64CoordinateGrid) VerticalLine(index int) (source []Float64Point) {
	source = make([]Float64Point, 0, f.horizontalGridlineLen())
	for hLine := 0; hLine < f.horizontalGridlineLen(); hLine++ {
		temp, err := f.Point(hLine, index)
		if err == nil {
			source = append(source, temp)
		}
	}
	return
}

func (f *float64CoordinateGrid) AllCubicCatmullRomSplines(vertical bool, alpha float64, totSteps int) (splines []*sortedFloat64Line, nSplines int, err error) {
	splines = nil
	nSplines = 0
	err = nil
	nLoops := f.horizontalGridlineLen()
	if vertical {
		nLoops = f.verticalGridlineLen()
	}
	for i := 0; i < nLoops; i++ {
		var sourcePts []Float64Point = nil
		if vertical {
			sourcePts = f.VerticalLine(i)
		} else {
			sourcePts = f.HorizontalLine(i)
		}
		if len(sourcePts) > 2 {
			sourceLine, err := CubicCatmullRomInterpolation(sourcePts, alpha, totSteps)
			if err != nil {
				return nil, 0, err
			}
			splines = append(splines, newSortedFloat64Line(!vertical))
			splines[len(splines)-1].AddPoints(sourceLine)
			nSplines++
		}
	}
	return
}