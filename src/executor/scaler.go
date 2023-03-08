package executor

import (
	"bytes"
	"golang.org/x/image/draw"
	"image"
	"image/png"
)

func Scale(buf []byte, x1 int, y1 int) ([]byte, error) {
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
	dr := image.Rect(0, 0, x1, y1)
	res := scaleTo(img, dr, sc)

	writer := bytes.NewBuffer(nil)

	if err = png.Encode(writer, res); err != nil {
		return nil, err
	}
	return writer.Bytes(), nil
}

func scaleTo(src image.Image, rect image.Rectangle, scale draw.Scaler) image.Image {
	dst := image.NewRGBA(rect)
	scale.Scale(dst, rect, src, src.Bounds(), draw.Over, nil)
	return dst
}
