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
// timeInterp - function to use to interpolate cross-fading grids over time
// nominalTimeConversion - function to covert actual time frame of grid to nominal time used
// in cross fading. The parameter and returned value must lie in the range [0.0, 1.0]
func Morph(numMorphs int, start, dest image.Image, mGrid MorphGrid, timeInterp InterpolationFunc, nominalTimeConversion func(float64) float64) ([]image.Image, error) {
	startBounds := start.Bounds()
	destBounds := dest.Bounds()
	if !startBounds.Min.Eq(destBounds.Min) || !startBounds.Max.Eq(destBounds.Max) {
		return nil, errors.New("Morph: image bounds do not match")
	}
	// go?
	results := make([]image.Image, 0, numMorphs)
	for i := 1; i <= numMorphs; i++ {
		baseTimeFrac := float64(i) / float64(numMorphs+1)
		intermedGrid := mGrid.interpolatedGrid(timeInterp, baseTimeFrac)
		auxGridSource := newFloat64CoordinateGrid()
		auxGridDest := newFloat64CoordinateGrid()

		// TODO: Factory creation function?
		auxSourceImage := image.NewRGBA64(startBounds)
		auxDestImage := image.NewRGBA64(destBounds)
		intermedSourceImage := image.NewRGBA64(startBounds)
		intermedDestImage := image.NewRGBA64(startBounds)

		// TODO: Do this somewhere else?
		maxX := intermedGrid.verticalGridlineCount()
		maxY := intermedGrid.horizontalGridlineCount()
		for x := 0; x < maxX; x++ {
			for y := 0; y < maxY; y++ {
				auxYPt, err := intermedGrid.point(y, x)
				if err != nil {
					return nil, err
				}
				auxXSourcePt, auxXDestPt, err := mGrid.Points(y, x)
				if err != nil {
					return nil, err
				}
				auxGridSource.addPoint(y, x, Float64Point{float64(auxXSourcePt.X), auxYPt.Y})
				auxGridDest.addPoint(y, x, Float64Point{float64(auxXDestPt.X), auxYPt.Y})
			}
		}

		// Calculate Cubic Catmull-Rom spline equations for each vertical line in
		//   both original (source, dest) and aux (source, dest) images
		sourceOriginalSplines, destOriginalSplines, nSplinesGrid, err := mGrid.allCubicCatmullRomSplines(true, 0.5, startBounds.Max.Y-startBounds.Min.Y)
		if err != nil {
			return nil, err
		}
		sourceAuxSplines, nSplinesAuxSource, err := auxGridSource.allCubicCatmullRomSplines(true, 0.5, startBounds.Max.Y-startBounds.Min.Y)
		if err != nil {
			return nil, err
		}
		if nSplinesGrid != nSplinesAuxSource {
			return nil, errors.New("Given MorphGrid and source auxilary grid do not have the same number of splines.")
		}
		destAuxSplines, nSplinesDestSource, err := auxGridDest.allCubicCatmullRomSplines(true, 0.5, startBounds.Max.Y-startBounds.Min.Y)
		if err != nil {
			return nil, err
		}
		if nSplinesGrid != nSplinesDestSource {
			return nil, errors.New("Given MorphGrid and destination auxilary grid do not have the same number of splines.")
		}

		err = stretchPixelsHorizontally(startBounds.Min.Y, startBounds.Max.Y, sourceOriginalSplines, sourceAuxSplines, start, auxSourceImage)
		if err != nil {
			return nil, err
		}
		err = stretchPixelsHorizontally(startBounds.Min.Y, startBounds.Max.Y, destOriginalSplines, destAuxSplines, dest, auxDestImage)
		if err != nil {
			return nil, err
		}

		// Auxiliary to intermediate, stretching vertically
		sourceAuxSplines, nSplinesAuxSource, err = auxGridSource.allCubicCatmullRomSplines(false, 0.5, startBounds.Max.X-startBounds.Min.X)
		if err != nil {
			return nil, err
		}
		destAuxSplines, nSplinesDestSource, err = auxGridDest.allCubicCatmullRomSplines(false, 0.5, startBounds.Max.X-startBounds.Min.X)
		if err != nil {
			return nil, err
		}
		intermedSplines, nIntermedSplines, err := intermedGrid.allCubicCatmullRomSplines(false, 0.5, startBounds.Max.X-startBounds.Min.X)
		if nIntermedSplines != nSplinesAuxSource {
			return nil, errors.New("Auxilary grid and intermediate grid do not have the same number of splines.")
		}

		err = stretchPixelsVertically(startBounds.Min.X, startBounds.Max.X, sourceAuxSplines, intermedSplines, auxSourceImage, intermedSourceImage)
		if err != nil {
			return nil, err
		}
		err = stretchPixelsVertically(startBounds.Min.X, startBounds.Max.X, destAuxSplines, intermedSplines, auxDestImage, intermedDestImage)
		if err != nil {
			return nil, err
		}

		// Cross dissolve the two intermediate (source, dest) images by
		//   using a weight (weight depends on i).
		results[i-1], err = CrossDissolve([]image.Image{intermedSourceImage, intermedDestImage}, []float64{nominalTimeConversion(baseTimeFrac), 1 - nominalTimeConversion(baseTimeFrac)})
	}
	return results, nil
}

