package geo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"../../hicup/hishape"
	"../../hicup/sphinx"
	"../shared"
)

// Search Result Count for geocoding
const SEARCH_COUNT = 1

type RComponent struct {
	Name string // contents name
	Type string // area level
}

type RGeometry struct {
	// request lon lat
	Location hishape.Point
	//
	ViewPort hishape.Box
}

// AreaMapper Result struct
type Mapper struct {
	s string
	p hishape.Point
	b hishape.Box
}
type Mappers []Mapper

func (a Mappers) Len() int           { return len(a) }
func (a Mappers) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Mappers) Less(i, j int) bool { return a[i].s > a[j].s }

type ResultPart struct {
	Components    []RComponent
	DetailAddress string
	Geometry      RGeometry
}

// Sorting Contents Len(), Swap(), Less()
func (a ResultPart) Len() int {
	return len(a.Components)
}
func (a ResultPart) Swap(i, j int) {
	a.Components[i], a.Components[j] = a.Components[j], a.Components[i]
}
func (a ResultPart) Less(i, j int) bool {
	return a.Components[i].Type < a.Components[j].Type
}

func (r *ResultPart) MakeAddress() {
	var lv2addr map[string]string
	lv2addr = make(map[string]string, len(r.Components))
	for _, v := range r.Components {
		lv2addr[v.Type] = v.Name
	}
	v, ok := lv2addr["Level1"] // state
	if ok {
		r.DetailAddress += v + " "
	}
	v, ok = lv2addr["Level2"] // city
	if ok {
		r.DetailAddress += v + " "
	}
	v, ok = lv2addr["Level3B"] // if townli
	if ok {
		v1, ok1 := lv2addr["Level3A"]
		if ok1 {
			r.DetailAddress += v1 + " "
		}
	}
	v, ok = lv2addr["Level4"] // street name
	if ok {
		r.DetailAddress += v + ", "
	}
	v, ok = lv2addr["Level7"] // building
	if ok {
		r.DetailAddress += v + " "
	}
	v, ok = lv2addr["Level3B"] // if townli
	if !ok {
		v1, ok1 := lv2addr["Level3A"]
		if ok1 {
			r.DetailAddress += "("
			r.DetailAddress += v1
			v2, ok2 := lv2addr["Level5"] // complex
			if ok2 {
				r.DetailAddress += ", " + v2
			}
			r.DetailAddress += ")"
		}
	}
}

type Result struct {
	Results []ResultPart
	Status  string
}

type SearchArea struct {
	areas []*Area
	// used for level7
	orderedfile map[string]int
	itemcounts  []int
}

func (s *SearchArea) Init() {
}

func (s SearchArea) Count() int {
	return len(s.areas)
}

func (s *SearchArea) Add(a *Area) {
	a.InitQuadTree()
	s.areas = append(s.areas, a)
}

// Call Once after all Add()
func (s *SearchArea) PrepareSearch() {
	s.itemcounts = make([]int, len(shared.Config.Level2File["Level7"])+1)
	s.orderedfile = make(map[string]int, len(shared.Config.Level2File["Level7"]))

	// orderd level4 !!
	// level2file is sorted in config.go
	for i, f := range shared.Config.Level2File["Level7"] {
		s.orderedfile[f] = i
		for _, a := range s.areas {
			if filepath.Base(a.Filename) == f {
				s.itemcounts[i+1] = s.itemcounts[i] + a.ItemCount()
			}
		}
	}
	log.Println("Items Count :", s.itemcounts)
}

func (s *SearchArea) CalID(nid, fid string) (int, int) {
	//itemcounts = itemcounts[:len?]
	// level71idx.bin : 0, level72idx.bin : 45xxxx
	// sid := baseid(fid)
	// id := nid - sid
	fidx, err1 := strconv.Atoi(fid)
	if err1 != nil {
		log.Println(err1)
	}
	//sid := s.itemcounts[s.orderedfile[fid]]
	sid := s.itemcounts[fidx]
	id, err2 := strconv.Atoi(nid)
	//log.Println("first id :", id)
	if err2 != nil {
		log.Println(err2)
	}
	id = id - sid
	//log.Println("last id :", id)
	return id, fidx
}

