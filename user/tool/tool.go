package tool

import (
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	. "github.com/digisan/go-generics/v2"
	fd "github.com/digisan/gotk/file-dir"
)

func ListField(objects ...any) (fields []string) {
	for _, obj := range objects {
		fields = append(fields, Fields(obj)...)
	}
	return
}

func ListValidator(objects ...any) (tags []string) {
	for _, obj := range objects {
		tags = append(tags, ValidatorTags(obj, "required", "email")...)
	}
	return Settify(tags...)
}

func LoadImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func Roi4Rgba(img image.Image, left, top, right, bottom int) *image.RGBA {
	rect := image.Rect(0, 0, right-left, bottom-top)
	rgba := image.NewRGBA(rect)
	draw.Draw(rgba, rect, img, image.Point{left, top}, draw.Src)
	return rgba
}

func SaveJPG(img image.Image, path string) (image.Image, error) {
	out, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	defer out.Close()

	var opts jpeg.Options
	opts.Quality = 100
	if err := jpeg.Encode(out, img, &opts); err != nil {
		return nil, err
	}
	return img, nil
}

func SavePNG(img image.Image, path string) (image.Image, error) {
	out, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	defer out.Close()

	if err := png.Encode(out, img); err != nil {
		return nil, err
	}
	return img, nil
}

// note must be 'crop:x,y,w,h'
func CropImage(fPath, note, outFmt string) (fCrop string, err error) {
	x, y, w, h := 0, 0, 0, 0
	if n, err := fmt.Sscanf(note, "crop:%d,%d,%d,%d", &x, &y, &w, &h); err == nil && n == 4 {
		img, err := LoadImage(fPath)
		if err != nil {
			return "", err
		}

		roi := Roi4Rgba(img, x, y, x+w, y+h)
		fCrop = fd.ChangeFileName(fPath, "", "-crop")
		fCrop = strings.TrimSuffix(fCrop, filepath.Ext(fCrop))

		switch outFmt {
		case ".png", "png":
			fCrop += ".png"
			if _, err := SavePNG(roi, fCrop); err != nil {
				return "", err
			}
		case ".jpg", "jpg":
			fCrop += ".jpg"
			if _, err := SaveJPG(roi, fCrop); err != nil {
				return "", err
			}
		default:
			fCrop += ".png"
			if _, err := SavePNG(roi, fCrop); err != nil {
				return "", err
			}
		}
		return fCrop, nil
	}
	return "", errors.New("note must be 'crop:x,y,w,h' to crop image")
}
