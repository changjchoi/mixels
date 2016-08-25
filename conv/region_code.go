package conv

import (
	"bytes"
	"encoding/csv"
	"io/ioutil"
	"log"

	. "../shared"
)

type RegionCode struct {
	RoadNameCode map[string]string
	Level1       map[string]string
	Level2       map[string]string
	Level3A      map[string]string
	Level3B      map[string]string
}

func NewRegionCode() *RegionCode {
	return &RegionCode{}
}

func (r *RegionCode) Open() error {
	/*
		filename := Config.OutputPath + "/" + Config.RegionCodeFile["RoadNameCode"]
		records := r.LoadCSV(filename)
		r.RoadNameCode = make(map[string]string, len(records))
		for _, v := range records {
			r.RoadNameCode[v[0]] = v[1]
		}

		filename = Config.OutputPath + "/" + Config.RegionCodeFile["Level1"]
		records = r.LoadCSV(filename)
		r.Level1 = make(map[string]string, len(records))
		for _, v := range records {
			r.Level1[v[0]] = v[1]
		}

		filename = Config.OutputPath + "/" + Config.RegionCodeFile["Level2"]
		records = r.LoadCSV(filename)
		r.Level2 = make(map[string]string, len(records))
		for _, v := range records {
			r.Level2[v[0]] = v[1]
		}

		filename = Config.OutputPath + "/" + Config.RegionCodeFile["Level3A"]
		records = r.LoadCSV(filename)
		r.Level3A = make(map[string]string, len(records))
		for _, v := range records {
			r.Level3A[v[0]] = v[1]
		}

		filename = Config.OutputPath + "/" + Config.RegionCodeFile["Level3B"]
		records = r.LoadCSV(filename)
		r.Level3B = make(map[string]string, len(records))
		for _, v := range records {
			r.Level3B[v[0]] = v[1]
		}
	*/
	r.RoadNameCode = r.LoadCode("RoadNameCode")
	r.Level1 = r.LoadCode("Level1")
	r.Level2 = r.LoadCode("Level2")
	r.Level3A = r.LoadCode("Level3A")
	r.Level3B = r.LoadCode("Level3B")
	return nil
}

func (r RegionCode) LoadCode(level string) map[string]string {
	//var code map[string]string
	filename := Config.InputPath + "/" + Config.RegionCodeFile[level]
	records := r.LoadCSV(filename)
	code := make(map[string]string, len(records))
	for _, v := range records {
		code[v[0]] = v[1]
	}
	return code
}

func (r *RegionCode) LoadCSV(filename string) [][]string {
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal("file open error :", filename)
	}
	c := csv.NewReader(bytes.NewReader(f))
	c.Comma = '|'
	records, err := c.ReadAll()
	if err != nil {
		log.Fatal("read records error")
	}
	return records
}
