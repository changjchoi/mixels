package main

import (
	"flag"
	// "fmt"
	"log"
	"os"
	_ "path/filepath"
	"runtime/pprof"
	_ "strconv"
	_ "strings"
	"sync"
	"time"
	// "unsafe"

	"../hicup/hishape"
	"./geo"
	"./shared"
)

var profile_list = []struct {
	p hishape.Point
	v string
}{
	{hishape.Point{Y: 37.49884, X: 127.04750},
		"서울특별시,강남구,역삼동,개나리 래미안 106"},
	{hishape.Point{Y: 37.512342, X: 127.058881},
		"서울특별시,강남구,삼성1동,코엑스"},
	{hishape.Point{Y: 37.468275, X: 126.631915},
		"인천광역시,중구,신흥동,1가 10번지"},
	{hishape.Point{Y: 37.451545, X: 126.653168},
		"인천광역시,남구,인하대학교 5호관"},
	{hishape.Point{Y: 37.359630, X: 127.105449},
		"경기도,성남시,분당구,네이버"},
	{hishape.Point{Y: 36.372326, X: 127.361640},
		"대전광역시,유성구,구성동,카이스트"},
	{hishape.Point{Y: 34.298957, X: 126.527923},
		"전라남도,해남군,송지면,땅끝파출소"},
}

var log_write = flag.Bool("w", false, "write a log in the file")
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

// Called from golang for each files
func init() {
	if shared.Config.Init("./config.json") == false {
		log.Fatal("map config.json loading error")
	}
}

func main() {
	t0 := time.Now()
	// parse command-line arguments
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	// if you set, "mixels -w". it writes a log in the file
	if *log_write {
		f, err := os.OpenFile("logfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		log.SetOutput(f)
	}
	sa := geo.SearchArea{}
	// Init() for Text Search
	sa.Init()
	loading(&sa)
	sa.PrepareSearch()
	t1 := time.Now()
	log.Println("Total Loading Time = ", t1.Sub(t0))
	log.Println("The Num of Area = ", sa.Count())
	//http service
	if *cpuprofile != "" {
		for _, v := range profile_list {
			sa.SearchAddress(v.p)
		}
	} else {
		geo.Service(&sa)
	}
}

func loading(sa *geo.SearchArea) {
	log.Println("Loading all index data in a file... wait...")
	var wg sync.WaitGroup
	wg.Add(len(shared.Config.AreaLevel))
	for key, _ := range shared.Config.AreaLevel {
		go func(key string) {
			t2 := time.Now()
			// Empty Area
			area := geo.Area{}
			area.Filename = shared.Config.OutputPath + "/" + key
			//
			defer func(filename string) {
				t3 := time.Now()
				log.Println("Map :", filename, "loading time =", t3.Sub(t2))
				wg.Done()
			}(area.Filename)
			err_f := area.Open()
			if err_f != nil {
				log.Fatal("Open Error :", err_f)
			}
			sa.Add(&area)
		}(key)
	}
	wg.Wait()
}
