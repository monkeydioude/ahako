package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/monkeydioude/hresize/pkg/resize"
)

const (
	defaultOutputDir = "resized"
	outputDirRights  = 0766
)

func getCwd() string {
	ex, err := os.Executable()

	if err != nil {
		log.Fatal("[ERR ] Could not resolve os.Executable")
	}

	return filepath.Dir(ex)
}

func getDiskEntity(dir *string) resize.DiskEntity {
	if dir != nil && *dir != "" {
		return resize.NewDirectory(*dir)
	}
	return nil
}

func getOutputDir(outputDir *string) string {
	if outputDir != nil && *outputDir != "" {
		return *outputDir
	}
	return getCwd() + "/" + defaultOutputDir + "/" + strconv.Itoa(int(time.Now().UnixNano()))
}

func getOutputConfig(outputDir string, r *float64, rexif *bool, rotate *bool) resize.OutputConfig {
	if r != nil && *r > 0 {
		return resize.NewRatio(*r, outputDir, !(*rexif), *rotate)
	}
	return nil
}

func main() {
	d := flag.String("d", "", "Directory of photos to be resized")
	o := flag.String("o", "", "Directory to output photos")
	r := flag.Float64("r", 0, "Define the ratio of resizing")
	rotate := flag.Bool("auto-rotate", false, "Automatic rotation of the picture using EXIF metadata")
	exif := flag.Bool("remove-exif", false, "Remove output EXIF metadata")
	m := flag.Uint("m", outputDirRights, "FileMode for output directory")
	flag.Parse()

	diskEntity := getDiskEntity(d)
	if diskEntity == nil {
		log.Fatal("[ERR ] Could not define any disk entity from flags")
	}

	outputDir := getOutputDir(o)

	err := os.MkdirAll(outputDir, os.FileMode(*m))
	if err != nil {
		log.Println(err)
	}

	outputConfig := getOutputConfig(outputDir, r, exif, rotate)
	if outputConfig == nil {
		log.Fatal("[ERR ] Could not define any output config from flags")
	}

	err = diskEntity.Apply(outputConfig)

	if err != nil {
		log.Fatal(err)
	}
}
