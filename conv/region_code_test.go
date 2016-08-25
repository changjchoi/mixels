package conv_test

import (
	"log"
	_ "os"
	"strings"
	"testing"

	"../conv"
	"../shared"
)

func TestRegionCode(t *testing.T) {
	if shared.Config.Init("../config.json") == false {
		t.Fatal("map config.json loading error")
	}
	rc := conv.RegionCode{}
	rc.Open()
	test_value := map[string]string{
		"1000033": "용인서울고속도로",
		"2000003": "남부순환로",
		"2000006": "양재대로",
		"3122008": "역삼로",
		"3122010": "테헤란로",
		"2000010": "중앙대로",
		"2007002": "달구벌대로",
		"1000028": "경인고속도로",
		"3020001": "연명로",
	}
	for k, v := range test_value {
		// is exist key?
		tv, ok := rc.RoadNameCode[k]
		if !ok {
			t.Fail()
		}
		if strings.Compare(tv, v) != 0 {
			t.Fail()
		}
		t.Log("org:", v, "conv:", tv)
	}
	t.Log("seoul =", rc.Level1["11"])
}
