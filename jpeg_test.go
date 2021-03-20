package jpeg

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/color"
	nativeJPEG "image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var naturalImageFiles = []string{
	"images/cosmos.jpg",
	"images/kinkaku.jpg",
}

var subsampledImageFiles = []string{
	"images/checkerboard_444.jpg",
	"images/checkerboard_440.jpg",
	"images/checkerboard_422.jpg",
	"images/checkerboard_420.jpg",
}

// WritePNG encodes and writes image into PNG file.
func WritePNG(img image.Image, name string) {
	f, err := os.Create(GetOutFilePath(name))
	if err != nil {
		panic(err)
	}
	b := bufio.NewWriter(f)
	defer func() {
		b.Flush()
		f.Close()
	}()

	if err := png.Encode(b, img); err != nil {
		panic(err)
	}
}

// GetOutFilePath returns the path of specified out file.
func GetOutFilePath(name string) string {
	p := filepath.Join("images", "out", name)
	os.MkdirAll(filepath.Dir(p), 0666)
	return p
}

func TestMain(m *testing.M) {
	result := m.Run()
	if SourceManagerMapLen() > 0 {
		fmt.Println("sourceManager leaked")
		result = 2
	}
	if DestinationManagerMapLen() > 0 {
		fmt.Println("destinationManager leaked")
		result = 2
	}
	os.Exit(result)
}

func delta2(u0, u1 uint32) int {
	d := int(u0) - int(u1)
	if d < 0 {
		return -d
	}
	return d
}

func BenchmarkDecode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, file := range naturalImageFiles {
			r, err := os.Open(file)
			if err != nil {
				b.Errorf("opening file: %w", err)
			}
			img, err := Decode(r, &DecoderOptions{})
			if img == nil {
				b.Error("Got nil")
			}
			if err != nil {
				b.Errorf("Got Error: %v", err)
			}
		}
	}
}

func BenchmarkDecodeIntoRGB(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, file := range naturalImageFiles {
			r, err := os.Open(file)
			if err != nil {
				b.Errorf("opening file: %w", err)
			}
			img, err := DecodeIntoRGB(r, &DecoderOptions{})
			if img == nil {
				b.Error("Got nil")
			}
			if err != nil {
				b.Errorf("Got Error: %v", err)
			}
		}
	}
}

func BenchmarkDecodeWithNativeJPEG(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, file := range naturalImageFiles {
			r, err := os.Open(file)
			if err != nil {
				b.Errorf("opening file: %w", err)
			}
			img, err := nativeJPEG.Decode(r)
			if img == nil {
				b.Error("Got nil")
			}
			if err != nil {
				b.Errorf("Got Error: %v", err)
			}
		}
	}
}

func TestDecode(t *testing.T) {
	for _, file := range naturalImageFiles {
		r, err := os.Open(file)
		if err != nil {
			t.Errorf("opening file: %w", err)
		}
		fmt.Printf(" - test: %s\n", file)

		img, err := Decode(r, &DecoderOptions{})
		if err != nil {
			t.Errorf("Got Error: %v", err)
		}

		WritePNG(img, fmt.Sprintf("TestDecode_%s.png", file))
	}
}

func TestDecodeScaled(t *testing.T) {
	for _, file := range naturalImageFiles {
		r, err := os.Open(file)
		if err != nil {
			t.Errorf("opening file: %w", err)
		}
		fmt.Printf(" - test: %s\n", file)

		img, err := Decode(r, &DecoderOptions{ScaleTarget: image.Rect(0, 0, 100, 100)})
		if err != nil {
			t.Errorf("Got Error: %v", err)
		}
		if got := img.Bounds().Dx(); got != 256 {
			t.Errorf("Wrong scaled width: %v, expect: 128 (=1024/8)", got)
		}
		if got := img.Bounds().Dy(); got != 192 {
			t.Errorf("Wrong scaled height: %v, expect: 192 (=768/8)", got)
		}

		WritePNG(img, fmt.Sprintf("TestDecodeScaled_%s.png", file))
	}
}

