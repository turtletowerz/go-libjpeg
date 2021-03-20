//
// The test implemented in this file are originally from the tests for the
// source of Go. The portions are:
//
// Copyright (c) 2012 The Go Authors. All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//    * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//    * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//    * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//

package jpeg

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"testing"
)

type imageTest struct {
	refFilename string
	filename    string
	tolerance   int
}

var imageTests = []imageTest{
	{"images/testdata/video-001.221212.png", "images/testdata/video-001.221212.jpeg", 8 << 8},
	// {"images/testdata/video-001.cmyk.png", "images/testdata/video-001.cmyk.jpeg", 8 << 8},
	{"images/testdata/video-001.png", "images/testdata/video-001.jpeg", 8 << 8},
	{"images/testdata/video-001.png", "images/testdata/video-001.progressive.jpeg", 8 << 8},
	{"images/testdata/video-001.png", "images/testdata/video-001.rgb.jpeg", 8 << 16},
	{"images/testdata/video-005.gray.png", "images/testdata/video-005.gray.jpeg", 8 << 8},
}

func withinTolerance(c0, c1 color.Color, tolerance int) bool {
	r0, g0, b0, a0 := c0.RGBA()
	r1, g1, b1, a1 := c1.RGBA()
	r := delta2(r0, r1)
	g := delta2(g0, g1)
	b := delta2(b0, b1)
	a := delta2(a0, a1)
	return r <= tolerance && g <= tolerance && b <= tolerance && a <= tolerance
}

func TestDecodeTestdataFromGoStdlib(t *testing.T) {
Loop:
	for _, it := range imageTests {
		io, err := os.Open(it.filename)
		if err != nil {
			t.Errorf("opening file 1: %w", err)
		}
		img, err := Decode(io, &DecoderOptions{})
		if err != nil {
			t.Errorf("%s: %v", it.filename, err)
			continue
		}
		if img == nil {
			t.Error("got nil")
			continue
		}

		WritePNG(img, fmt.Sprintf("TestDecode_testdata_%s.png", it.filename[len("testdata/"):]))

		io2, err := os.Open(it.refFilename)
		if err != nil {
			t.Errorf("opening file 1: %w", err)
		}
		ref, _, err := image.Decode(io2)
		if err != nil {
			t.Errorf("%s: %v", it.refFilename, err)
		}

		gb := img.Bounds()
		rb := ref.Bounds()
		if !gb.Eq(rb) {
			t.Errorf("%s: got bounds %v want %v", it.filename, gb, rb)
			continue
		}

		for y := rb.Min.Y; y < rb.Max.Y; y++ {
			for x := rb.Min.X; x < rb.Max.X; x++ {
				if !withinTolerance(ref.At(x, y), img.At(x, y), it.tolerance) {
					t.Errorf("%s: at (%d, %d): got %v want %v", it.filename, x, y, ref.At(x, y), img.At(x, y))
					continue Loop
				}
			}
		}
	}
}
