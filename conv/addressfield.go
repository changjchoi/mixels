package conv

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"strings"

	"github.com/suapapa/go_hangul/encoding/cp949"
)

const FIELD_COUNT int = 12

type AddressField struct {
	RegionCode  *RegionCode
	field_count int
	address     string
	csv_address string
	// address set
	level1_name  string
	level2_name  string
	level3a_name string
	level3b_name string
	street_name  string // road_name + buid_conv_num + build_sub_num
	jibun_name   string // jibun_main_num + jibun_sub_num + is_mountain
	// address set end
	// field list start
	level1_code    string
	level2_code    string
	level3a_code   string
	level3b_code   string
	roadname_code  string
	build_main_num string
	build_sub_num  string
	build_name     string
	build_name_dtl string
	jibun_main_num string
	jibun_sub_num  string
	is_mountain    string
	floor          string
	// field list end
}

func NewAddressField(r *RegionCode) *AddressField {
	return &AddressField{RegionCode: r, build_main_num: "0", jibun_main_num: "0"}
}

func (o *AddressField) AddField(field string, value string) {
	switch field {
	case "level1_name":
		o.field_count++
		o.level1_name = value
	case "level2_name":
		o.field_count++
		o.level2_name = value
	case "level3a_name":
		o.field_count++
		o.level3a_name = value
	case "level3b_name":
		o.field_count++
		o.level3b_name = value
	case "level1_code":
		o.field_count++
		// level1 code setting here
		o.level1_code = value
	case "level2_code":
		o.field_count++
		o.level2_code = value
		// level1 code setting here
		o.level1_code = value[:2]
	case "level3a_code":
		o.field_count++
		o.level3a_code = value
	case "level3b_code":
		o.field_count++
		o.level3b_code = value
	case "build_main_num":
		o.field_count++
		o.build_main_num = value
	case "build_sub_num":
		o.field_count++
		o.build_sub_num = value
	case "build_name":
		o.field_count++
		got, err := cp949.From([]byte(value))
		if err != nil {
			log.Fatal("Error cp949 code :", err)
		}
		o.build_name = string(got)
		// Delete quote in the name
		o.build_name = strings.Replace(o.build_name, "\"", "", -1)
		// Delete tab
		o.build_name = strings.Replace(o.build_name, "\t", "", -1)
		// Delete carriage return
		o.build_name = strings.Replace(o.build_name, "\r", "", -1)
		// Delete line feed
		o.build_name = strings.Replace(o.build_name, "\n", "", -1)
		//log.Println(o.build_name)
	case "build_name_dtl":
		o.field_count++
		got, err := cp949.From([]byte(value))
		if err != nil {
			log.Fatal("Error cp949 code :", err)
		}
		o.build_name_dtl = string(got)
		o.build_name_dtl = strings.Replace(o.build_name_dtl, "\"", "", -1)
		o.build_name_dtl = strings.Replace(o.build_name_dtl, "\t", "", -1)
		o.build_name_dtl = strings.Replace(o.build_name_dtl, "\r", "", -1)
		o.build_name_dtl = strings.Replace(o.build_name_dtl, "\n", "", -1)
		//log.Println(o.build_name_dtl)
	case "roadname_code":
		o.field_count++
		o.roadname_code = value
		o.roadname_code = strings.Replace(o.roadname_code, "\"", "", -1)
		o.roadname_code = strings.Replace(o.roadname_code, "\t", "", -1)
		o.roadname_code = strings.Replace(o.roadname_code, "\r", "", -1)
		o.roadname_code = strings.Replace(o.roadname_code, "\n", "", -1)
	case "jibun_main_num":
		o.field_count++
		o.jibun_main_num = value
	case "jibun_sub_num":
		o.field_count++
		o.jibun_sub_num = value
	case "is_mountain":
		o.field_count++
		o.is_mountain = value
	case "floor":
		o.field_count++
		o.floor = value
	default:
		log.Fatal("AddField Error item? :", field)
	}
}

func (o *AddressField) FieldCount() int {
	return o.field_count
}

// Address Output Test Code
func (o *AddressField) Address() []byte {
	if o.RegionCode == nil {
		log.Fatal("You must Set RegionCode instance")
	}
	extra := ""
	var key1 string
	if o.level2_code != "" {
		key1 = o.level2_code[:2]
	}
	key2 := o.level2_code
	key3a := key1 + o.level3a_code
	key3b := key2 + o.level3b_code
	frt_addr := o.RegionCode.Level1[key1]
	frt_addr += " " + o.RegionCode.Level2[key2]
	//log.Println("lvl 3b :", o.level3b_code)
	//if strings.Compare(o.level3b_code, "00") == 0 {
	if o.level3b_code == "00" {
		extra += o.RegionCode.Level3A[key3a]
	} else {
		frt_addr += " " + o.RegionCode.Level3B[key3b]
	}
	if o.build_name != "" {
		if extra == "" {
			extra += o.build_name
		} else {
			extra += ", " + o.build_name
		}
	}
	frt_addr += " " + o.RegionCode.RoadNameCode[o.roadname_code]
	frt_addr += " " + o.build_main_num
	//if strings.Compare(o.build_sub_num, "0") != 0 {
	if o.build_sub_num == "0" {
		frt_addr += "-" + o.build_sub_num
	}
	//addr += o.build_name_dtl
	if extra == "" {
		if o.build_name_dtl == "" {
			o.address = fmt.Sprintf("%s", frt_addr)
		} else {
			o.address = fmt.Sprintf("%s, %s", frt_addr, o.build_name_dtl)
		}
	} else {
		if o.build_name_dtl == "" {
			o.address = fmt.Sprintf("%s, (%s)", frt_addr, extra)
		} else {
			o.address = fmt.Sprintf("%s, %s(%s)", frt_addr, o.build_name_dtl, extra)
		}
	}
	//log.Println(o.address)
	return []byte(o.address)
}