func TestDecodeIntoRGBA(t *testing.T) {
	if SupportRGBA() != true {
		t.Skipf("This build is not support DecodeIntoRGBA.")
		return
	}
	for _, file := range naturalImageFiles {
		r, err := os.Open(file)
		if err != nil {
			t.Errorf("opening file: %w", err)
		}
		fmt.Printf(" - test: %s\n", file)

		img, err := DecodeIntoRGBA(r, &DecoderOptions{})
		if err != nil {
			t.Errorf("Got Error: %v", err)
			continue
		}

		WritePNG(img, fmt.Sprintf("TestDecodeIntoRGBA_%s.png", file))
	}
}

func TestDecodeScaledIntoRGBA(t *testing.T) {
	if SupportRGBA() != true {
		t.Skipf("This build is not support DecodeIntoRGBA.")
		return
	}
	for _, file := range naturalImageFiles {
		r, err := os.Open(file)
		if err != nil {
			t.Errorf("opening file: %w", err)
		}
		fmt.Printf(" - test: %s\n", file)

		img, err := DecodeIntoRGBA(r, &DecoderOptions{ScaleTarget: image.Rect(0, 0, 100, 100)})
		if err != nil {
			t.Errorf("Got Error: %v", err)
			continue
		}
		if got := img.Bounds().Dx(); got != 256 {
			t.Errorf("Wrong scaled width: %v, expect: 128 (=1024/8)", got)
		}
		if got := img.Bounds().Dy(); got != 192 {
			t.Errorf("Wrong scaled height: %v, expect: 192 (=768/8)", got)
		}

		WritePNG(img, fmt.Sprintf("TestDecodeIntoRGBA_%s.png", file))
	}
}

func TestDecodeScaledIntoRGB(t *testing.T) {
	for _, file := range naturalImageFiles {
		r, err := os.Open(file)
		if err != nil {
			t.Errorf("opening file: %w", err)
		}
		fmt.Printf(" - test: %s\n", file)

		img, err := DecodeIntoRGB(r, &DecoderOptions{ScaleTarget: image.Rect(0, 0, 100, 100)})
		if err != nil {
			t.Errorf("Got Error: %v", err)
		}
		if got := img.Bounds().Dx(); got != 256 {
			t.Errorf("Wrong scaled width: %v, expect: 128 (=1024/8)", got)
		}
		if got := img.Bounds().Dy(); got != 192 {
			t.Errorf("Wrong scaled height: %v, expect: 192 (=768/8)", got)
		}

		WritePNG(img, fmt.Sprintf("TestDecodeIntoRGB_%s.png", file))
	}
}

func TestDecodeSubsampledImage(t *testing.T) {
	for _, file := range subsampledImageFiles {
		r, err := os.Open(file)
		if err != nil {
			t.Errorf("opening file: %w", err)
		}
		fmt.Printf(" - test: %s\n", file)

		img, err := Decode(r, &DecoderOptions{})
		if err != nil {
			t.Errorf("Got Error: %v", err)
		}

		WritePNG(img, fmt.Sprintf("TestDecodeSubsampledImage_%s.png", file))
	}
}

func TestDecodeAndEncode(t *testing.T) {
	for _, file := range naturalImageFiles {
		r, err := os.Open(file)
		if err != nil {
			t.Errorf("opening file: %w", err)
		}
		fmt.Printf(" - test: %s\n", file)

		img, err := Decode(r, &DecoderOptions{})
		if err != nil {
			t.Errorf("Decode returns error: %v", err)
		}

		// Create output file
		f, err := os.Create(GetOutFilePath(fmt.Sprintf("TestDecodeAndEncode_%s", file)))
		if err != nil {
			panic(err)
		}
		w := bufio.NewWriter(f)
		defer func() {
			w.Flush()
			f.Close()
		}()

		if err := Encode(w, img, &EncoderOptions{Quality: 90}); err != nil {
			t.Errorf("%s: Encode returns error: %v", file, err)
		}
	}
}