func (s *SearchArea) ReadAddress(n int, a *Area) (*Address, error) {
	value, err := a.ReadAttribute(n)
	if err != nil {
		log.Println("ReadAttribute Error :", n, " ", a.Filename)
	}
	//log.Println("level mapping = ", Config.AreaLevel[a.File])
	//log.Println("Name = ", name)
	log.Println("value :", value)
	area_address := NewAddress()
	area_address.CSVIn(value)
	return area_address, nil
}

/*
type ChanResult struct {
	v string
	p hishape.Point
	b hishape.Box
}
*/

func (s *SearchArea) AreaMapper(v *Area, p hishape.Point, c1 chan Mapper) {
	candidate := v.Query(p)
	searched := false
	for _, i := range candidate {
		target := hishape.NewPolygon()
		pts, err := v.ReadPoints(i)
		if err != nil {
			log.Println("Error ReadPoints().. It's not Possible : ", err)
			continue
		}
		target.AddInt32(pts)
		if target.PointInPolygon(p) {
			value, err := v.ReadAttribute(i)
			if err != nil {
				log.Println("ReadAttribute Error :", i, " ", v.Filename)
			}
			log.Println("value :", value)
			c1 <- Mapper{value, p, target.Box}
			searched = true
			break
		}
	}
	// no candidate
	if !searched {
		c1 <- Mapper{"", hishape.Point{}, hishape.Box{}}
	}
}

func (s *SearchArea) SearchAddress(p hishape.Point) []byte {
	rt := Result{}
	rt_part := ResultPart{}
	mapper_list := Mappers{}
	count_chan := len(s.areas)
	//log.Println("chan count =", count_chan)
	buf_c1 := make(chan Mapper, count_chan)
	for _, v := range s.areas {
		go s.AreaMapper(v, p, buf_c1)
	}
	//for range s.areas {
	for i := 0; i < count_chan; i++ {
		t1 := <-buf_c1
		if t1.s == "" {
			continue
		}
		area_address := Address{}
		area_address.CSVIn(t1.s)
		for _, lv := range shared.Config.LevelList {
			name := area_address.PriotityAddress(lv.Priority)
			if name == "" {
				continue
			}
			rt_part.Components = append(rt_part.Components, RComponent{name, lv.Level})
			mapper_list = append(mapper_list, Mapper{lv.Level, t1.p, t1.b})
		}
	}
	if len(mapper_list) > 0 {
		sort.Sort(mapper_list)
		sort.Sort(rt_part)
		// Delete Duplicate Items
		var prev_comp RComponent
		var tmp_rcomps []RComponent
		var tmp_mapper_list Mappers
		for i, v := range rt_part.Components {
			if v.Type == prev_comp.Type {
				continue
			}
			// Change Level3B and Level3A to Level3
			if v.Type == "Level3B" || v.Type == "Level3A" {
				v.Type = "Level3"
			}
			tmp_rcomps = append(tmp_rcomps, v)
			tmp_mapper_list = append(tmp_mapper_list, mapper_list[i])
			prev_comp = v
		}
		// copy object !
		rt_part.Components = tmp_rcomps
		mapper_list = tmp_mapper_list
		//log.Println("sort 0:", geo_list[0].s)
		//log.Println("sort 1:", geo_list[1].s)
		rt_part.Geometry = RGeometry{mapper_list[0].p, mapper_list[0].b}
	}
	// Make Detail Address ?
	rt_part.MakeAddress()
	//
	rt.Results = append(rt.Results, rt_part)
	// @todo
	if rt.Status != "Error" {
		rt.Status = "OK"
	}
	//rt_json, err := json.Marshal(rt)
	rt_json, err := json.MarshalIndent(rt, "", "  ")
	if err != nil {
		log.Println("error:", err)
	}
	return rt_json
}

var gstate []string = []string{"특별시", "광역시", "도", "특별자치시"}
var gcity []string = []string{"시", "군", "구"}
var gtown []string = []string{"읍", "면", "동"}
var gtownli []string = []string{"리"}

