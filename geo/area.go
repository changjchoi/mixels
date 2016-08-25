package geo

import (
	"log"

	"../../hicup/hishape"
	"../shared"
	//"github.com/volkerp/goquadtree/quadtree"
)

type BoxInterface struct {
	Box   hishape.BoxInt32
	Index int // polygon index
}

// Implementation of the QuadTree package
func (b BoxInterface) BoundingBox() hishape.BoxInt32 {
	return b.Box
}

type Area struct {
	shared.RegionFile
	QT hishape.QuadTree
}

// Make QuadTree
func (a *Area) InitQuadTree() {
	g := &a.Header.Box
	qt_world := hishape.BoxInt32{g.MinX, g.MinY, g.MaxX, g.MaxY}
	a.QT = hishape.NewQuadTree(qt_world)
	// Loading Boxes
	a.LoadBoxs()
	// Iterate all Polygon
	for i, p := range a.Bounds {
		obj := BoxInterface{hishape.BoxInt32{p.MinX, p.MinY, p.MaxX, p.MaxY}, i}
		a.QT.Add(obj)
	}
	// Release Boxes
	a.ReleseBoxs()
}

// Convert User Point from float64 to int32
// Search Boxes
func (a *Area) Query(p hishape.Point) []int {
	p32 := p.Int32()
	delta := int32(100)
	bound := hishape.BoxInt32{
		p32.X - delta, p32.Y - delta, p32.X + delta, p32.Y + delta,
	}
	val := a.QT.Query(bound)
	ret := make([]int, len(val))
	// There are no way to convert from []interface{} to []SomethingType
	log.Println("QuadTree Searched Count = ", len(val))
	for i, v := range val {
		ret[i] = v.(BoxInterface).Index
	}
	return ret
}
