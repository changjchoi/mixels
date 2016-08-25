package geo

import (
	"bytes"
	"encoding/csv"
	"log"
	"strings"
)

type Address struct {
	// address set
	level1_name    string
	level2_name    string
	level3a_name   string
	level3b_name   string
	street_name    string // road_name + buid_conv_num + build_sub_num
	build_name     string
	build_name_dtl string
	jibun_name     string // jibun_main_num + jibun_sub_num + is_mountain
	// address set end
}

func NewAddress() *Address {
	return &Address{}
}

func (o Address) LevelAddress(lv string) string {
	switch lv {
	case "Level1":
		return o.State()
	case "Level2":
		return o.City()
	case "Level3A":
		return o.Town()
	case "Level3B":
		return o.TownLi()
	case "Level4": // road name and number
		return o.StreetName()
	case "Level5": // Complex level or complex name ?
		return o.BuildingName()
	case "Level6": // a lot number..jibun
		return o.Jibun()
	case "Level7": // a building name
		return o.BuildingNameSub()
	default:
		return "No Value"
	}
}

func (o Address) PriotityAddress(p int) string {
	if p < 5 {
		if p < 3 {
			if p == 1 {
				return o.State()
			} else if p == 2 {
				return o.City()
			}
		} else {
			if p == 3 {
				return o.Town()
			} else if p == 4 {
				return o.TownLi()
			}
		}
	} else {
		if p < 7 {
			if p == 5 {
				return o.StreetName()
			} else if p == 6 {
				return o.BuildingName()
			}
		} else {
			if p == 7 {
				return o.Jibun()
			} else if p == 8 {
				return o.BuildingNameSub()
			}
		}
	}
	return "No Value"
}

// Test All print
func (o Address) Print() {
	log.Println(
		o.level1_name,
		o.level2_name,
		o.level3a_name,
		o.level3b_name,
		o.street_name,
		o.build_name,
		o.build_name_dtl,
		o.jibun_name,
	)
}

// Used by dump_address
func (o *Address) CSVIn(in string) {
	r := csv.NewReader(strings.NewReader(in))
	r.Comma = '|'
	r.Comment = '#'
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal("csv read error:", err)
	}
	for _, record := range records {
		/*
			if len(record) != FIELD_COUNT {
				log.Fatal("Not match Field Count :", len(record), in, record)
			}
		*/
		o.level1_name = record[0]
		o.level2_name = record[1]
		o.level3a_name = record[2]
		o.level3b_name = record[3]
		o.street_name = record[4]
		o.build_name = record[5]
		o.build_name_dtl = record[6]
		o.jibun_name = record[7]
	}
}

func (o Address) CSVOut() string {
	record := []string{
		o.State(),
		o.City(),
		o.Town(),
		o.TownLi(),
		o.StreetName(),
		o.BuildingName(),
		o.BuildingNameSub(),
		o.Jibun(),
	}
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	w.Comma = '|'
	// using '\n' instead '\r\n'
	w.UseCRLF = false
	if err := w.Write(record); err != nil {
		log.Println("Error write csv")
	}
	w.Flush()
	// Cut '\n' byte
	return strings.TrimRight(buf.String(), "\n")
}

func (o Address) State() string {
	return o.level1_name
}

func (o *Address) SetState(i string) {
	o.level1_name = i
}

func (o Address) City() string {
	return o.level2_name
}

func (o *Address) SetCity(i string) {
	o.level2_name = i
}

func (o Address) Town() string {
	return o.level3a_name
}

func (o *Address) SetTown(i string) {
	o.level3a_name = i
}

func (o Address) TownLi() string {
	return o.level3b_name
}

func (o *Address) SetTownLi(i string) {
	o.level3b_name = i
}

func (o Address) StreetName() string {
	return o.street_name
}

func (o *Address) SetStreeName(i string) {
	o.street_name = i
}

// @todo add 번지 or not ?
func (o Address) Jibun() string {
	return o.jibun_name
}

func (o *Address) SetJibun(i string) {
	o.jibun_name = i
}

func (o Address) BuildingName() string {
	return o.build_name
}

func (o *Address) SetBuildingName(i string) {
	o.build_name = i
}

func (o Address) BuildingNameSub() string {
	return o.build_name_dtl
}

func (o *Address) SetBuildingNameSub(i string) {
	o.build_name_dtl = i
}
