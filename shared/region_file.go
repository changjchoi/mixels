package shared

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"sync"
	"unsafe"

	. "../../hicup/hishape"
)

var Lock sync.Mutex

const SPLITCOUNT int = 1000000

// Make a file within 4G
type AreaRecord struct {
	PointBase     int32
	PointSize     int32
	AttributeBase int32
	AttributeSize int32
}

type RegionHeader struct {
	ItemCount int32
	Padding   int32
	Box       BoxInt32
}

type Attribute struct {
	ValueSize int32
	Padding   int32
	Value     []byte
}

type RegionFile struct {
	Header      RegionHeader
	Bounds      []BoxInt32
	AreaRecords []AreaRecord
	Points      []PointInt32 // read only
	Attributes  []Attribute
	//
	CPoints   [][]PointInt32
	Filename  string
	saveCount int
	idxFile   *os.File
	dataFile  []*os.File
}

func NewRegionFile(f string) *RegionFile {
	return &RegionFile{Filename: f}
}

func (r RegionFile) ItemCount() int {
	return int(r.Header.ItemCount)
}

// sorting contents
func (s RegionFile) Len() int {
	return len(s.Attributes)
}
func (s RegionFile) Swap(i, j int) {
	s.Attributes[i], s.Attributes[j] = s.Attributes[j], s.Attributes[i]
	s.CPoints[i], s.CPoints[j] = s.CPoints[j], s.CPoints[i]
	// swap bound and arearecords ?
	k := s.saveCount * SPLITCOUNT
	s.AreaRecords[i+k], s.AreaRecords[j+k] =
		s.AreaRecords[j+k], s.AreaRecords[i+k]
	s.Bounds[i+k], s.Bounds[j+k] = s.Bounds[j+k], s.Bounds[i+k]
}

// []byte ? less ?
func (s RegionFile) Less(i, j int) bool {
	return bytes.Compare(s.Attributes[i].Value, s.Attributes[j].Value) == -1
}

// end sorting

func (r RegionFile) PrintHeader() {
	log.Println("ItemCount :", r.ItemCount())
	log.Println("Bound :", r.Header.Box)
	//log.Println("Attributes base =", r.Header.AttributeBase)
}

func (r RegionFile) PrintBounds(n int) {
	if len(r.Bounds) != 0 {
		log.Println("Bounds :", r.Bounds[n])
	}
}

func (r RegionFile) PrintAreaRecords(n int) {
	ar, err := r.ReadAreaRecord(n)
	if err != nil {
		log.Println("Error : ", err)
	} else {
		log.Println("AreaRecords :", ar)
	}
}

func (r *RegionFile) SaveAll(file string) error {
	// Header !
	//value_count := len(o.target.Bounds)
	//header_size := int64(unsafe.Sizeof(o.target.Header))
	//bounds_size := int64(value_count * int(unsafe.Sizeof(Box{})))
	//ars_size := int64(value_count * int(unsafe.Sizeof(AreaRecord{})))
	//points_size := int64(len(o.target.Points) * int(unsafe.Sizeof(Point{})))

	r.Header.ItemCount = int32(len(r.Bounds))
	//r.Header.PointBase = header_size + bounds_size + ars_size
	//r.Header.AttributeBase = o.target.Header.PointBase + points_size

	h, err := os.Create(file)
	if err != nil {
		log.Println("file create error :", err)
		return err
	}
	h.Seek(0, os.SEEK_SET)
	binary.Write(h, binary.LittleEndian, r.Header)
	binary.Write(h, binary.LittleEndian, r.Bounds)
	binary.Write(h, binary.LittleEndian, r.AreaRecords)
	// @todo need flatten
	for _, v := range r.CPoints {
		binary.Write(h, binary.LittleEndian, v)
	}
	//binary.Write(h, binary.LittleEndian, r.Attributes)
	for _, v := range r.Attributes {
		binary.Write(h, binary.LittleEndian, v.ValueSize)
		binary.Write(h, binary.LittleEndian, v.Padding)
		binary.Write(h, binary.LittleEndian, v.Value)
	}
	h.Close()
	return nil
}

func (r *RegionFile) SaveIndex(fn string) error {
	idxf := fn[0:len(fn)-4] + "idx.bin"
	r.Header.ItemCount = int32(len(r.Bounds))
	h, err := os.Create(idxf)
	if err != nil {
		log.Println("file create error :", err)
		return err
	}
	h.Seek(0, os.SEEK_SET)
	binary.Write(h, binary.LittleEndian, r.Header)
	binary.Write(h, binary.LittleEndian, r.Bounds)
	binary.Write(h, binary.LittleEndian, r.AreaRecords)
	h.Close()
	return nil
}

