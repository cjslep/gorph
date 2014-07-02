package gorph

import (
	"image"
)

//
// Gridlines:
//  +------> x
//  |
//  |          -- Horizontal Gridlines
//  |
// \/          --
//  y          --
//    |  |  |
//    Vertical Gridlines

type MorphGrid struct {
	start *CoordinateGrid
	dest  *CoordinateGrid
}

func NewMorphGrid() *MorphGrid {
	return &MorphGrid{NewCoordinateGrid(), NewCoordinateGrid()}
}

func (m *MorphGrid) AddPoints(horizLine, vertLine int, startPt, destPt image.Point) {
	m.start.AddPoint(horizLine, vertLine, startPt)
	m.dest.AddPoint(horizLine, vertLine, destPt)
}

func (m *MorphGrid) RemovePoints(horizLine, vertLine int) error {
	err := m.start.RemovePoint(horizLine, vertLine)
	if err != nil {
		return err
	}
	err = m.dest.RemovePoint(horizLine, vertLine)
	return err
}

func (m *MorphGrid) VerticalGridlineCount() int {
	return m.start.VerticalGridlineCount()
}

func (m *MorphGrid) HorizontalGridlineCount() int {
	return m.start.HorizontalGridlineCount()
}

func (m *MorphGrid) Points(horizLine, vertLine int) (image.Point, image.Point, error) {
	startPt, err := m.start.Point(horizLine, vertLine)
	if err != nil {
		return startPt, image.ZP, err
	}
	endPt, err := m.dest.Point(horizLine, vertLine)
	return startPt, endPt, err
}

func (m *MorphGrid) HorizontalLine(index int) (source, dest []image.Point) {
	source = make([]image.Point, 0, m.start.verticalGridlineLen())
	dest = make([]image.Point, 0, m.start.verticalGridlineLen())
	for vLine := 0; vLine < m.start.verticalGridlineLen(); vLine++ {
		temp, err := m.start.Point(index, vLine)
		if err == nil {
			source = append(source, temp)
		}
		temp, err = m.dest.Point(index, vLine)
		if err == nil {
			dest = append(dest, temp)
		}
	}
	return
}

func (m *MorphGrid) VerticalLine(index int) (source, dest []image.Point) {
	source = make([]image.Point, 0, m.start.horizontalGridlineLen())
	dest = make([]image.Point, 0, m.start.horizontalGridlineLen())
	for hLine := 0; hLine < m.start.horizontalGridlineLen(); hLine++ {
		temp, err := m.start.Point(hLine, index)
		if err == nil {
			source = append(source, temp)
		}
		temp, err = m.dest.Point(hLine, index)
		if err == nil {
			dest = append(dest, temp)
		}
	}
	return
}

func (m *MorphGrid) AllCubicCatmullRomSplines(vertical bool, alpha float64, totSteps int) (source, dest []*sortedFloat64Line, nSplines int, err error) {
	source = nil
	dest = nil
	nSplines = 0
	err = nil
	nLoops := m.start.horizontalGridlineLen()
	if vertical {
		nLoops = m.start.verticalGridlineLen()
	}
	for i := 0; i < nLoops; i++ {
		var sourcePts []image.Point = nil
		var destPts []image.Point = nil
		if vertical {
			sourcePts, destPts = m.VerticalLine(i)
		} else {
			sourcePts, destPts = m.HorizontalLine(i)
		}
		if len(sourcePts) > 2 && len(destPts) > 2 {
			sourceLine, err := CubicCatmullRomInterpolationImagePoints(sourcePts, alpha, totSteps)
			if err != nil {
				return nil, nil, 0, err
			}
			destLine, err := CubicCatmullRomInterpolationImagePoints(destPts, alpha, totSteps)
			if err != nil {
				return nil, nil, 0, err
			}
			source = append(source, newSortedFloat64Line(!vertical))
			dest = append(dest, newSortedFloat64Line(!vertical))
			source[len(source)-1].AddPoints(sourceLine)
			dest[len(dest)-1].AddPoints(destLine)
			nSplines++
		}
	}
	return
}

func (m *MorphGrid) interpolatedGrid(interpFn InterpolationFunc, fractionFromStart float64) *float64CoordinateGrid {
	interpGrid := newFloat64CoordinateGrid()
	maxX := MaxInt(m.start.verticalGridlineLen(), m.dest.verticalGridlineLen())
	maxY := MaxInt(m.start.horizontalGridlineLen(), m.dest.horizontalGridlineLen())
	for x := 0; x < maxX; x++ {
		for y := 0; y < maxY; y++ {
			start, end, err := m.Points(y, x)
			if err == nil {
				resultPt := interpFn(start, end, fractionFromStart)
				interpGrid.AddPoint(y, x, resultPt)
			}
		}
	}
	return interpGrid
}
