package geo_test

import (
	"testing"

	"../geo"
	//"../shared"
)

func TestAddress(t *testing.T) {
	/*
		if shared.Config.Init("../config.json") == false {
			t.Fatal("map config.json loading error")
		}
	*/

	a := geo.NewAddress()
	buf := `서울특별시|강남구|역삼동||역삼로 306|개나리 래미안|106동|753번지`
	a.CSVIn(buf)
	a.Print()
	// Change State
	a.SetState("경기도")
	a.Print()
	t.Log(a.CSVOut())
}
