package gorph

import (
	"errors"
	"math"
	"sort"
	"strconv"
)

type sortedFloat64Line struct {
	sortedPoints []Float64Point
	sortOnX      bool
}

func newSortedFloat64Line(sortHorizontally bool) *sortedFloat64Line {
	return &sortedFloat64Line{nil, sortHorizontally}
}

func (s *sortedFloat64Line) Len() int {
	return len(s.sortedPoints)
}

func (s *sortedFloat64Line) Less(i, j int) bool {
	if s.sortOnX {
		return s.sortedPoints[i].X < s.sortedPoints[j].X
	} else {
		return s.sortedPoints[i].Y < s.sortedPoints[j].Y
	}
}

func (s *sortedFloat64Line) Swap(i, j int) {
	temp := s.sortedPoints[j]
	s.sortedPoints[j] = s.sortedPoints[i]
	s.sortedPoints[i] = temp
}

func (s *sortedFloat64Line) AddPoint(pt Float64Point) int {
	s.sortedPoints = append(s.sortedPoints, pt)
	sort.Sort(s)
	return sort.Search(len(s.sortedPoints), func(i int) bool {
		return s.sortedPoints[i].X >= pt.X && s.sortedPoints[i].Y >= pt.Y
	})
}

func (s *sortedFloat64Line) AddPoints(pts []Float64Point) {
	for i := 0; i < len(pts); i++ {
		s.AddPoint(pts[i])
	}
}

func (s *sortedFloat64Line) PointWithXValue(xValue float64) (Float64Point, error) {
	if 0 == len(s.sortedPoints) {
		return Float64Point{0, 0}, errors.New("PointWithXValue: Line has no points.")
	}
	index := sort.Search(len(s.sortedPoints), func(i int) bool {
		return s.sortedPoints[i].X >= xValue
	})
	if index == len(s.sortedPoints) || s.sortedPoints[index].X != xValue {
		return Float64Point{0, 0}, errors.New("PointWithXValue: No point with value = " + strconv.FormatFloat(xValue, 'g', 8, 64))
	}
	return s.sortedPoints[index], nil
}

func (s *sortedFloat64Line) ClosestPointToSortedValue(value float64) (Float64Point, error) {
	if 0 == len(s.sortedPoints) {
		return Float64Point{0, 0}, errors.New("ClosestPointWithYValue: Line has no points.")
	}
	minVal := -1.0
	index := sort.Search(len(s.sortedPoints), func(i int) bool {
		var temp float64 = 0.0
		if s.sortOnX {
			temp = math.Abs(s.sortedPoints[i].X - value)
		} else {
			temp = math.Abs(s.sortedPoints[i].Y - value)
		}
		if minVal == -1.0 || temp < minVal {
			minVal = temp
			return true
		}
		return false
	})
	return s.sortedPoints[index], nil
}

func (s *sortedFloat64Line) RemovePoint(index int) error {
	if index < 0 || index >= s.Len() {
		return errors.New("RemovePoint: index to remove is out of bounds.")
	}
	copy(s.sortedPoints[index:], s.sortedPoints[index+1:])
	s.sortedPoints[len(s.sortedPoints)-1] = Float64Point{0, 0}
	s.sortedPoints = s.sortedPoints[:len(s.sortedPoints)-1]
	return nil
}

func (s *sortedFloat64Line) RemovePointWithXValue(xValue float64) error {
	index := sort.Search(len(s.sortedPoints), func(i int) bool {
		return s.sortedPoints[i].X >= xValue
	})
	if index == len(s.sortedPoints) || s.sortedPoints[index].X != xValue {
		return errors.New("RemovePointWithXValue: No point with value = " + strconv.FormatFloat(xValue, 'f', -1, 64))
	}
	copy(s.sortedPoints[index:], s.sortedPoints[index+1:])
	s.sortedPoints[len(s.sortedPoints)-1] = Float64Point{0, 0}
	s.sortedPoints = s.sortedPoints[:len(s.sortedPoints)-1]
	return nil
}

func (s *sortedFloat64Line) Point(index int) (Float64Point, error) {
	if index < 0 || index >= s.Len() {
		return Float64Point{0, 0}, errors.New("Point: index to remove is out of bounds.")
	}
	return s.sortedPoints[index], nil
}

func (s *sortedFloat64Line) HasPoints() bool {
	return len(s.sortedPoints) > 0
}
