package conv

import (
	//"bytes"
	//"encoding/binary"
	"log"
	//"os"
	"path/filepath"
	"sort"
	"strings"
	"unsafe"

	"../../hicup/hishape"
	. "../shared"
	"github.com/jonas-p/go-shp"
	"github.com/pebbe/go-proj-4/proj"
)

const POINTINT32_SIZE = int32(unsafe.Sizeof(hishape.PointInt32{}))

// Projection parameters
const WGS84 string = "+proj=latlong +ellps=WGS84 +datum=WGS84 +no_defs"

// Reference : http://www.osgeo.kr/17
const UTM_K1 string = "+proj=tmerc +lat_0=38 +lon_0=127.5 +k=0.9996 " +
	"+x_0=1000000 +y_0=2000000 +ellps=bessel +units=m +no_defs " +
	"+towgs84=-115.80,474.99,674.11,1.16,-2.31,-1.63,6.43"

// Reference : http://www.slideshare.net/jangbi882/proj4-32605736
const UTM_K2 string = "+proj=tmerc +lat_0=38 +lon_0=127.5 +k=0.9996 " +
	"+x_0=1000000 +y_0=2000000 +ellps=bessel +units=m +no_defs " +
	"+towgs84=-145.907,505.034,685.756,-1.162,2.347,1.592,6.342"

// Reference : http://www.biz-gis.com/
// http://spatialreference.org/ref/sr-org/7901/
// PCS_ITRF2000_TM : ok
const UTM_K3 string = "+proj=tmerc +lat_0=38 +lon_0=127.5 +k=0.9996 " +
	"+x_0=1000000 +y_0=2000000 +ellps=GRS80 +units=m +no_defs"
const GOOGLE string = "+proj=merc +a=6378137 +b=6378137 +lat_ts=0.0 " +
	"+lon_0=0.0 +x_0=0.0 +y_0=0 +k=1.0 +units=m +nadgrids=@null +no_defs"

type ConvertShapeFile struct {
	InputFile         string
	OutputFile        string
	isdirconv         bool
	from              *proj.Proj
	to                *proj.Proj
	shape             *shp.Reader
	regioncode        *RegionCode
	curField          []shp.Field
	curShape          *shp.Polygon
	curIndex          int
	prevAttributeBase int32
	prevPointBase     int32
	target            RegionFile
	projection        string
	//regionCode        RegionCode
}

func NewConvertShapeFile(i string, o string) *ConvertShapeFile {
	return &ConvertShapeFile{InputFile: i, OutputFile: o}
}

// set true : directory convertion
// set false : file convertion
func (o *ConvertShapeFile) SetDirConv(t bool) {
	o.isdirconv = t
}

func (o *ConvertShapeFile) Open(projection string) {
	var err error
	// set projection type
	o.projection = projection
	if projection == "UTM_K3" {
		o.from, err = proj.NewProj(UTM_K3)
		//defer g_from.Close()
		if err != nil {
			log.Fatal(err)
		}
	} else if projection == "WGS84" {
		o.from, err = proj.NewProj(WGS84)
		//defer g_from.Close()
		if err != nil {
			log.Fatal(err)
		}
	}
	// output coordination is alway WGS84
	o.to, err = proj.NewProj(WGS84)
	//defer g_to.Close()
	if err != nil {
		log.Fatal(err)
	}
	// Open a shapefile for reading
	if len(o.InputFile) != 0 {
		o.shape, err = shp.Open(o.InputFile)
		if err != nil {
			o.shape = nil
		}
	}
	// Loading region code data
	o.regioncode = NewRegionCode()
	o.regioncode.Open()
	// Default value in target Header Box
	o.target.Header.Box = *hishape.NewBoxInt32()
}

func (o *ConvertShapeFile) Reopen(filename string) {
	var err error
	if o.shape != nil {
		o.shape.Close()
	}
	o.shape = nil
	o.InputFile = filename
	// open a shapefile for reading
	o.shape, err = shp.Open(o.InputFile)
	if err != nil {
		log.Fatal(err)
	}
}

func (o *ConvertShapeFile) Iterate() {
	if o.shape == nil {
		log.Fatal("ShapeFile open Error")
	}
	// fields from the attribute table (DBF)
	o.curField = o.shape.Fields()
	// loop through all features in the shapefile
	for o.shape.Next() {
		// Shape() return : shape index, shape interface
		var shape_obj shp.Shape
		o.curIndex, shape_obj = o.shape.Shape()
		// @todo do a job when only polygon
		//log.Println("shape index = ", n)
		switch shape_obj.(type) {
		case *shp.Polygon:
			o.curShape = shape_obj.(*shp.Polygon)
			o.visitShape()
		default:
			log.Println("Unknown Format?")
		}
	}
}