// CrossDissolve weights a series of images on a pixel-by-pixel basis in order to
// produce a resulting image. Returns an error if any of the image bounds do not
// match, if one or no images are provided, or the number of images do not match
// the number of weights.
func CrossDissolve(dissolving []image.Image, weights []float64) (image.Image, error) {
	nImages := len(dissolving)
	nWeights := len(weights)
	if nImages != nWeights {
		return nil, errors.New("CrossDissolve: Number of images to dissolve does not match the number of weights")
	}
	if nImages <= 1 {
		return nil, errors.New("CrossDissolve: Two or more images must be provided")
	}
	startBounds := dissolving[0].Bounds()
	for i := 1; i < nImages; i++ {
		if !startBounds.Min.Eq(dissolving[i].Bounds().Min) || !startBounds.Max.Eq(dissolving[i].Bounds().Max) {
			return nil, errors.New("CrossDissolve: Image bounds do not match")
		}
	}
	result := image.NewRGBA64(startBounds)
	// go?
	for x := startBounds.Min.X; x < startBounds.Max.X; x++ {
		for y := startBounds.Min.Y; y < startBounds.Max.Y; y++ {
			colorToSet := weightColor(dissolving[0].At(x, y), weights[0])
			for i := 1; i < nImages; i++ {
				colorToSet = addColors(colorToSet, weightColor(dissolving[i].At(x, y), weights[i]))
			}
			result.Set(x, y, colorToSet)
		}
	}
	return result, nil
}

func stretchPixelsHorizontally(yStart, yEnd int, originalSplines, auxSplines []*parametricLineFloat64, start image.Image, final *image.RGBA64) error {
	nSplines := len(originalSplines)
	if nSplines != len(auxSplines) {
		return errors.New("stretchPixelsHorizontally: Spline count does not match between start and final images")
	}
	for y := yStart; y < yEnd; y++ {
		for iSpline := 0; iSpline < nSplines-1; iSpline++ {
			origStart, err := originalSplines[iSpline].InterpolatePointsAtY(float64(y))
			if err != nil {
				return err
			}
			origEnd, err := originalSplines[iSpline+1].InterpolatePointsAtY(float64(y))
			if err != nil {
				return err
			}
			destStart, err := auxSplines[iSpline].InterpolatePointsAtY(float64(y))
			if err != nil {
				return err
			}
			destEnd, err := auxSplines[iSpline+1].InterpolatePointsAtY(float64(y))
			if err != nil {
				return err
			}
			if len(origStart) != 1 || len(origEnd) != 1 || len(destStart) != 1 || len(destEnd) != 1 {
				return errors.New("stretchPixelsHorizontally: Invalid spline length (folds back on itself, or no length)")
			}
			mergePixelsInLine(true, y, iSpline != 0, iSpline != nSplines-1, origStart[0].X, origEnd[0].X, destStart[0].X, destEnd[0].X, start, final)
		}
	}
	return nil
}

func stretchPixelsVertically(xStart, xEnd int, originalSplines, auxSplines []*parametricLineFloat64, start image.Image, final *image.RGBA64) error {
	nSplines := len(originalSplines)
	if nSplines != len(auxSplines) {
		return errors.New("stretchPixelsVertically: Spline count does not match between start and final images")
	}
	for x := xStart; x < xEnd; x++ {
		for iSpline := 0; iSpline < nSplines-1; iSpline++ {
			origStart, err := originalSplines[iSpline].InterpolatePointsAtX(float64(x))
			if err != nil {
				return err
			}
			origEnd, err := originalSplines[iSpline+1].InterpolatePointsAtX(float64(x))
			if err != nil {
				return err
			}
			destStart, err := auxSplines[iSpline].InterpolatePointsAtX(float64(x))
			if err != nil {
				return err
			}
			destEnd, err := auxSplines[iSpline+1].InterpolatePointsAtX(float64(x))
			if err != nil {
				return err
			}
			if len(origStart) > 1 || len(origEnd) > 1 || len(destStart) > 1 || len(destEnd) > 1 {
				return errors.New("stretchPixelsVertically: Spline folds back on itself")
			}
			mergePixelsInLine(false, x, iSpline != 0, iSpline != nSplines-1, origStart[0].Y, origEnd[0].Y, destStart[0].Y, destEnd[0].Y, start, final)
		}
	}
	return nil
}

func mergePixelsInLine(horizontally bool, line int, fadeStartPixel, fadeEndPixel bool, origStart, origEnd, destStart, destEnd float64, original image.Image, dest *image.RGBA64) {
	pixelOrigSnapStart := int(math.Floor(origStart)) + 1
	pixelOrigSnapEnd := int(math.Floor(origEnd)) + 1
	var origColor color.Color
	lastColoredDestPixel := int(math.Floor(destStart))
	fmt.Printf("**mergePixelsInLine**\norigStart=%v, origEnd=%v, destStart=%v, destEnd=%v\n", origStart, origEnd, destStart, destEnd)
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