func TestDecodeAndEncodeSubsampledImages(t *testing.T) {
	for _, file := range subsampledImageFiles {
		r, err := os.Open(file)
		if err != nil {
			t.Errorf("opening file: %w", err)
		}
		fmt.Printf(" - test: %s\n", file)

		img, err := Decode(r, &DecoderOptions{})
		if err != nil {
			t.Errorf("Decode returns error: %v", err)
		}

		// Create output file
		f, err := os.Create(GetOutFilePath(fmt.Sprintf("TestDecodeAndEncodeSubsampledImages_%s", file)))
		if err != nil {
			panic(err)
		}
		w := bufio.NewWriter(f)
		defer func() {
			w.Flush()
			f.Close()
		}()

		if err := Encode(w, img, &EncoderOptions{Quality: 90}); err != nil {
			t.Errorf("Encode returns error: %v", err)
		}
	}
}

func TestEncodeGrayImage(t *testing.T) {
	w, h := 400, 200
	img := image.NewGray(image.Rect(0, 0, w, h))

	// make gradient
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			img.SetGray(x, y, color.Gray{uint8(float64(x*y) / float64(w*h) * 255)})
		}
	}

	// encode gray gradient
	f, err := os.Create(GetOutFilePath(fmt.Sprintf("TestEncodeGrayImage_%dx%d.jpg", w, h)))
	if err != nil {
		panic(err)
	}
	wr := bufio.NewWriter(f)
	defer func() {
		wr.Flush()
		f.Close()
	}()
	if err := Encode(wr, img, &EncoderOptions{Quality: 90}); err != nil {
		t.Errorf("Encode returns error: %v", err)
	}
	wr.Flush()

	// rewind to first
	f.Seek(0, 0)

	// decode file
	decoded, err := Decode(f, &DecoderOptions{})
	if err != nil {
		t.Errorf("Decode returns error: %v", err)
	}
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			r, g, b, _ := decoded.At(x, y).RGBA()
			ref := uint32(float64(x*y) / float64(w*h) * 255)
			if delta2((r>>8), ref) > 1 || delta2((g>>8), ref) > 1 || delta2((b>>8), ref) > 1 {
				t.Errorf("(%d, %d): got (%d, %d, %d) want %v", x, y, r, g, b, ref)
			}
		}
	}
}

func TestDecodeConfig(t *testing.T) {
	for _, file := range naturalImageFiles {
		r, err := os.Open(file)
		if err != nil {
			t.Errorf("opening file: %w", err)
		}
		fmt.Printf(" - test: %s\n", file)

		config, err := DecodeConfig(r)
		if err != nil {
			t.Errorf("Got error: %v", err)
		}

		if got := config.ColorModel; got != color.YCbCrModel {
			t.Errorf("got wrong ColorModel: %v, expect: color.YCbCrModel", got)
		}
		if got := config.Width; got != 1024 {
			t.Errorf("got wrong width: %d, expect: 1024", got)
		}
		if got := config.Height; got != 768 {
			t.Errorf("got wrong height: %d, expect: 768", got)
		}
	}
}

func TestNewYCbCrAlignedWithLandscape(t *testing.T) {
	got := NewYCbCrAligned(image.Rect(0, 0, 125, 25), image.YCbCrSubsampleRatio444)

	if len(got.Y) != 6912 {
		t.Errorf("wrong array size Y: %d, expect: 6912", len(got.Y))
	}
	if len(got.Cb) != 6912 {
		t.Errorf("wrong array size Cb: %d, expect: 6912", len(got.Cb))
	}
	if len(got.Cr) != 6912 {
		t.Errorf("wrong array size Cr: %d, expect: 6912", len(got.Cr))
	}
	if got.YStride != 144 {
		t.Errorf("got wrong YStride: %d, expect: 128", got.YStride)
	}
	if got.CStride != 144 {
		t.Errorf("got wrong CStride: %d, expect: 128", got.CStride)
	}
}

func TestNewYCbCrAlignedWithPortrait(t *testing.T) {
	got := NewYCbCrAligned(image.Rect(0, 0, 25, 125), image.YCbCrSubsampleRatio444)

	if len(got.Y) != 6912 {
		t.Errorf("wrong array size Y: %d, expect: 6912", len(got.Y))
	}
	if len(got.Cb) != 6912 {
		t.Errorf("wrong array size Cb: %d, expect: 6912", len(got.Cb))
	}
	if len(got.Cr) != 6912 {
		t.Errorf("wrong array size Cr: %d, expect: 6912", len(got.Cr))
	}
	if got.YStride != 48 {
		t.Errorf("got wrong YStride: %d, expect: 128", got.YStride)
	}
	if got.CStride != 48 {
		t.Errorf("got wrong CStride: %d, expect: 128", got.CStride)
	}
}

