package main_test

import (
	"flag"
	// "io"
	// "io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"../hicup/hicom"
	"./conv"
	"./shared"
)

var index_list = flag.String("index", "", "index list to ...??")

func init() {
	if shared.Config.Init("./config.json") == false {
		//log.Println("map config.json loading error")
	}
}

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

func TestDumpAddress(t *testing.T) {
	time_start := time.Now()
	log.Println("address string........")
	var idx_list []int
	if len(*index_list) == 0 {
	} else {
		for _, v := range strings.Split(*index_list, ",") {
			l, err := strconv.Atoi(v)
			if err != nil {
				log.Println("Argument parsing error :", err)
			}
			idx_list = append(idx_list, l)
		}
	}
	conv.DumpAddress(idx_list)
	time_end := time.Now()
	log.Println("Finished :", time_end.Sub(time_start))
}

func TestPushSphinx(t *testing.T) {
	time_start := time.Now()
	log.Println("Starting to make Sphinx Index ........")
	conv.PushSphinx("start")
	time_end := time.Now()
	log.Println("Finished :", time_end.Sub(time_start))
}

func TestConvAll(t *testing.T) {
	time_start := time.Now()
	log.Println("Converting........")
	convertFile(t)
	time_end := time.Now()
	log.Println("Conversion Finished :", time_end.Sub(time_start))
}

func convertFile(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(len(shared.Config.Mapping))
	for key, value := range shared.Config.Mapping {
		// key : filename, outfile : output filename, prj : projection name
		go func(key string, outfile string, prj string) {
			t2 := time.Now()
			mapfile := shared.Config.InputPath + "/" + key
			defer func(filename string) {
				t3 := time.Now()
				log.Println("Map :", filename)
				log.Println("converting time =", t3.Sub(t2))
				wg.Done()
			}(key)
			basefile := filepath.Base(key)
			basedir := shared.Config.InputPath + "/" + filepath.Dir(key)
			if basefile == "*" {
				files := hicom.DirListSuffix(basedir, ".shp")
				if len(files) == 0 {
					log.Println("There are no processing file")
					return
				}
				output := shared.Config.OutputPath + "/" + outfile
				conv := conv.NewConvertShapeFile("", output)
				conv.SetDirConv(true)
				// open proj4
				conv.Open(prj)
				for i, v := range files {
					// Projection Info..
					log.Println("Adding...", "(", i+1, ")", filepath.Base(v))
					conv.Reopen(v)
					conv.Iterate()
				}
				conv.Close()
				conv.Save()
			} else {
				output := shared.Config.OutputPath + "/" + outfile
				conv := conv.NewConvertShapeFile(mapfile, output)
				// Projection Info..
				conv.Open(prj)
				conv.Iterate()
				conv.Close()
				conv.Save()
			}
		}(key, value.File, value.Projection)
	}
	// Waiting here while all thread comming here !
	wg.Wait()
}
