package gorph

import (
	"errors"
	"image"
	"image/color"
	"math"
)

type InterpolationFunc func(start, end image.Point, fractionFromStart float64) Float64Point

//
// numMorphs - the number of morph images to create
// start - starting image
// dest - resulting image
// mGrid - two grids with start/end points
// timeInterp - function to use to interpolate time-wise
//
func Morph(numMorphs int, start, dest image.Image, mGrid MorphGrid, timeInterp InterpolationFunc) error {
	startBounds := start.Bounds()
	destBounds := dest.Bounds()
	if !startBounds.Min.Eq(destBounds.Min) || !startBounds.Max.Eq(destBounds.Max) {
		return errors.New("Morph: image bounds do not match")
	}
	// go?
	for i := 1; i <= numMorphs; i++ {
		baseTimeFrac := float64(i) / float64(numMorphs+1)
		intermedGrid := mGrid.interpolatedGrid(timeInterp, baseTimeFrac)
		auxGridSource := newFloat64CoordinateGrid()
		auxGridDest := newFloat64CoordinateGrid()

		// TODO: Factory creation function
		auxSourceImage := image.NewNRGBA64(startBounds)
		auxDestImage := image.NewNRGBA64(destBounds)

		// TODO: Do this somewhere else
		maxX := intermedGrid.VerticalGridlineCount()
		maxY := intermedGrid.HorizontalGridlineCount()
		for x := 0; x < maxX; x++ {
			for y := 0; y < maxY; y++ {
				auxYPt, err := intermedGrid.Point(y, x)
				if err != nil {
					return err
				}
				auxXSourcePt, auxXDestPt, err := mGrid.Points(y, x)
				if err != nil {
					return err
				}
				auxGridSource.AddPoint(y, x, Float64Point{float64(auxXSourcePt.X), auxYPt.Y})
				auxGridDest.AddPoint(y, x, Float64Point{float64(auxXDestPt.X), auxYPt.Y})
			}
		}

		// Calculate Cubic Catmull-Rom spline equations for each vertical line in
		//   both original (source, dest) and aux (source, dest) images
		sourceOriginalSplines, destOriginalSplines, nSplinesGrid, err := mGrid.AllCubicCatmullRomSplines(true, 0.5, startBounds.Max.Y-startBounds.Min.Y)
		if err != nil {
			return err
		}
		sourceAuxSplines, nSplinesAuxSource, err := auxGridSource.AllCubicCatmullRomSplines(true, 0.5, startBounds.Max.Y-startBounds.Min.Y)
		if err != nil {
			return err
		}
		if nSplinesGrid != nSplinesAuxSource {
			return errors.New("Given MorphGrid and source auxilary grid do not have the same number of splines.")
		}
		destAuxSplines, nSplinesDestSource, err := auxGridDest.AllCubicCatmullRomSplines(true, 0.5, startBounds.Max.Y-startBounds.Min.Y)
		if err != nil {
			return err
		}
		if nSplinesGrid != nSplinesDestSource {
			return errors.New("Given MorphGrid and destination auxilary grid do not have the same number of splines.")
		}

		err = stretchPixelsHorizontally(startBounds.Min.Y, startBounds.Max.Y, sourceOriginalSplines, sourceAuxSplines, start, auxSourceImage)
		err = stretchPixelsHorizontally(startBounds.Min.Y, startBounds.Max.Y, destOriginalSplines, destAuxSplines, dest, auxDestImage)

		// Calculate Cubic Catmull-Rom spline equations for each horizontal line in
		//   both aux (source, dest) and intermediate (source, dest) images
		// Iterate over vertical rows in intermediate (source, dest) image
		for x := startBounds.Min.X; x < startBounds.Max.X; x++ {
			// For each line:
			//  Get intercept of spline for aux (source, dest) and intermediate (source, dest) image
			//  Map pixels from aux (source, dest) to intermediate (source, dest) image, using
			//    fractional weights as necessary
		}

		// Cross dissolve the two intermediate (source, dest) images by
		//   using a weight (weight depends on i).
		// Advanced: weight changes value non-linearly depending on i

	}
	return nil
}

