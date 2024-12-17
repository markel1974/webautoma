package executor

import (
	"bytes"
	"golang.org/x/image/draw"
	"image"
	"image/png"
)

// AspectRatio adjusts the width and height to maintain the original aspect ratio within the target dimensions.
// Returns the new width and height as integers.
// If the source height is zero or less, the target width is used for both dimensions.
func AspectRatio(srcW int, srcH int, targetW int, targetH int) (int, int) {
	if srcH <= 0 {
		return targetW, targetW
	}
	r := float64(srcW) / float64(srcH)
	widthT := float64(targetH) * r
	heightT := widthT / r
	return int(widthT), int(heightT)
}

// Scale resizes an input image to the specified dimensions, optionally preserving the aspect ratio if enabled.
// buf is the image data to be scaled, x1 and y1 define the target width and height, and aspectRatio determines scaling behavior.
// Returns the scaled image as a PNG-encoded byte slice or an error if the operation fails.
func Scale(buf []byte, x1 int, y1 int, aspectRatio bool) ([]byte, error) {
	reader := bytes.NewReader(buf)
	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}
	//draw.NearestNeighbor
	//draw.ApproxBiLinear
	//draw.BiLinear
	//draw.CatmullRom
	sc := draw.BiLinear

	if aspectRatio {
		s := img.Bounds().Size()
		x1, y1 = AspectRatio(s.X, s.Y, x1, y1)
	}

	dr := image.Rect(0, 0, x1, y1)
	res := scaleTo(img, dr, sc)

	writer := bytes.NewBuffer(nil)

	if err = png.Encode(writer, res); err != nil {
		return nil, err
	}
	return writer.Bytes(), nil
}

// scaleTo scales the source image to fit within the specified rectangle using the provided scaling algorithm.
func scaleTo(src image.Image, rect image.Rectangle, scale draw.Scaler) image.Image {
	dst := image.NewRGBA(rect)
	scale.Scale(dst, rect, src, src.Bounds(), draw.Over, nil)
	return dst
}
