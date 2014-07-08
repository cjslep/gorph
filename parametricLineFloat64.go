package gorph

import (
	"errors"
	"strconv"
	"fmt"
)

type parametricLineFloat64 struct {
	parametricPoints []Float64Point
}

func newParametricLineFloat64() *parametricLineFloat64 {
	return &parametricLineFloat64{nil}
}

func (p *parametricLineFloat64) Len() int {
	return len(p.parametricPoints)
}

func (p *parametricLineFloat64) AddPoint(pt Float64Point) int {
	p.parametricPoints = append(p.parametricPoints, pt)
	return len(p.parametricPoints) - 1
}

func (p *parametricLineFloat64) AddPoints(pts []Float64Point) {
	for i := 0; i < len(pts); i++ {
		p.AddPoint(pts[i])
	}
}

func (p *parametricLineFloat64) RemovePoint(index int) error {
	if index < 0 || index >= p.Len() {
		return errors.New("RemovePoint: index to remove is out of bounds.")
	}
	copy(p.parametricPoints[index:], p.parametricPoints[index+1:])
	p.parametricPoints[len(p.parametricPoints)-1] = Float64Point{0, 0}
	p.parametricPoints = p.parametricPoints[:len(p.parametricPoints)-1]
	return nil
}

func (p *parametricLineFloat64) Point(index int) (Float64Point, error) {
	if index < 0 || index >= p.Len() {
		return Float64Point{0, 0}, errors.New("Point: index to remove is out of bounds.")
	}
	return p.parametricPoints[index], nil
}

func (p *parametricLineFloat64) HasPoints() bool {
	return len(p.parametricPoints) > 0
}

func (p *parametricLineFloat64) InterpolatePointsAtX(xValue float64) (points []Float64Point, err error) {
	if p.Len() < 2 {
		return nil, errors.New("InterpolatePointsAtX: Line has fewer than 2 points.")
	}
	points = nil
	err = nil
	for i := 1; i < len(p.parametricPoints); i++ {
		if p.parametricPoints[i].X >= xValue  && p.parametricPoints[i-1].X <= xValue{
			points = append(points, LinearInterpolation(p.parametricPoints[i-1], p.parametricPoints[i], (xValue - p.parametricPoints[i-1].X)/(p.parametricPoints[i].X - p.parametricPoints[i-1].X)))
			fmt.Printf("InterpolatePointsAtX, appended: (%v,%v)\n", points[len(points)-1].X, points[len(points)-1].Y)
		}
	}
	if len(points) == 0 {
		err = errors.New("InterpolatePointsAtX: No points interpolated for value = " + strconv.FormatFloat(xValue, 'g', 8, 64))
	}
	return
}

func (p *parametricLineFloat64) InterpolatePointsAtY(yValue float64) (points []Float64Point, err error) {
	if p.Len() < 2 {
		return nil, errors.New("InterpolatePointsAtY: Line has fewer than 2 points.")
	}
	points = nil
	err = nil
	for i := 1; i < len(p.parametricPoints); i++ {
		if p.parametricPoints[i].Y >= yValue  && p.parametricPoints[i-1].Y <= yValue {
			fmt.Printf("Linearly interpolating %v and %v over %v", p.parametricPoints[i-1], p.parametricPoints[i], (yValue - p.parametricPoints[i-1].Y)/(p.parametricPoints[i].Y - p.parametricPoints[i-1].Y))
			points = append(points, LinearInterpolation(p.parametricPoints[i-1], p.parametricPoints[i], (yValue - p.parametricPoints[i-1].Y)/(p.parametricPoints[i].Y - p.parametricPoints[i-1].Y)))
			fmt.Printf("InterpolatePointsAtY, appended for i=%v: (%v,%v)\n", i, points[len(points)-1].X, points[len(points)-1].Y)
		}
	}
	if len(points) == 0 {
		err = errors.New("InterpolatePointsAtY: No points interpolated for value = " + strconv.FormatFloat(yValue, 'g', 8, 64))
	}
	return
}
