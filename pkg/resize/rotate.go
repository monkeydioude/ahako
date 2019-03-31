package resize

import (
	"fmt"
	"image"
	"image/color"

	"github.com/disintegration/imaging"
	"github.com/monkeydioude/ahako/pkg/exif"
)

var rotateMap = map[uint16]float64{
	8: 90,
}

func RotateImg(e *exif.Exif, imgOutput image.Image) (image.Image, error) {
	o := e.GetOrientation()

	if o == 1 {
		return imgOutput, nil
	}

	if _, ok := rotateMap[o]; !ok {
		return nil, fmt.Errorf("[WARN] %d is not an available orientation", o)
	}

	return imaging.Rotate(imgOutput, rotateMap[o], color.Black), nil
}
