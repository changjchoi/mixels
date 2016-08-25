package shared

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sort"
)

type MapConfig struct {
	ServicePort string
	InputPath   string
	OutputPath  string
	Mapping     map[string]struct {
		File       string
		FieldType  map[string]string
		Projection string
	}
	AreaLevel      map[string]string
	RegionCodeFile map[string]string
	LevelList      []struct {
		Level    string
		Priority int
	}
	TSVFile    []string
	Level2File map[string][]string
}

var Config MapConfig

func (s *MapConfig) Init(file string) bool {
	// Working Global ?
	log.SetFlags(log.Flags() | log.Lshortfile)

	f, err := ioutil.ReadFile(file)
	if err != nil {
		log.Println("file open error")
		return false
	}
	err = json.Unmarshal(f, s)
	if err != nil {
		log.Println("json parsing error", err)
		return false
	}
	//log.Println("AreaLevel size =", len(s.AreaLevel))
	s.Level2File = make(map[string][]string, len(s.AreaLevel))
	for key, value := range s.AreaLevel {
		// search key by value
		v, _ := s.Level2File[value]
		s.Level2File[value] = append(v, key)
	}
	// level2file sorting
	if !sort.StringsAreSorted(s.Level2File["Level7"]) {
		sort.Strings(s.Level2File["Level7"])
	}
	return true
}

func (s MapConfig) Print() {
	log.Printf("%#v", s)
}
