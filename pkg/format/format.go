package format

import (
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"os"
)

var JpegQuality = 98

type Image interface {
	Decode(*os.File) (image.Image, error)
	Encode(*os.File, image.Image) error
	DecodeConfig(*os.File) (image.Config, error)
}

func GetType(file *os.File) string {
	b := make([]byte, 4)
	n, _ := file.ReadAt(b, 0)
	if n < 4 {
		return ""
	}

	if b[0] == 0x89 && b[1] == 0x50 && b[2] == 0x4E && b[3] == 0x47 {
		return "png"
	}
	if b[0] == 0xFF && b[1] == 0xD8 {
		return "jpg"
	}
	if b[0] == 0x47 && b[1] == 0x49 && b[2] == 0x46 && b[3] == 0x38 {
		return "gif"
	}
	if b[0] == 0x42 && b[1] == 0x4D {
		return "bmp"
	}
	return ""
}

func GetImage(file *os.File) (Image, error) {
	t := GetType(file)

	if t == "" {
		return nil, errors.New("Could not define image type")
	}

	switch t {
	case "jpg":
		return NewJpeg(), nil
	case "png":
		return &Png{}, nil
	}

	return nil, errors.New("Could not get Image instance")
}

type Jpeg struct {
	Options *jpeg.Options
}

func NewJpeg() *Jpeg {
	return &Jpeg{
		Options: &jpeg.Options{
			Quality: JpegQuality,
		},
	}
}

func (j *Jpeg) Decode(file *os.File) (image.Image, error) {
	file.Seek(0, 0)
	return jpeg.Decode(file)
}

func (j *Jpeg) Encode(file *os.File, img image.Image) error {
	file.Seek(0, 0)
	return jpeg.Encode(file, img, j.Options)
}

func (j *Jpeg) DecodeConfig(file *os.File) (image.Config, error) {
	file.Seek(0, 0)
	return jpeg.DecodeConfig(file)
}

type Png struct {
}

func (p *Png) Decode(file *os.File) (image.Image, error) {
	return png.Decode(file)
}

func (p *Png) Encode(file *os.File, img image.Image) error {
	return png.Encode(file, img)
}

func (p *Png) DecodeConfig(file *os.File) (image.Config, error) {
	return png.DecodeConfig(file)
}
