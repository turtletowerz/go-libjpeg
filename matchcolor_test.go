package jpeg

import (
	"image"
	"testing"
)

func TestMatchImage(t *testing.T) {
	for _, x := range colorMatches {
		a := image.NewNRGBA(image.Rect(0, 0, 1, 1))
		b := image.NewNRGBA(image.Rect(0, 0, 1, 1))
		a.Set(0, 0, x.a)
		b.Set(0, 0, x.b)

		if _, err := MatchImage(a, b, x.tolerance); (err == nil) != x.match {
			t.Errorf("MatchImage(a:%v b:%v, tolerance: %v) err:%v but want:%v", x.a, x.b, x.tolerance, err, x.match)
		}
	}
}

func TestMatchColor(t *testing.T) {
	for _, x := range colorMatches {
		if got := MatchColor(x.a, x.b, x.tolerance); x.match != got {
			t.Errorf("MatchColor(a:%v b:%v, tolerance: %v) got:%v but want:%v", x.a, x.b, x.tolerance, got, x.match)
		}
	}
}