func (s *SearchArea) WhereClause(token []string) string {
	var result bytes.Buffer
	if len(token) == 0 {
		return result.String()
	}
	// special case
	var city_searched []string
	for _, v := range token {
		search_state := false
		// check state
		for _, s := range gstate {
			if strings.HasSuffix(v, s) {
				result.WriteString(fmt.Sprintf("|state='%s'", v))
				search_state = true
				break
			}
		}
		if search_state == true {
			continue
		}
		// check city .. need all search
		for _, s := range gcity {
			if strings.HasSuffix(v, s) {
				city_searched = append(city_searched, v)
				search_state = true
			}
		}
		if search_state == true {
			continue
		}
		// check town
		for _, s := range gtown {
			if strings.HasSuffix(v, s) {
				result.WriteString(fmt.Sprintf("|town='%s'", v))
				search_state = true
				break
			}
		}
		if search_state == true {
			continue
		}
		// check townli
		for _, s := range gtownli {
			if strings.HasSuffix(v, s) {
				result.WriteString(fmt.Sprintf("|townli='%s'", v))
				search_state = true
				break
			}
		}
		if search_state == true {
			continue
		}
	}
	if len(city_searched) == 1 {
		result.WriteString(fmt.Sprintf("|city='%s'", city_searched[0]))
	} else if len(city_searched) == 2 {
		result.WriteString(fmt.Sprintf("|city='%s %s'", city_searched[0],
			city_searched[1]))
	}
	//result = strings.TrimSpace(result)
	tmp := result.String()
	tmp = strings.Replace(tmp, "|", " and ", -1)
	return tmp
}

func (s *SearchArea) DoTokenJob(ss string) (string, string) {
	token := strings.Split(ss, ",")
	last_token := token[len(token)-1]
	// split by space
	quote_last_token := strings.Split(last_token, " ")
	// append a quote
	for i, v := range quote_last_token {
		quote_last_token[i] = fmt.Sprintf("\"%s\"", v)
	}
	// join token
	last_token = strings.Join(quote_last_token, " ")
	//
	other_token := token[0 : len(token)-1]
	restrict := s.WhereClause(other_token)
	log.Println("where : ", restrict)
	log.Println("Last token :", last_token)
	log.Println("other token :", other_token)
	return last_token, restrict
}

func (s *SearchArea) SearchPoint(ss string) []byte {
	rt := Result{}
	// Make a token string for search
	last_token, restrict := s.DoTokenJob(ss)
	// Text Search
	searchResults, err := sphinx.SearchX(last_token, restrict)
	if err != nil {
		log.Println(err)
		rt.Status = "Error"
	}
	//
	count := 0
	for _, id_obj := range searchResults {
		split_obj := strings.Split(id_obj, " ")
		log.Println("id = ", split_obj[0])
		id, fidx := s.CalID(split_obj[0], split_obj[1])
		// Search More Detail
		rt_part := ResultPart{}
		for _, v := range s.areas {
			// @todo How to treat another level, and level2file
			if filepath.Base(v.Filename) == shared.Config.Level2File["Level7"][fidx] {
				addr, err := s.ReadAddress(id, v)
				if err != nil {
					log.Println("Error ReadAddress:", err)
				}
				for _, lv := range shared.Config.LevelList {
					l2a := addr.LevelAddress(lv.Level)
					if l2a != "" {
						rt_part.Components = append(rt_part.Components,
							RComponent{l2a, lv.Level})
					}
				}
				target := hishape.NewPolygon()
				pts, err := v.ReadPoints(id)
				// test code
				//for _, k := range pts {
				//	log.Printf("%f,%f\n", float64(k.Y)/hishape.LLCONV,
				//	  float64(k.X)/hishape.LLCONV)
				//}
				if err != nil {
					log.Println("ReadPoint Error Not Possible: ", err)
					continue
				}
				target.AddInt32(pts)
				rt_part.Geometry = RGeometry{target.Center3(), target.Box}
			}
		}
		// Sort RComponent
		if rt_part.Len() != 0 {
			sort.Sort(rt_part)
		} else {
			// Getting a information error ?
			count += 1
			if count == SEARCH_COUNT {
				break
			}
			continue
		}
		// Detail address
		rt_part.MakeAddress()
		//
		rt.Results = append(rt.Results, rt_part)
		// Adding Only 3 items
		count += 1
		if count == SEARCH_COUNT {
			break
		}
	}
	// @todo
	if rt.Status != "Error" {
		rt.Status = "OK"
	}
	//rt_json, err := json.Marshal(rt)
	rt_json, err := json.MarshalIndent(rt, "", "  ")
	if err != nil {
		log.Println("error:", err)
	}
	return rt_json
}
