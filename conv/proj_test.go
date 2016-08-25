package conv_test

import (
	"testing"

	"github.com/pebbe/go-proj-4/proj"
)

const WGS84 string = "+proj=latlong +ellps=WGS84 +datum=WGS84 +no_defs"

func TestProjection(t *testing.T) {

	from, err1 := proj.NewProj(WGS84)
	if err1 != nil {
		t.Fatal("init proj fail")
	}
	to, err2 := proj.NewProj(WGS84)
	if err2 != nil {
		t.Fatal("init proj fail")
	}
	x := 127.4
	y := 34.2
	rx, ry, err3 := proj.Transform2(from, to, x, y)
	if err3 != nil {
		t.Fatal("fail to transform")
	}
	if proj.RadToDeg(rx) == rx {
		t.Logf("x -> rx = %v %v", x, proj.RadToDeg(rx))
		t.Error("x != rx error")
	}
	if proj.RadToDeg(ry) == ry {
		t.Logf("y -> ry = %v %v", y, proj.RadToDeg(ry))
		t.Error("y != ry error")
	}
}
