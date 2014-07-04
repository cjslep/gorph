package gorph

import (
	"errors"
	"image"
)

type coordinateGrid struct {
	xGridLines   []*sortedLine
	yIndexLines  []*sortedLine
	nXGridLines  int
	nYIndexLines int
}

func newCoordinateGrid() *coordinateGrid {
	return &coordinateGrid{nil, nil, 0, 0}
}

func (c *coordinateGrid) verticalGridlineCount() int {
	return c.nXGridLines
}

func (c *coordinateGrid) horizontalGridlineCount() int {
	return c.nYIndexLines
}

func (c *coordinateGrid) verticalGridlineLen() int {
	return len(c.xGridLines)
}

func (c *coordinateGrid) horizontalGridlineLen() int {
	return len(c.yIndexLines)
}

func (c *coordinateGrid) addPoint(horizLine, vertLine int, pt image.Point) {
	for i := len(c.xGridLines); i <= vertLine; i++ {
		c.xGridLines = append(c.xGridLines, newSortedLine(true))
	}
	for i := len(c.yIndexLines); i <= horizLine; i++ {
		c.yIndexLines = append(c.yIndexLines, newSortedLine(true))
	}
	if !c.xGridLines[vertLine].HasPoints() {
		c.nXGridLines++
	}
	index := c.xGridLines[vertLine].AddPoint(pt)
	if !c.yIndexLines[horizLine].HasPoints() {
		c.nYIndexLines++
	}
	c.yIndexLines[horizLine].AddPoint(image.Point{vertLine, index})
}

func (c *coordinateGrid) removePoint(horizLine, vertLine int) error {
	err := c.checkBounds(horizLine, vertLine)
	if err != nil {
		return err
	}
	yIndexPoint, err := c.yIndexLines[horizLine].PointWithXValue(vertLine)
	if err != nil {
		return err
	}
	err = c.xGridLines[vertLine].RemovePoint(yIndexPoint.Y)
	if err != nil {
		return err
	}
	if !c.xGridLines[vertLine].HasPoints() {
		c.nXGridLines--
	}
	err = c.yIndexLines[horizLine].RemovePointWithXValue(vertLine)
	if err != nil {
		return err
	}
	if !c.yIndexLines[horizLine].HasPoints() {
		c.nYIndexLines--
	}
	for i := 0; i < len(c.yIndexLines); i++ {
		pt, err := c.yIndexLines[i].PointWithXValue(vertLine)
		if err == nil && pt.Y > yIndexPoint.Y {
			err := c.yIndexLines[i].RemovePointWithXValue(vertLine)
			if err == nil {
				c.yIndexLines[i].AddPoint(image.Point{pt.X, pt.Y - 1})
			}
		}
	}
	return nil
}

func (c *coordinateGrid) point(horizLine, vertLine int) (image.Point, error) {
	err := c.checkBounds(horizLine, vertLine)
	if err != nil {
		return image.ZP, err
	}
	yIndexPoint, err := c.yIndexLines[horizLine].PointWithXValue(vertLine)
	if err != nil {
		return image.ZP, err
	}
	return c.xGridLines[vertLine].Point(yIndexPoint.Y)
}

func (c *coordinateGrid) checkBounds(horizLine, vertLine int) error {
	if horizLine < 0 || horizLine >= c.horizontalGridlineLen() {
		return errors.New("checkBounds: horizLine is out of bounds.")
	} else if vertLine < 0 || vertLine >= c.verticalGridlineLen() {
		return errors.New("checkBounds: vertLine is out of bounds.")
	}
	return nil
}