func (o *ConvertShapeFile) visitDBF() int32 {
	attri := Attribute{}
	af := NewAddressField(o.regioncode)
	// Process file name
	var key string
	if o.isdirconv {
		key = filepath.Base(filepath.Dir(o.InputFile)) + "/*"
	} else {
		key = filepath.Base(filepath.Dir(o.InputFile)) + "/" +
			filepath.Base(o.InputFile)
	}
	for k, f := range o.curField {
		val := o.shape.ReadAttribute(o.curIndex, k)
		// Trim string?
		field_name := f.String()
		n := strings.IndexByte(field_name, 0)
		field_name = field_name[:n]
		// Get Mapping List : a file to field naame
		param, ok := Config.Mapping[key]
		if ok {
			//if strings.Compare(field_name, param.FieldType) == 0 {
			// Make a attribute
			//}
			if fn, good := param.FieldType[field_name]; !good {
				//log.Println("a file : ", param.File)
				//log.Fatal("There're no Field Name? :", field_name, " : ", key)
			} else {
				af.AddField(fn, val)
			}
		} else {
			log.Fatal("There're no Field Name? :", field_name, " : ", key)
		}
	}
	if af.FieldCount() == FIELD_COUNT || af.FieldCount() == 1 {
		attri.Value = append(attri.Value, af.CSVAddress()...)
		attri.ValueSize += int32(af.CSVAddressSize())
	} else {
		log.Fatal("Something Wrong :", key)
	}
	o.target.Attributes = append(o.target.Attributes, attri)
	return attri.ValueSize
}

func (o *ConvertShapeFile) visitShape() {
	// Make All Point slice
	var idx0, idx1 int32
	parts_len := len(o.curShape.Parts)
	for i := 0; i < parts_len; i += 1 {
		idx0 = o.curShape.Parts[i]
		if i < parts_len-1 {
			idx1 = o.curShape.Parts[i+1]
		} else {
			idx1 = o.curShape.NumPoints
		}
		bound := *hishape.NewBoxInt32()
		points := make([]hishape.PointInt32, idx1-idx0)
		for j := idx0; j < idx1; j += 1 {
			trans_x, trans_y := o.transform(o.curShape.Points[j].X,
				o.curShape.Points[j].Y)
			points[j-idx0] = hishape.PointInt32{trans_x, trans_y}
			bound.ExtendWithPoint(points[j-idx0])
		}
		attri_size := o.visitDBF() + 4 + 4
		point_size := (idx1 - idx0) * POINTINT32_SIZE
		ar := AreaRecord{PointBase: o.prevPointBase, PointSize: point_size,
			AttributeBase: o.prevAttributeBase, AttributeSize: attri_size}
		// Save prevAddress
		o.prevPointBase += point_size
		o.prevAttributeBase += attri_size
		// Append Data
		o.target.Header.Box.Extend(bound)
		o.target.CPoints = append(o.target.CPoints, points)
		o.target.Bounds = append(o.target.Bounds, bound)
		o.target.AreaRecords = append(o.target.AreaRecords, ar)
		// Everytime Do a job when an Item Count SPLITCOUNT
		if len(o.target.Attributes) == SPLITCOUNT {
			//
			o.modifyAttributeAddr()
			//
			o.target.SavePointAttri(o.OutputFile)
			// clear base
			o.prevPointBase, o.prevAttributeBase = 0, 0
			// clear contents
			o.target.CPoints = o.target.CPoints[:0]
			o.target.Attributes = o.target.Attributes[:0]
		}
	}
}

// @todo sorting and base address
func (o *ConvertShapeFile) modifyAttributeAddr() {
	//log.Println("prev Point base :", o.prevPointBase)
	n := len(o.target.AreaRecords)
	slot := (n - 1) / SPLITCOUNT
	start := slot * SPLITCOUNT
	end := start + SPLITCOUNT
	if n-start != SPLITCOUNT {
		end = n
	}
	// Sorting by Attribute Value
	sort.Sort(o.target)
	// Can it change target AreaRecord ?
	var prev_point_base int32
	prev_attri_base := o.prevPointBase
	for i, j := start, 0; i < end; i, j = i+1, j+1 {
		// bound don't need
		// point don't need
		// only arearecords ?
		o.target.AreaRecords[i].PointBase = prev_point_base
		o.target.AreaRecords[i].AttributeBase = prev_attri_base

		// The size don't be changed
		prev_point_base += o.target.AreaRecords[i].PointSize
		prev_attri_base += o.target.AreaRecords[i].AttributeSize
	}
}

func (o *ConvertShapeFile) transform(x, y float64) (ox, oy int32) {
	if o.projection == "WGS84" {
		//log.Printf("X: %v, Y: %v\n", x, y)
		return int32(x * hishape.LLCONV), int32(y * hishape.LLCONV)
	}
	cx, cy, err := proj.Transform2(o.from, o.to, x, y)
	if err != nil {
		log.Fatal(err)
	}
	return int32(proj.RadToDeg(cx) * hishape.LLCONV),
		int32(proj.RadToDeg(cy) * hishape.LLCONV)
}

func (o *ConvertShapeFile) Save() {
	// Modify atrribute base address
	o.modifyAttributeAddr()
	//
	if err := o.target.SaveIndex(o.OutputFile); err != nil {
		log.Fatal("Save Index Error :", err)
	}
	ac := len(o.target.Attributes)
	if ac < SPLITCOUNT && ac != 0 {
		if err := o.target.SavePointAttri(o.OutputFile); err != nil {
			log.Fatal("Save Points and Attributes Error :", err)
		}
	}
}

func (o *ConvertShapeFile) Close() {
	if o.from != nil {
		o.from.Close()
	}
	o.from = nil
	if o.to != nil {
		o.to.Close()
	}
	o.to = nil
	if o.shape != nil {
		o.shape.Close()
	}
	o.shape = nil
}
