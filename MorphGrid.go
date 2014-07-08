package gorph

import (
	"image"
)

// MorphGrid is a metadata structure usually used alongside images for specifying
// parameters used in transformations. The MorphGrid consists of two grids in order
// to hold data for a pre-morph and post-morph state. It manages these pairs of
// points so accompanying morphing algorithms can easily use the parametric information
// between homogulous lines. For details on how a MorphGrid is exactly used, please
// refer to a morphing algorithm's documentation.
type MorphGrid struct {
	start *coordinateGrid
	dest  *coordinateGrid
}

// NewMorphGrid supplies a new instance of a MorphGrid.
func NewMorphGrid() *MorphGrid {
	return &MorphGrid{newCoordinateGrid(), newCoordinateGrid()}
}

// AddPoints adds two homogulous points for a before and after image on the
// specified horizontal and vertical line indices. Replaces any preexisting
// points.
func (m *MorphGrid) AddPoints(horizLine, vertLine int, startPt, destPt image.Point) {
	m.start.addPoint(horizLine, vertLine, startPt)
	m.dest.addPoint(horizLine, vertLine, destPt)
}

// RemovePoints removes the homogulous points that belong to the specified
// horizontal and vertical line. Returns an error if the operation is not
// able to complete successfully.
func (m *MorphGrid) RemovePoints(horizLine, vertLine int) error {
	err := m.start.removePoint(horizLine, vertLine)
	if err != nil {
		return err
	}
	err = m.dest.removePoint(horizLine, vertLine)
	return err
}

// VerticalGridlineCount determines the number of vertical grid lines that have
// been specified.
func (m *MorphGrid) VerticalGridlineCount() int {
	return m.start.verticalGridlineCount()
}

// HorizontalGridlineCount determines the number of horizontal grid lines that
// have been specified.
func (m *MorphGrid) HorizontalGridlineCount() int {
	return m.start.horizontalGridlineCount()
}

// Points returns the homogulous pair of points at the intersection of the two
// lines given by their indices.
func (m *MorphGrid) Points(horizLine, vertLine int) (image.Point, image.Point, error) {
	startPt, err := m.start.point(horizLine, vertLine)
	if err != nil {
		return startPt, image.ZP, err
	}
	endPt, err := m.dest.point(horizLine, vertLine)
	return startPt, endPt, err
}

// HorizontalLine takes an index of a horizontal line and returns all points
// associated with the line in both grids. The points are sorted in increasing
// x-values.
func (m *MorphGrid) HorizontalLine(index int) (source, dest []image.Point) {
	source = make([]image.Point, 0, m.start.verticalGridlineLen())
	dest = make([]image.Point, 0, m.start.verticalGridlineLen())
	for vLine := 0; vLine < m.start.verticalGridlineLen(); vLine++ {
		temp, err := m.start.point(index, vLine)
		if err == nil {
			source = append(source, temp)
		}
		temp, err = m.dest.point(index, vLine)
		if err == nil {
			dest = append(dest, temp)
		}
	}
	return
}

// VerticalLine takes an index of a vertical line and returns all points
// associated with the line in both grids. The points are sorted in increasing
// y-values.
func (m *MorphGrid) VerticalLine(index int) (source, dest []image.Point) {
	source = make([]image.Point, 0, m.start.horizontalGridlineLen())
	dest = make([]image.Point, 0, m.start.horizontalGridlineLen())
	for hLine := 0; hLine < m.start.horizontalGridlineLen(); hLine++ {
		temp, err := m.start.point(hLine, index)
		if err == nil {
			source = append(source, temp)
		}
		temp, err = m.dest.point(hLine, index)
		if err == nil {
			dest = append(dest, temp)
		}
	}
	return
}

// allCubicCatmullRomSplines
func (m *MorphGrid) allCubicCatmullRomSplines(vertical bool, alpha float64, totSteps int) (source, dest []*parametricLineFloat64, nSplines int, err error) {
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
			source = append(source, newParametricLineFloat64())
			dest = append(dest, newParametricLineFloat64())
			source[len(source)-1].AddPoints(sourceLine)
			dest[len(dest)-1].AddPoints(destLine)
			nSplines++
		}
	}
	return
}

// interpolatedGrid
func (m *MorphGrid) interpolatedGrid(interpFn InterpolationFunc, fractionFromStart float64) *float64CoordinateGrid {
	interpGrid := newFloat64CoordinateGrid()
	maxX := MaxInt(m.start.verticalGridlineLen(), m.dest.verticalGridlineLen())
	maxY := MaxInt(m.start.horizontalGridlineLen(), m.dest.horizontalGridlineLen())
	for x := 0; x < maxX; x++ {
		for y := 0; y < maxY; y++ {
			start, end, err := m.Points(y, x)
			if err == nil {
				resultPt := interpFn(start, end, fractionFromStart)
				interpGrid.addPoint(y, x, resultPt)
			}
		}
	}
	return interpGrid
}
