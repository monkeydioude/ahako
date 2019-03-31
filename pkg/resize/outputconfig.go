package resize

import (
	"fmt"
	"image"
	"math"
)

type OutputConfig interface {
	GetDimensions() (uint, uint)
	GetSourceDimensions() (uint, uint)
	GetDir() string
	GetName() string
	SetName(string)
	SetImageConfig(image.Config)
	ShouldRotate() bool
	fmt.Stringer
}

type Ratio struct {
	Ratio       float64
	Dir         string
	Name        string
	imageConfig image.Config
	Rotate      bool
}

func NewRatio(ratio float64, outputdir string, exif bool, rotate bool) *Ratio {
	return &Ratio{
		Ratio:  ratio,
		Dir:    outputdir,
		Rotate: rotate,
	}
}

func (r *Ratio) GetDimensions() (uint, uint) {
	w := uint(math.Round(float64(r.imageConfig.Width) * r.Ratio))
	h := uint(math.Round(float64(r.imageConfig.Height) * r.Ratio))
	return w, h
}

func (r *Ratio) GetSourceDimensions() (uint, uint) {
	return uint(r.imageConfig.Width), uint(r.imageConfig.Height)
}

func (r *Ratio) SetImageConfig(c image.Config) {
	r.imageConfig = c
}

func (r *Ratio) GetDir() string {
	return r.Dir
}

func (r *Ratio) GetName() string {
	return r.Name
}

func (r *Ratio) SetName(n string) {
	r.Name = n
}

func (r *Ratio) String() string {
	return fmt.Sprintf("%s/%s", r.Dir, r.Name)
}

func (r *Ratio) ShouldRotate() bool {
	return r.Rotate
}
