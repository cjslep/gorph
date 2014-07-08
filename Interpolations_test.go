package gorph

import (
	"testing"
)

func TestCentripetalCubicCatmullRomInterpolationThreePoints(t *testing.T) {
	alpha := 0.5
	totSteps := 30
	points := make([]Float64Point, 0)
	points = append(points, Float64Point{0, 0})
	points = append(points, Float64Point{1, 1})
	points = append(points, Float64Point{2, 0})
	pts, err := CubicCatmullRomInterpolation(points, alpha, totSteps)
	if err != nil {
		LogVerbose(t, err.Error())
		t.Fail()
	}
	AssertEqualsInt(t, len(pts), totSteps+1, "Interpolation step length failed")
	AssertEqualsFloat64PointTolerance(t, pts[0], Float64Point{0, 0}, .000001, "Point 0 incorrect")
	AssertEqualsFloat64PointTolerance(t, pts[1], Float64Point{0.07663060253, 0.08077875068}, .000001, "Point 1 incorrect")
	AssertEqualsFloat64PointTolerance(t, pts[2], Float64Point{0.1505160391, 0.1659234465}, .000001, "Point 2 incorrect")
	AssertEqualsFloat64PointTolerance(t, pts[7], Float64Point{0.4894413772, 0.6055895254}, .000001, "Point 7 incorrect")
	AssertEqualsFloat64PointTolerance(t, pts[14], Float64Point{0.934045043, 0.9921191171}, .000001, "Point 14 incorrect")
	AssertEqualsFloat64PointTolerance(t, pts[15], Float64Point{1, 1}, .000001, "Middle point incorrect")
	AssertEqualsFloat64PointTolerance(t, pts[16], Float64Point{1.065954957, 0.9921191171}, .000001, "Point 16 incorrect")
	AssertEqualsFloat64PointTolerance(t, pts[29], Float64Point{1.923369397, 0.08077875068}, .000001, "Point 29 incorrect")
	LogVerbose(t, pts)
}

func TestCentripetalCubicCatmullRomInterpolationThreePointsSimple(t *testing.T) {
	alpha := 0.5
	totSteps := 5
	points := make([]Float64Point, 0)
	points = append(points, Float64Point{0, 0})
	points = append(points, Float64Point{3, 0})
	points = append(points, Float64Point{4, 0})
	pts, err := CubicCatmullRomInterpolation(points, alpha, totSteps)
	if err != nil {
		LogVerbose(t, err.Error())
		t.Fail()
	}
	AssertEqualsInt(t, len(pts), totSteps, "Interpolation step length failed")
	AssertEqualsFloat64PointTolerance(t, pts[0], Float64Point{0, 0}, .000001, "Point 0 incorrect")
	AssertEqualsFloat64PointTolerance(t, pts[3], Float64Point{3, 0}, .000001, "Point 3 incorrect")
	AssertEqualsFloat64PointTolerance(t, pts[4], Float64Point{4, 0}, .000001, "Point 4 incorrect")
}

func TestCentripetalCubicCatmullRomInterpolationFourPoints(t *testing.T) {
	alpha := 0.5
	totSteps := 30
	points := make([]Float64Point, 0)
	points = append(points, Float64Point{0, 0})
	points = append(points, Float64Point{0.75, 0.5})
	points = append(points, Float64Point{1.25, 0.5})
	points = append(points, Float64Point{2, 1})
	pts, err := CubicCatmullRomInterpolation(points, alpha, totSteps)
	if err != nil {
		LogVerbose(t, err.Error())
		t.Fail()
	}
	AssertEqualsInt(t, len(pts), totSteps+1, "Interpolation step length failed")
	AssertEqualsFloat64PointTolerance(t, pts[0], Float64Point{0, 0}, .000001, "Point 0 incorrect")
	AssertEqualsFloat64PointTolerance(t, pts[1], Float64Point{0.07179755079, 0.0494979255}, .000001, "Point 1 incorrect")
	AssertEqualsFloat64PointTolerance(t, pts[11], Float64Point{0.6914760989, 0.478945874}, .000001, "Point 11 incorrect")
	AssertEqualsFloat64PointTolerance(t, pts[12], Float64Point{0.75, 0.5}, .000001, "Point 12 incorrect")
	AssertEqualsFloat64PointTolerance(t, pts[13], Float64Point{0.83564892, 0.51471849}, .000001, "Point 13 incorrect")
	AssertEqualsFloat64PointTolerance(t, pts[18], Float64Point{1.25, 0.5}, .000001, "Point 18 incorrect")
	AssertEqualsFloat64PointTolerance(t, pts[19], Float64Point{1.308523901, 0.521054126}, .000001, "Point 19 incorrect")
	AssertEqualsFloat64PointTolerance(t, pts[20], Float64Point{1.3668031, 0.548179856}, .000001, "Point 20 incorrect")
	AssertEqualsFloat64PointTolerance(t, pts[30], Float64Point{2, 1}, .000001, "Point 30 incorrect")
	LogVerbose(t, pts)
}
