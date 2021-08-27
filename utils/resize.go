package utils

import (
	"bytes"
	"github.com/disintegration/imaging"
	"image"
	"io"
)

func ResizeImage(file io.Reader, width, height int) (*[]byte, error) {
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	dstImage := imaging.Resize(img, width, height, imaging.CatmullRom)
	buff := new(bytes.Buffer)
	if err := imaging.Encode(buff, dstImage, imaging.JPEG); err != nil {
		return nil, err
	}

	imgBytes := buff.Bytes()

	return &imgBytes, nil
}
