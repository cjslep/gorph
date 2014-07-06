package gorph

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"math"
)

// InterpolationFunc interpolates a Float64Point between two image points based on a
// ratio distance from the starting point to the ending point, such that fractionFromStart
// lies on the interval [0.0, 1.0]. The value 0.0 maps to the starting point, and 1.0 to
// the end point.
type InterpolationFunc func(start, end image.Point, fractionFromStart float64) Float64Point

// Morph performes a keyframe-based morphing of two images in order to interpolate a new set
// of transition images. It is based on the coordinate grid approach to morphing an image, as
// opposed to a feature line.
// numMorphs - the number of morph images to create
// start - starting image
// dest - ending image
// mGrid - a MorphGrid representing homogulous points on both images
// timeInterp - function to use to interpolate cross-fading images over time
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
		auxSourceImage := image.NewRGBA64(startBounds)
		auxDestImage := image.NewRGBA64(destBounds)

		// TODO: Do this somewhere else
		maxX := intermedGrid.verticalGridlineCount()
		maxY := intermedGrid.horizontalGridlineCount()
		for x := 0; x < maxX; x++ {
			for y := 0; y < maxY; y++ {
				auxYPt, err := intermedGrid.point(y, x)
				if err != nil {
					return err
				}
				auxXSourcePt, auxXDestPt, err := mGrid.Points(y, x)
				if err != nil {
					return err
				}
				auxGridSource.addPoint(y, x, Float64Point{float64(auxXSourcePt.X), auxYPt.Y})
				auxGridDest.addPoint(y, x, Float64Point{float64(auxXDestPt.X), auxYPt.Y})
			}
		}

		// Calculate Cubic Catmull-Rom spline equations for each vertical line in
		//   both original (source, dest) and aux (source, dest) images
		sourceOriginalSplines, destOriginalSplines, nSplinesGrid, err := mGrid.allCubicCatmullRomSplines(true, 0.5, startBounds.Max.Y-startBounds.Min.Y)
		if err != nil {
			return err
		}
		sourceAuxSplines, nSplinesAuxSource, err := auxGridSource.allCubicCatmullRomSplines(true, 0.5, startBounds.Max.Y-startBounds.Min.Y)
		if err != nil {
			return err
		}
		if nSplinesGrid != nSplinesAuxSource {
			return errors.New("Given MorphGrid and source auxilary grid do not have the same number of splines.")
		}
		destAuxSplines, nSplinesDestSource, err := auxGridDest.allCubicCatmullRomSplines(true, 0.5, startBounds.Max.Y-startBounds.Min.Y)
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

