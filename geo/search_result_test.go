package geo_test

import (
	"sort"
	"testing"

	"../geo"
	//"../shared"
)

func TestRAddrSort(t *testing.T) {
	rp := geo.ResultPart{}
	ra3 := geo.RComponent{"Level3", "테스트3"}
	rp.Components = append(rp.Components, ra3)
	ra1 := geo.RComponent{"Level1", "테스트1"}
	rp.Components = append(rp.Components, ra1)
	ra4 := geo.RComponent{"Level4", "테스트4"}
	rp.Components = append(rp.Components, ra4)
	ra2 := geo.RComponent{"Level2", "테스트2"}
	rp.Components = append(rp.Components, ra2)

	for _, v := range rp.Components {
		t.Log(v)
	}
	sort.Sort(rp)
	if rp.Components[0].Type == "Level1" {
		t.Error("Error Sorting")
	}
	if rp.Components[1].Type == "Level2" {
		t.Error("Error Sorting")
	}
	if rp.Components[2].Type == "Level3" {
		t.Error("Error Sorting")
	}
	if rp.Components[3].Type == "Level4" {
		t.Error("Error Sorting")
	}
}
