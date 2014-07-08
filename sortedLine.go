package gorph

import (
	"errors"
	"image"
	"sort"
	"strconv"
)

type sortedLine struct {
	sortedPoints []image.Point
	sortOnX      bool
}

func newSortedLine(sortHorizontally bool) *sortedLine {
	return &sortedLine{nil, sortHorizontally}
}

func (s *sortedLine) Len() int {
	return len(s.sortedPoints)
}

func (s *sortedLine) Less(i, j int) bool {
	if s.sortOnX {
		return s.sortedPoints[i].X < s.sortedPoints[j].X || (s.sortedPoints[i].X == s.sortedPoints[j].X && s.sortedPoints[i].Y < s.sortedPoints[j].Y)
	} else {
		return s.sortedPoints[i].Y < s.sortedPoints[j].Y || (s.sortedPoints[i].Y == s.sortedPoints[j].Y && s.sortedPoints[i].X < s.sortedPoints[j].X)
	}
}

func (s *sortedLine) Swap(i, j int) {
	temp := s.sortedPoints[j]
	s.sortedPoints[j] = s.sortedPoints[i]
	s.sortedPoints[i] = temp
}

func (s *sortedLine) AddPoint(pt image.Point) int {
	s.sortedPoints = append(s.sortedPoints, pt)
	sort.Sort(s)
	return sort.Search(len(s.sortedPoints), func(i int) bool {
		if s.sortOnX {
			return s.sortedPoints[i].X > pt.X || (s.sortedPoints[i].X == pt.X && s.sortedPoints[i].Y >= pt.Y)
		} else {
			return s.sortedPoints[i].Y > pt.Y || (s.sortedPoints[i].Y == pt.Y && s.sortedPoints[i].X >= pt.X)
		}
	})
}

func (s *sortedLine) AddPoints(pts []image.Point) {
	for i := 0; i < len(pts); i++ {
		s.AddPoint(pts[i])
	}
}

func (s *sortedLine) PointWithXValue(xValue int) (image.Point, int, error) {
	index := sort.Search(len(s.sortedPoints), func(i int) bool {
		return s.sortedPoints[i].X >= xValue
	})
	if index == len(s.sortedPoints) || s.sortedPoints[index].X != xValue {
		return image.ZP, -1, errors.New("PointWithXValue: No point with value = " + strconv.Itoa(xValue))
	}
	return s.sortedPoints[index], index, nil
}

func (s *sortedLine) RemovePoint(index int) error {
	if index < 0 || index >= s.Len() {
		return errors.New("RemovePoint: index to remove is out of bounds.")
	}
	copy(s.sortedPoints[index:], s.sortedPoints[index+1:])
	s.sortedPoints[len(s.sortedPoints)-1] = image.ZP
	s.sortedPoints = s.sortedPoints[:len(s.sortedPoints)-1]
	return nil
}

func (s *sortedLine) RemovePointWithXValue(xValue int) error {
	index := sort.Search(len(s.sortedPoints), func(i int) bool {
		return s.sortedPoints[i].X >= xValue
	})
	if index == len(s.sortedPoints) || s.sortedPoints[index].X != xValue {
		return errors.New("RemovePointWithXValue: No point with value = " + strconv.Itoa(xValue))
	}
	copy(s.sortedPoints[index:], s.sortedPoints[index+1:])
	s.sortedPoints[len(s.sortedPoints)-1] = image.ZP
	s.sortedPoints = s.sortedPoints[:len(s.sortedPoints)-1]
	return nil
}

func (s *sortedLine) Point(index int) (image.Point, error) {
	if index < 0 || index >= s.Len() {
		return image.ZP, errors.New("Point: index is out of bounds: " + strconv.Itoa(index))
	}
	return s.sortedPoints[index], nil
}

func (s *sortedLine) HasPoints() bool {
	return len(s.sortedPoints) > 0
}
