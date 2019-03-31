package resize

import (
	"image"
	"log"
	"os"

	"github.com/monkeydioude/hresize/pkg/format"
	"github.com/nfnt/resize"
)

func ResizeOutputImage(o OutputConfig, i image.Image, ifunc resize.InterpolationFunction) image.Image {
	w, h := o.GetDimensions()

	m := resize.Resize(w, h, i, ifunc)
	sw, sh := o.GetSourceDimensions()
	log.Printf("[INFO] Resized %s from %dx%d to %dx%d\n", o.GetName(), sw, sh, w, h)
	return m
}

func WriteOutputImage(o OutputConfig, imgOutput image.Image, imgFormat format.Image) error {
	out, err := os.Create(o.String())
	if err != nil {
		return err
	}
	defer out.Close()

	err = imgFormat.Encode(out, imgOutput)

	if err != nil {
		return err
	}
	return nil
}
