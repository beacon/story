package story

import (
	"bytes"
	"image"
)

type Image struct {
	Node
	image image.Image
}

func NewImage(imageBytes []byte) (*Image, error) {
	// Decode an image from the image file's byte slice.
	img, _, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		return nil, err
	}
	return &Image{
		image: img,
	}, nil
}