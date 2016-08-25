package shared_test

import (
	. "../shared"
	"log"
	"testing"
	"time"
)

func init() {
	if Config.Init("../config.json") == false {
		log.Fatal("map config.json loading error")
	}
}

func TestRegionIndex(t *testing.T) {
	t.Log("Test Region Index")
	file := Config.OutputPath + "/" + Config.Level2File["Level3B"][0]
	bin := RegionFile{Filename: file}
	err_f := bin.Open()
	if err_f != nil {
		t.Fatal("Open Error")
	}
	bin.PrintHeader()
	bin.PrintBounds(0)
	bin.PrintAreaRecords(0)
	// last item
	last := int(bin.Header.ItemCount)
	bin.PrintAreaRecords(last - 1)
	r, _ := bin.ReadAttribute(0)
	//t.Log("data :", r)
	t.Log("Attribute :", r)
	bin.Close()
}

func TestRegionIndex7(t *testing.T) {
	t.Log("Test Region Index Level7")
	file := Config.OutputPath + "/" + Config.Level2File["Level7"][0]
	t.Log(file)
	bin := RegionFile{Filename: file}
	err_f := bin.Open()
	if err_f != nil {
		t.Fatal("Open Error")
	}
	//bin.PrintHeader()
	//bin.PrintBounds(0)
	// Test Edge Base Address
	count := int(bin.Header.ItemCount) / SPLITCOUNT
	for i := 0; i < count-1; i += 1 {
		start, _ := bin.ReadAreaRecord(i * SPLITCOUNT)
		end, _ := bin.ReadAreaRecord(i*SPLITCOUNT + SPLITCOUNT - 1)
		if start.AttributeBase != (end.PointBase + end.PointSize) {
			t.Error("Must be the same")
		}
	}
	/*
		bin.PrintAreaRecords(000000)
		bin.PrintAreaRecords(999999)

		bin.PrintAreaRecords(1000000)
		bin.PrintAreaRecords(1999999)

		bin.PrintAreaRecords(2000000)
		bin.PrintAreaRecords(2999999)

		bin.PrintAreaRecords(7000000)
		bin.PrintAreaRecords(7999999)

		bin.PrintAreaRecords(8000000)
		bin.PrintAreaRecords(8999999)

		bin.PrintAreaRecords(9000000)
		bin.PrintAreaRecords(9999999)
	*/

	r, err := bin.ReadAttribute(8280951)
	if err != nil {
		t.Log("Error : ", err)
	}
	//t.Log("data :", r)
	t.Log("Attribute :", r)
	r, err = bin.ReadAttribute(0)
	t.Log("ID: 0", "midong : ", r)
	r, err = bin.ReadAttribute(2)
	t.Log("ID: 2", "midong : ", r)
	r, err = bin.ReadAttribute(270549)
	t.Log("ID: 10", "midong : ", r)
	k, e := bin.ReadPoints(0)
	t.Log("id 0: ", k, e)
	bin.Close()
}

func TestRegionAttribute(t *testing.T) {
	t.Log("Test Region Attribute")
	file := Config.OutputPath + "/" + Config.Level2File["Level7"][0]
	bin := RegionFile{Filename: file}
	err_f := bin.Open()
	if err_f != nil {
		t.Fatal("Open Error")
	}
	t0 := time.Now()
	for i, _ := range bin.AreaRecords {
		bin.ReadAttribute(i)
	}
	t1 := time.Now()
	t.Log("Read All Attribute time :", t1.Sub(t0))
	bin.Close()
}
