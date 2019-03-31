package resize

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/monkeydioude/hresize/pkg/exif"
	"github.com/monkeydioude/hresize/pkg/format"
	"github.com/nfnt/resize"
)

const (
	defaultInterpolationFunction = resize.Lanczos3
)

type DiskEntity interface {
	Apply(OutputConfig) error
}

type Directory struct {
	Path string
}

func NewDirectory(path string) *Directory {
	return &Directory{
		Path: path,
	}
}

func (d *Directory) Apply(o OutputConfig) error {
	files, err := ioutil.ReadDir(d.Path)
	if err != nil {
		return err
	}

	for _, f := range files {

		if f.IsDir() {
			continue
		}

		err = (NewFile(d.Path, f.Name(), f.Size())).Apply(o)
		if err != nil {
			log.Printf("[WARN] %s\n", err)
		}
	}

	return nil
}

type File struct {
	Name string
	Path string
	Size int64
}

func NewFile(path, name string, size int64) *File {
	return &File{
		Name: name,
		Path: path,
		Size: size,
	}
}

func (f *File) Apply(o OutputConfig) error {
	fi, err := os.Open(fmt.Sprintf("%s/%s", f.Path, f.Name))
	defer fi.Close()

	if err != nil {
		return err
	}

	b := make([]byte, f.Size)
	n, err := fi.Read(b)
	if n <= 0 {
		return errors.New("Could not read any byte from file")
	}
	if err != nil {
		return err
	}

	exifData, err := exif.GetExifFromFile(fi)
	if err != nil {
		return err
	}

	imgFormat, err := format.GetImage(fi)
	if err != nil {
		return err
	}

	c, err := imgFormat.DecodeConfig(fi)
	if err != nil {
		return err
	}

	o.SetImageConfig(c)
	o.SetName(f.Name)

	i, err := imgFormat.Decode(fi)

	if err != nil {
		return err
	}

	imgOutput := ResizeOutputImage(o, i, defaultInterpolationFunction)

	if o.ShouldRotate() {
		imgOutput, err = RotateImg(exifData, imgOutput)
		if err != nil {
			return err
		}
	}

	return WriteOutputImage(o, imgOutput, imgFormat)
}