func TestDecodeFailsWithBlankFile(t *testing.T) {
	blank := bytes.NewBuffer(nil)
	_, err := Decode(blank, &DecoderOptions{})
	if err == nil {
		t.Errorf("got no error with blank file")
	}
}

func TestEncodeFailsWithEmptyImage(t *testing.T) {
	dummy := &image.YCbCr{}
	w := bytes.NewBuffer(nil)
	err := Encode(w, dummy, &EncoderOptions{})
	if err == nil {
		t.Errorf("got no error with empty image")
	}
}

func newRGBA() *image.RGBA {
	rgba := image.NewRGBA(image.Rect(0, 0, 4, 8))
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			rgba.SetRGBA(i, j, color.RGBA{255, 0, 0, 255})
		}
		for j := 4; j < 8; j++ {
			rgba.SetRGBA(i, j, color.RGBA{0, 0, 255, 255})
		}
	}
	return rgba
}

func TestEncodeRGBA(t *testing.T) {
	rgba := newRGBA()
	w := bytes.NewBuffer(nil)

	err := Encode(w, rgba, &EncoderOptions{
		Quality: 100,
	})
	if err != nil {
		t.Fatalf("failed to encode: %v", err)
	}

	decoded, err := Decode(w, &DecoderOptions{})
	if err != nil {
		t.Fatalf("failed to decode: %v", err)
	}

	diff, err := MatchImage(rgba, decoded, 1)
	if err != nil {
		t.Errorf("match image: %v", err)
		WritePNG(rgba, "TestEncodeRGBA.want.png")
		WritePNG(decoded, "TestEncodeRGBA.got.png")
		WritePNG(diff, "TestEncodeRGBA.diff.png")
	}
}

// See: https://github.com/pixiv/go-libjpeg/issues/36
func TestDecodeAndEncodeRGBADisableFancyUpsampling(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, 3000, 2000))

	w, err := ioutil.TempFile("", "jpeg_test_")
	if err != nil {
		t.Fatalf("failed to create a file: %v", err)
	}
	name := w.Name()
	defer os.Remove(w.Name())

	err = Encode(w, src, &EncoderOptions{Quality: 95})
	w.Close()
	if err != nil {
		t.Fatalf("faled to encode: %v", err)
	}

	r, err := os.Open(name)
	if err != nil {
		t.Fatalf("failed to open: %v", err)
	}
	defer r.Close()

	_, err = DecodeIntoRGBA(r, &DecoderOptions{
		DisableBlockSmoothing:  true,
		DisableFancyUpsampling: true,
	})
	if err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
}

// https://github.com/pixiv/go-libjpeg/issues/51
func TestIssue51(t *testing.T) {
	data := []byte("\xff\xd8\xff\xdb\x00C\x000000000000000" +
		"00000000000000000000" +
		"00000000000000000000" +
		"00000000000\xff\xc9\x00\v\b00\x000" +
		"\x01\x01\x14\x00\xff\xda\x00\b\x01\x010\x00?\x0000")

	Decode(bytes.NewReader(data), &DecoderOptions{})
}

// https://github.com/pixiv/go-libjpeg/issues/55
func TestIssue55(t *testing.T) {
	data := []byte("\xff\xd8\xff\xdb\x00C\x000000000000000" +
		"00000000000000000000" +
		"00000000000000000000" +
		"00000000000\xff\xc0\x00\x11\b\x00000" +
		"\x03R\"\x00G\x11\x00B\x11\x00\xff\xda\x00\f\x03R\x00G\x00B" +
		"\x00")

	img, err := Decode(bytes.NewReader(data), &DecoderOptions{})
	if err != nil {
		t.Log(err)
		return
	}

	var w bytes.Buffer
	err = Encode(&w, img, &EncoderOptions{})
	if err != nil {
		t.Errorf("encoding after decoding failed: %v", err)
	}
}