func stretchPixelsHorizontally(yStart, yEnd int, originalSplines, auxSplines []*sortedFloat64Line, start image.Image, aux *image.NRGBA64) error {
	nSplinesGrid := len(originalSplines)
	for y := yStart; y < yEnd; y++ {
		// For each line:
		//  Get intercept of spline for original (source, dest) and aux (source, dest) image
		//  Map pixels from original (source, dest) to aux (source, dest) image, using
		//    fractional weights as necessary for antialiasing
		originalStart, err := originalSplines[0].ClosestPointToSortedValue(float64(y))
		if err != nil {
			return err
		}
		auxStart, err := auxSplines[0].ClosestPointToSortedValue(float64(y))
		if err != nil {
			return err
		}
		for splineIndex := 1; splineIndex < nSplinesGrid; splineIndex++ {
			originalEnd, err := originalSplines[splineIndex].ClosestPointToSortedValue(float64(y))
			if err != nil {
				return err
			}
			auxEnd, err := auxSplines[splineIndex].ClosestPointToSortedValue(float64(y))
			if err != nil {
				return err
			}
			deltaSourceOriginal := originalEnd.X - originalStart.X
			deltaSourceAux := auxEnd.X - auxStart.X
			normSourceDist := deltaSourceOriginal / deltaSourceAux
			if normSourceDist < 1 { // Expanding smaller pixels into larger
				normSourceDist = deltaSourceAux / deltaSourceOriginal
				auxPivot := auxStart.X
				for auxX := auxPivot; auxX < auxPivot+normSourceDist; auxX += 1 {
					for x := originalStart.X; x <= originalEnd.X; x += 1 {
						colorRes := start.At(int(math.Floor(x)), y)
						if (auxX == auxPivot && splineIndex != 1) || (auxX+1 >= auxPivot+normSourceDist && splineIndex+1 != nSplinesGrid) {
							prev := aux.At(int(math.Floor(auxX)), y)
							ratioThis := auxX - math.Floor(auxX)
							colorRes = interpolateColors(colorRes, prev, ratioThis)
							aux.Set(int(math.Floor(auxX)), y, colorRes)
						} else {
							aux.Set(int(math.Floor(auxX)), y, colorRes)
						}
					}
					auxPivot += normSourceDist
				}
			} else { // Expanding larger pixels into smaller
				// TODO: handle end pixels that need full color, border pixels
				sourcePivot := originalStart.X
				for x := auxStart.X; x <= auxEnd.X; x += 1 {
					/*for sourceX := sourcePivot; sourceX < sourcePivot + normSourceDist; sourceX += 1 {
						colorRes := start.At(int(math.Floor(sourcePivot)), y)
						if (sourceX == sourcePivot && splineIndex != 1) || (sourceX + 1 >= sourcePivot + normSourceDist && splineIndex + 1 != nSplinesGrid) {
							prev := aux.At(int(math.Floor(auxX)), y)
							ratioThis := auxX - math.Floor(auxX)
							colorRes = interpolateColors(colorRes, prev, ratioThis)
							aux.Set(int(math.Floor(auxX)), y, colorRes)
						} else {
							aux.Set(int(math.Floor(auxX)), y, colorRes)
						}
					}*/

					colorRes := start.At(int(math.Floor(sourcePivot)), y)
					nextColor := start.At(int(math.Floor(sourcePivot))+1, y)
					ratioNext := 0.0
					if math.Floor(sourcePivot+normSourceDist/2) > math.Floor(sourcePivot) {
						ratioNext = (sourcePivot + normSourceDist/2 - math.Floor(sourcePivot+normSourceDist)) / normSourceDist
					}
					colorRes = interpolateColors(nextColor, colorRes, ratioNext)
					aux.Set(int(math.Floor(x)), y, colorRes)
					sourcePivot += normSourceDist
				}
			}
			originalStart = originalEnd
			auxStart = auxEnd
		}
	}
	return nil
}

/*
func mergePixels(horizontally, fadeStartPixel, fadeEndPixel bool, origStart, origEnd, destStart, destEnd float64, original image.Image, dest *image.NRGBA64) {
	pixelOrigStart := int(math.Floor(origStart))
	pixelOrigEnd := int(math.Floor(origEnd))
	pixelDestStart := int(math.Floor(destStart))
	pixelDestEnd := int(math.Floor(destEnd))
}*/

func interpolateColors(colorWeighted, colorOther color.Color, weight float64) color.Color {
	r, g, b, a := colorWeighted.RGBA()
	rOther, gOther, bOther, aOther := colorOther.RGBA()
	rRes := addCeilingOverflow(multiplyCeilingOverflow(r, weight), multiplyCeilingOverflow(rOther, 1.0-weight))
	gRes := addCeilingOverflow(multiplyCeilingOverflow(g, weight), multiplyCeilingOverflow(gOther, 1.0-weight))
	bRes := addCeilingOverflow(multiplyCeilingOverflow(b, weight), multiplyCeilingOverflow(bOther, 1.0-weight))
	aRes := addCeilingOverflow(multiplyCeilingOverflow(a, weight), multiplyCeilingOverflow(aOther, 1.0-weight))
	return color.NRGBA64{rRes, gRes, bRes, aRes}
}

func multiplyCeilingOverflow(value uint32, weight float64) uint16 {
	ret := uint16(float64(value) * weight)
	if math.Floor(weight) != 0.0 && uint32(float64(ret)/math.Floor(weight)) != value {
		ret = 1<<16 - 1 // Overflow, return largest number possible
	}
	return ret
}

func addCeilingOverflow(value uint16, value2 uint16) uint16 {
	ret := value + value2
	if ret < value || ret < value2 {
		ret = 1<<16 - 1 // Overflow, return largest number possible
	}
	return ret
}
