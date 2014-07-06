package gorph

import (
	"errors"
	"image"
	"math"
)

// LinearInterpolation linearly interpolates the Float64Point between two image points.
func LinearInterpolation(start, end image.Point, fractionFromStart float64) Float64Point {
	interpX := float64(start.X)*(1-fractionFromStart) + float64(end.X)*fractionFromStart
	interpY := float64(start.Y)*(1-fractionFromStart) + float64(end.Y)*fractionFromStart
	return Float64Point{interpX, interpY}
}

// CubicCatmullRomInterpolationImagePoints computes the Catmull-Rom spline from a given set
// of image points. There must be at least three points passed in, the alpha value must
// be in the range of [0.0, 1.0] and the total steps must be 2 or greater. The
// alpha parameter dictates the kind of Catmull-Rom spline generated; a value of
// 0 yields a Uniform curve, a value of 0.5 yields a Centripetal curve (which will
// not form loops), and a value of 1.0 creates a Chordal curve.
func CubicCatmullRomInterpolationImagePoints(points []image.Point, alpha float64, totSteps int) ([]Float64Point, error) {
	floatPoints := make([]Float64Point, len(points))
	for i := 0; i < len(floatPoints); i++ {
		floatPoints[i] = Float64Point{float64(points[i].X), float64(points[i].Y)}
	}
	return CubicCatmullRomInterpolation(floatPoints, alpha, totSteps)
}

// CubicCatmullRomInterpolation computes the Catmull-Rom spline from a given set
// of points. There must be at least three points passed in, the alpha value must
// be in the range of [0.0, 1.0] and the total steps must be 2 or greater. The
// alpha parameter dictates the kind of Catmull-Rom spline generated; a value of
// 0 yields a Uniform curve, a value of 0.5 yields a Centripetal curve (which will
// not form loops), and a value of 1.0 creates a Chordal curve.
func CubicCatmullRomInterpolation(points []Float64Point, alpha float64, totSteps int) ([]Float64Point, error) {
	nPoints := len(points)
	if nPoints < 3 {
		return nil, errors.New("CubicCatmullRomInterpolation: Two or less points passed in")
	}
	if alpha < 0 || alpha > 1 {
		return nil, errors.New("CubicCatmullRomInterpolation: Alpha must be in the range of [0.0, 1.0]")
	}
	if totSteps < 2 {
		return nil, errors.New("CubicCatmullRomInterpolation: Total steps must be 2 or greater")
	}
	// Precompute the number of "steps" to take between each point
	stepsAtRange := make([]uint, nPoints-1)
	sumDist := 0.0
	for i := 0; i < nPoints-1; i++ {
		sumDist += Distance(points[i], points[i+1])
	}
	for i := 0; i < nPoints-1; i++ {
		stepsAtRange[i] = uint(math.Floor(float64(totSteps)*Distance(points[i], points[i+1])/sumDist + 0.5))
	}

	// Linearly extrapolate
	postpendControl := Float64Point{points[nPoints-1].X + 2*(points[nPoints-1].X-points[nPoints-2].X), points[nPoints-1].Y + 2*(points[nPoints-1].Y-points[nPoints-2].Y)}
	pt0 := Float64Point{points[0].X - 2*(points[1].X-points[0].X), points[0].Y - 2*(points[1].Y-points[0].Y)} // "prependControl"
	pt1 := points[0]
	pt2 := points[1] // Want to iterate until this is nPoints (inclusive)
	pt3 := points[2]
	tPrev := 0.0
	tStart := math.Pow(Distance(pt0, pt1), alpha)
	tEnd := tStart + math.Pow(Distance(pt1, pt2), alpha)
	tNext := tEnd + math.Pow(Distance(pt2, pt3), alpha)

	resultPts := make([]Float64Point, 0, totSteps)
	for i := 0; i < nPoints-1; i++ {
		var j uint = 0
		for ; j < stepsAtRange[i]; j++ {
			// Use Barry and Goldman's pyramid to interpolate
			t := tStart + (float64(j)*(tEnd-tStart))/float64(stepsAtRange[i])
			L01 := Float64Point{float64(pt0.X)*((tStart-t)/(tStart-tPrev)) + float64(pt1.X)*((t-tPrev)/(tStart-tPrev)), float64(pt0.Y)*((tStart-t)/(tStart-tPrev)) + float64(pt1.Y)*((t-tPrev)/(tStart-tPrev))}
			L12 := Float64Point{float64(pt1.X)*((tEnd-t)/(tEnd-tStart)) + float64(pt2.X)*((t-tStart)/(tEnd-tStart)), float64(pt1.Y)*((tEnd-t)/(tEnd-tStart)) + float64(pt2.Y)*((t-tStart)/(tEnd-tStart))}
			L23 := Float64Point{float64(pt2.X)*((tNext-t)/(tNext-tEnd)) + float64(pt3.X)*((t-tEnd)/(tNext-tEnd)), float64(pt2.Y)*((tNext-t)/(tNext-tEnd)) + float64(pt3.Y)*((t-tEnd)/(tNext-tEnd))}
			L012 := Float64Point{float64(L01.X)*((tEnd-t)/(tEnd-tPrev)) + float64(L12.X)*((t-tPrev)/(tEnd-tPrev)), float64(L01.Y)*((tEnd-t)/(tEnd-tPrev)) + float64(L12.Y)*((t-tPrev)/(tEnd-tPrev))}
			L123 := Float64Point{float64(L12.X)*((tNext-t)/(tNext-tStart)) + float64(L23.X)*((t-tStart)/(tNext-tStart)), float64(L12.Y)*((tNext-t)/(tNext-tStart)) + float64(L23.Y)*((t-tStart)/(tNext-tStart))}
			C12 := Float64Point{float64(L012.X)*((tEnd-t)/(tEnd-tStart)) + float64(L123.X)*((t-tStart)/(tEnd-tStart)), float64(L012.Y)*((tEnd-t)/(tEnd-tStart)) + float64(L123.Y)*((t-tStart)/(tEnd-tStart))}
			resultPts = append(resultPts, C12)
		}
		// Iterate over next set of points in curve
		pt0 = pt1
		pt1 = pt2
		pt2 = pt3
		if i+3 >= nPoints {
			pt3 = postpendControl
		} else {
			pt3 = points[i+3]
		}
		// Update t parameters for next iteration
		tPrev = tStart
		tStart = tEnd
		tEnd = tNext
		tNext = tNext + math.Pow(Distance(pt2, pt3), alpha)
	}
	return resultPts, nil
}

// Distance computes the distance between two floating-point points.
func Distance(p1, p2 Float64Point) float64 {
	return math.Pow(math.Pow(p1.X-p2.X, 2.0)+math.Pow(p1.Y-p2.Y, 2.0), 0.5)
}

// DistanceImagePoint computes the distance between two integer points.
func DistanceImagePoint(p1, p2 image.Point) float64 {
	return math.Pow(float64(p1.X*p1.X+p2.Y*p2.Y), 0.5)
}
