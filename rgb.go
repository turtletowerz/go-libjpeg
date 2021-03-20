// Package rgb provides RGB image which implements image.Image interface.
package jpeg

import (
	"image"
	"image/color"
)

// RGBModel is an RGB color model instance
var RGBModel = color.ModelFunc(rgbModel)

func rgbModel(c color.Color) color.Color {
	if _, ok := c.(ColorRGB); ok {
		return c
	}
	r, g, b, _ := c.RGBA()
	return ColorRGB{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8)}
}

type ColorRGB struct {
	R, G, B uint8
}

func (c ColorRGB) RGBA() (r, g, b, a uint32) {
	r = uint32(c.R)
	r |= r << 8
	g = uint32(c.G)
	g |= g << 8
	b = uint32(c.B)
	b |= b << 8
	a = uint32(0xFFFF)
	return
}

// RGB represents an RGBA model where the Alpha value is always 0
type RGB struct {
	// Pix holds the image's pixels, in R, G, B order. The pixel at
	// (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*3].
	Pix []uint8

	// Stride is the Pix stride (in bytes) between vertically adjacent pixels.
	Stride int

	// Rect is the image's bounds.
	Rect image.Rectangle
}

func (r *RGB) ColorModel() color.Model {
	return RGBModel
}

func (r *RGB) Bounds() image.Rectangle {
	return r.Rect
}

func (r *RGB) At(x, y int) color.Color {
	return r.RGBAt(x, y)
}

// RGBAt returns the color of the pixel at (x, y) as ColorRGB.
func (r *RGB) RGBAt(x, y int) ColorRGB {
	if !(image.Point{x, y}.In(r.Rect)) {
		return ColorRGB{}
	}
	i := r.PixOffset(x, y)
	s := r.Pix[i : i+3 : i+3] // Small cap improves performance, see https://golang.org/issue/27857
	return ColorRGB{s[0], s[1], s[2]}
}

// PixOffset returns the index of the first element of Pix that corresponds to the pixel at (x, y).
func (r *RGB) PixOffset(x, y int) int {
	return (y-r.Rect.Min.Y)*r.Stride + (x-r.Rect.Min.X)*3
}

// Set allows for the implementation of draw.Image.
func (p *RGB) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}

	i := p.PixOffset(x, y)
	c1 := RGBModel.Convert(c).(ColorRGB)
	s := p.Pix[i : i+3 : i+3] // Small cap improves performance, see https://golang.org/issue/27857
	s[0] = c1.R
	s[1] = c1.G
	s[2] = c1.B
}

// SetRGB is for a faster method of Setting pixels, avoiding interface conversions
func (p *RGB) SetRGB(x, y int, c ColorRGB) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}

	i := p.PixOffset(x, y)
	s := p.Pix[i : i+3 : i+3] // Small cap improves performance, see https://golang.org/issue/27857
	s[0] = c.R
	s[1] = c.G
	s[2] = c.B
}

func NewRGB(r image.Rectangle) *RGB {
	w, h := r.Dx(), r.Dy()
	return &RGB{
		Pix:    make([]uint8, 3*w*h), // TODO: will this be a problem if the bounds are negative or can overflow?
		Stride: 3 * w,
		Rect:   r,
	}
}

// Make sure Image implements image.Image.
// See https://golang.org/doc/effective_go.html#blank_implements.
var _ image.Image = new(RGB)
