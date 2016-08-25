package conv_test

import (
	"testing"

	"../conv"
	"../shared"
)

func TestAddressField(t *testing.T) {
	if shared.Config.Init("../config.json") == false {
		t.Fatal("map config.json loading error")
	}

	//	`11680|101|00|3122008|306|0|개나리 래미안|106동|754|0|0|22`,
	//	`41281|123|00|3192088|12|0|고양경찰서||1008|0|0|1`,
	test_case := []map[string]string{
		{"level2_code": "11680",
			"level3a_code":   "101",
			"level3b_code":   "00",
			"roadname_code":  "3122008",
			"build_main_num": "306",
			"build_sub_num":  "0",
			"build_name":     "CP949Name",
			"build_name_dtl": "106Dong",
			"jibun_main_num": "754",
			"jibun_sub_num":  "0",
			"is_mountain":    "0",
			"floor":          "22"},
		{"level2_code": "41281",
			"level3a_code":   "123",
			"level3b_code":   "00",
			"roadname_code":  "3192088",
			"build_main_num": "12",
			"build_sub_num":  "0",
			"build_name":     "GoYang Police",
			"build_name_dtl": "",
			"jibun_main_num": "1008",
			"jibun_sub_num":  "0",
			"is_mountain":    "0",
			"floor":          "1"},
	}
	code := conv.NewRegionCode()
	if err := code.Open(); err != nil {
		t.Fatal("region code open fail")
	}
	for _, v := range test_case {
		building := conv.NewAddressField(code)
		for key, value := range v {
			building.AddField(key, value)
		}
		t.Log(building.State())
		t.Log(building.City())
		t.Log(building.Town())
		t.Log(building.TownLi())
		t.Log(building.StreetName())
		t.Log(building.Jibun())
	}
	nulltest := conv.NewAddressField(code)
	t.Log(string(nulltest.CSVAddress()))
	if nulltest.CSVAddressSize() != 7 {
		t.Error("The size must be 7")
	}

}