func (o *AddressField) AddressSize() int {
	return len(o.address)
}

// Used by dump_address
func (o *AddressField) AddCSVField(in string) {
	r := csv.NewReader(strings.NewReader(in))
	r.Comma = '|'
	r.Comment = '#'
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal("csv read error:", err)
	}
	for _, record := range records {
		if len(record) != 8 {
			log.Fatal("Not match Field Count :", len(record), in, record)
		}
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

// Test All print
func (o AddressField) Print() {
	log.Println(
		o.level1_code,
		o.level2_code,
		o.level3a_code,
		o.level3b_code,
		o.level1_name,
		o.level2_name,
		o.level3a_name,
		o.level3b_name,
		o.roadname_code,
		o.build_main_num,
		o.build_sub_num,
		o.build_name,
		o.build_name_dtl,
		o.jibun_main_num,
		o.jibun_sub_num,
		o.is_mountain,
		o.floor,
	)
}

// Map Address Output Format
func (o *AddressField) CSVAddress() []byte {
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
	o.csv_address = strings.TrimRight(buf.String(), "\n")
	//o.csv_address = buf.String()
	return []byte(o.csv_address)
}

func (o *AddressField) CSVAddressSize() int {
	return len(o.csv_address)
}

func (o AddressField) State() string {
	//return strings.Split(o.RegionCode.Level2[o.level2_code], " ")[0]
	if o.level1_name == "" {
		if o.level1_code != "" {
			o.level1_name = o.RegionCode.Level1[o.level1_code]
		}
	}
	return o.level1_name
}

// case 세종시
// case 성남시 분당구, 성남시 수정구
func (o AddressField) City() string {
	/*
		city := strings.Split(o.RegionCode.Level2[o.level2_code], " ")
		if len(city) == 2 {
			return city[1]
		} else if len(city) == 3 {
			return city[1] + " " + city[2]
		}
	*/
	//log.Println("3 or 1 :", len(city))
	//log.Println("Addr :", o.RegionCode.Level2[o.level2_code])
	if o.level2_name == "" {
		if o.level2_code != "" {
			o.level2_name = o.RegionCode.Level2[o.level2_code]
		}
	}
	return o.level2_name
}

func (o AddressField) Town() string {
	if o.level3a_name == "" {
		var key bytes.Buffer
		key.WriteString(o.level2_code)
		key.WriteString(o.level3a_code)
		if key.Len() != 0 {
			o.level3a_name = o.RegionCode.Level3A[key.String()]
		}
	}
	return o.level3a_name
}

func (o AddressField) TownLi() string {
	if o.level3b_name == "" {
		var key bytes.Buffer
		key.WriteString(o.level2_code)
		key.WriteString(o.level3a_code)
		key.WriteString(o.level3b_code)
		if key.Len() != 0 {
			o.level3b_name = o.RegionCode.Level3B[key.String()]
		}
	}
	return o.level3b_name
}

func (o AddressField) StreetName() string {
	if o.street_name == "" {
		var buffer bytes.Buffer
		if o.roadname_code != "" {
			// road name
			v, ok := o.RegionCode.RoadNameCode[o.roadname_code]
			if !ok {
				// @todo Really Error Log ?
				//log.Println("roadname ?: ", o.roadname_code)
				o.roadname_code = ""
			}
			buffer.WriteString(v)
			// space
			buffer.WriteString(" ")
		}
		if o.build_main_num != "0" {
			// building main number
			buffer.WriteString(o.build_main_num)
			// building sub number
			if o.build_sub_num != "0" {
				buffer.WriteString("-")
				buffer.WriteString(o.build_sub_num)
			}
		}
		// jibun address !!
		if o.roadname_code == "" {
			o.street_name = ""
		} else {
			o.street_name = buffer.String()
		}
	}
	return o.street_name
}

/*
func (o AddressField) StreetNum() string {
	var buffer bytes.Buffer
	buffer.WriteString(o.build_main_num)
	if o.build_sub_num != "0" {
		buffer.WriteString("-")
		buffer.WriteString(o.build_sub_num)
	}
	return buffer.String()
}
*/

// @todo add 번지 or not ?
func (o AddressField) Jibun() string {
	if o.jibun_name == "" {
		var buffer bytes.Buffer
		if o.is_mountain == "1" {
			buffer.WriteString("산")
		}
		if o.jibun_main_num != "0" {
			buffer.WriteString(o.jibun_main_num)
			if o.jibun_sub_num != "0" {
				buffer.WriteString("-")
				buffer.WriteString(o.jibun_sub_num)
			}
			buffer.WriteString("번지")
		}
		o.jibun_name = buffer.String()
	}
	return o.jibun_name
}

func (o AddressField) BuildingName() string {
	return o.build_name
}

func (o AddressField) BuildingNameSub() string {
	return o.build_name_dtl
}