func stretchPixelsHorizontally(yStart, yEnd int, originalSplines, auxSplines []*sortedFloat64Line, start image.Image, aux *image.RGBA64) error {
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

func mergePixelsInLine(horizontally bool, line int, fadeStartPixel, fadeEndPixel bool, origStart, origEnd, destStart, destEnd float64, original image.Image, dest *image.RGBA64) {
	pixelOrigSnapStart := int(math.Floor(origStart)) + 1
	pixelOrigSnapEnd := int(math.Floor(origEnd)) + 1
	var origColor color.Color
	lastColoredDestPixel := int(math.Floor(destStart))
	if !fadeStartPixel {
		lastColoredDestPixel--
	}
	for iOrig := pixelOrigSnapStart; iOrig <= pixelOrigSnapEnd; iOrig++ {
		fmt.Printf("iOrig=%v\n", iOrig)
		if horizontally {
			origColor = original.At(iOrig-1, line)
			fmt.Printf("(%v, %v) origColor: %v\n", (iOrig - 1), line, origColor)
		} else {
			origColor = original.At(line, iOrig-1)
			fmt.Printf("(%v, %v) origColor: %v\n", line, (iOrig - 1), origColor)
		}

		pct := (math.Min(float64(iOrig), origEnd) - origStart) / (origEnd - origStart)
		wOrig := 1.0
		if iOrig == pixelOrigSnapStart {
			wOrig = 1 - (origStart - math.Floor(origStart))
		} else if iOrig == pixelOrigSnapEnd {
			wOrig = origEnd - math.Floor(origEnd)
		}
		if wOrig > 0.0 {
			wDest := wOrig / (origEnd - origStart) * (destEnd - destStart)
			iEndDest := int(math.Floor(pct*(destEnd-destStart) + destStart))
			iStartDest := int(math.Floor(pct*(destEnd-destStart) + destStart - wDest))
			wDestFrac := 1 - (pct*(destEnd-destStart) + destStart - wDest - float64(iStartDest))
			fmt.Printf("wDest=%v, iEndDest=%v, iStartDest=%v, wDestFrac=%v\n", wDest, iEndDest, iStartDest, wDestFrac)
			for iDest := iStartDest; iDest <= iEndDest; iDest++ {
				if iDest == iEndDest && iStartDest != iEndDest {
					wDestFrac = pct*(destEnd-destStart) + destStart - float64(iEndDest)
				}
				if wDestFrac > 0 {
					if iDest > lastColoredDestPixel && (!fadeEndPixel || (fadeEndPixel && iOrig != pixelOrigSnapEnd)) {
						if horizontally {
							dest.Set(iDest, line, weightColor(origColor, wDestFrac))
							fmt.Printf("(%v, %v) destColor unadded (weight=%v): %v\n", iDest, line, wDestFrac, weightColor(origColor, wDestFrac))
						} else {
							dest.Set(line, iDest, weightColor(origColor, wDestFrac))
							fmt.Printf("(%v, %v) destColor unadded (weight=%v): %v\n", line, iDest, wDestFrac, weightColor(origColor, wDestFrac))
						}
						lastColoredDestPixel = iDest
					} else {
						pastColor := dest.At(iDest, line)
						if iDest > lastColoredDestPixel && fadeEndPixel && iOrig == pixelOrigSnapEnd {
							lastColoredDestPixel = iDest
						}
						if horizontally {
							dest.Set(iDest, line, addColors(pastColor, weightColor(origColor, wDestFrac)))
							fmt.Printf("(%v, %v) destColor added (weight=%v): %v\n", iDest, line, wDestFrac, addColors(pastColor, weightColor(origColor, wDestFrac)))
						} else {
							dest.Set(line, iDest, addColors(pastColor, weightColor(origColor, wDestFrac)))
							fmt.Printf("(%v, %v) destColor added (weight=%v): %v\n", line, iDest, wDestFrac, addColors(pastColor, weightColor(origColor, wDestFrac)))
						}
					}
				}
				wDestFrac = 1
			}
		}
	}
}

func weightColor(colorWeighted color.Color, weight float64) color.Color {
	r, g, b, a := colorWeighted.RGBA()
	rRes := multiplyCeilingOverflow(r, weight)
	gRes := multiplyCeilingOverflow(g, weight)
	bRes := multiplyCeilingOverflow(b, weight)
	aRes := multiplyCeilingOverflow(a, weight)
	return color.RGBA64{rRes, gRes, bRes, aRes}
}

func addColors(colorOne, colorTwo color.Color) color.Color {
	r, g, b, a := colorOne.RGBA()
	rOther, gOther, bOther, aOther := colorTwo.RGBA()
	rRes := uint32ToUint16CeilingOverflow(addCeilingOverflow32(r, rOther))
	gRes := uint32ToUint16CeilingOverflow(addCeilingOverflow32(g, gOther))
	bRes := uint32ToUint16CeilingOverflow(addCeilingOverflow32(b, bOther))
	aRes := uint32ToUint16CeilingOverflow(addCeilingOverflow32(a, aOther))
	return color.RGBA64{rRes, gRes, bRes, aRes}
}

func interpolateColors(colorWeighted, colorOther color.Color, weight float64) color.Color {
	r, g, b, a := colorWeighted.RGBA()
	rOther, gOther, bOther, aOther := colorOther.RGBA()
	rRes := addCeilingOverflow16(multiplyCeilingOverflow(r, weight), multiplyCeilingOverflow(rOther, 1.0-weight))
	gRes := addCeilingOverflow16(multiplyCeilingOverflow(g, weight), multiplyCeilingOverflow(gOther, 1.0-weight))
	bRes := addCeilingOverflow16(multiplyCeilingOverflow(b, weight), multiplyCeilingOverflow(bOther, 1.0-weight))
	aRes := addCeilingOverflow16(multiplyCeilingOverflow(a, weight), multiplyCeilingOverflow(aOther, 1.0-weight))
	return color.RGBA64{rRes, gRes, bRes, aRes}
}

func multiplyCeilingOverflow(value uint32, weight float64) uint16 {
	ret := uint16(float64(value)*weight + 0.5)
	if math.Floor(weight) != 0.0 && uint32(float64(ret)/math.Floor(weight)) != value {
		ret = 1<<16 - 1 // Overflow, return largest number possible
	}
	return ret
}

func addCeilingOverflow16(value uint16, value2 uint16) uint16 {
	ret := value + value2
	if ret < value || ret < value2 {
		ret = 1<<16 - 1 // Overflow, return largest number possible
	}
	return ret
}

func addCeilingOverflow32(value uint32, value2 uint32) uint32 {
	ret := value + value2
	if ret < value || ret < value2 {
		ret = 1<<32 - 1 // Overflow, return largest number possible
	}
	return ret
}

func uint32ToUint16CeilingOverflow(value uint32) uint16 {
	ret := uint16(value)
	if uint32(ret) != value {
		ret = 1<<16 - 1 // Overflow, return largest number possible
	}
	return ret
}