func (r *RegionFile) SavePointAttri(fn string) error {
	paf := fmt.Sprintf("%s%03d.bin", fn[0:len(fn)-4], r.saveCount)
	h, err := os.Create(paf)
	if err != nil {
		log.Println("file create error :", err)
		return err
	}
	h.Seek(0, os.SEEK_SET)
	// @todo flatten
	for _, v := range r.CPoints {
		binary.Write(h, binary.LittleEndian, v)
	}
	//binary.Write(h, binary.LittleEndian, r.Attributes)
	for _, v := range r.Attributes {
		binary.Write(h, binary.LittleEndian, v.ValueSize)
		binary.Write(h, binary.LittleEndian, v.Padding)
		binary.Write(h, binary.LittleEndian, v.Value)
	}
	h.Close()
	// Increase the Count of Save
	r.saveCount += 1
	return nil
}

func (r *RegionFile) Open() error {
	var err error
	//r.fileName = filename
	r.idxFile, err = os.Open(r.Filename)
	if err != nil {
		log.Println("Open Error : ", r.Filename)
		return err
	}
	r.loadHeader()
	// Point and Attribute file open
	count_file := (r.ItemCount() - 1) / SPLITCOUNT
	r.dataFile = make([]*os.File, count_file+1)
	sf := r.Filename[0 : len(r.Filename)-7]
	for i := 0; i < count_file+1; i++ {
		f := fmt.Sprintf("%s%03d.bin", sf, i)
		//log.Println(f)
		df, err := os.Open(f)
		if err != nil {
			log.Fatal("error :", err)
		}
		//r.dataFile = append(r.dataFile, df)
		r.dataFile[i] = df
	}
	return nil
}

func (r *RegionFile) Close() {
	if r.idxFile != nil {
		r.idxFile.Close()
	}
	for _, v := range r.dataFile {
		if v != nil {
			v.Close()
			v = nil
		}
	}
}

func (r *RegionFile) loadHeader() {
	r.idxFile.Seek(0, os.SEEK_SET)
	binary.Read(r.idxFile, binary.LittleEndian, &r.Header)
}

func (r *RegionFile) LoadBoxs() {
	header_size := int64(unsafe.Sizeof(RegionHeader{}))
	r.idxFile.Seek(header_size, os.SEEK_SET)
	r.Bounds = make([]BoxInt32, r.ItemCount())
	binary.Read(r.idxFile, binary.LittleEndian, &r.Bounds)
}

func (r *RegionFile) ReleseBoxs() {
	// @todo really released ??
	r.Bounds = r.Bounds[:0]
}

func (r *RegionFile) ReadAreaRecord(n int) (AreaRecord, error) {
	ar := AreaRecord{}
	if r.ItemCount() < n {
		log.Println("Out of range ", r.ItemCount(), "<", n)
		return ar, fmt.Errorf("Out of range")
	}
	header_size := int64(unsafe.Sizeof(RegionHeader{}))
	box_size := int64(unsafe.Sizeof(BoxInt32{}))
	area_size := int64(unsafe.Sizeof(AreaRecord{}))
	base_addr := header_size + box_size*int64(r.Header.ItemCount) +
		int64(n)*area_size
	_, err := r.idxFile.Seek(base_addr, os.SEEK_SET)
	if err != nil {
		log.Println("Seek Error :", err)
		return ar, fmt.Errorf("Seek Error")
	}
	binary.Read(r.idxFile, binary.LittleEndian, &ar)
	return ar, nil
}

func (r *RegionFile) ReadAttribute(n int) (string, error) {
	attr := Attribute{}
	Lock.Lock()
	ar, err := r.ReadAreaRecord(n)
	if err != nil {
		Lock.Unlock()
		return "", err
	}
	addr := int64(ar.AttributeBase)
	size := ar.AttributeSize
	dump_bytes := make([]byte, size)
	fd := r.dataFile[n/SPLITCOUNT]
	//log.Println("selected file :", fd.Name())
	//log.Println("baseaddr : ", addr)
	//log.Println("size: ", size)
	//log.Println("n: ", n)
	_, err = fd.Seek(addr, os.SEEK_SET)
	if err != nil {
		log.Println("Seek Error :", err)
		Lock.Unlock()
		return "", fmt.Errorf("Seek Error")
	}
	_, err = fd.Read(dump_bytes)
	if err != nil {
		log.Println("Read Error :", err)
		Lock.Unlock()
		return "", fmt.Errorf("Read Error")
	}
	buf := bytes.NewReader(dump_bytes)
	binary.Read(buf, binary.LittleEndian, &attr.ValueSize)
	binary.Read(buf, binary.LittleEndian, &attr.Padding)
	//log.Println("value size = ", attr.ValueSize)
	attr.Value = make([]byte, attr.ValueSize)
	binary.Read(buf, binary.LittleEndian, &attr.Value)
	Lock.Unlock()
	return string(attr.Value), nil
}

func (r *RegionFile) ReadPoints(n int) ([]PointInt32, error) {
	ar, err := r.ReadAreaRecord(n)
	if err != nil {
		return []PointInt32{}, err
	}
	addr := int64(ar.PointBase)
	count := int32(ar.PointSize) / int32(unsafe.Sizeof(PointInt32{}))
	point := make([]PointInt32, count)
	fd := r.dataFile[n/SPLITCOUNT]
	fd.Seek(addr, os.SEEK_SET)
	binary.Read(fd, binary.LittleEndian, point)
	return point, nil
}
