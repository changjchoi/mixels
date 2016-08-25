package conv

import (
	//"bytes"
	//"encoding/csv"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	. "../shared"
)

func DumpAddress(idx_list []int) {
	// Delete Old adddress.tsv
	for _, v := range Config.TSVFile {
		if err := os.Remove(Config.OutputPath + "/" + v); err != nil {
			log.Println("file remove error :", err)
		}
	}
	// Only Once !
	rc := NewRegionCode()
	rc.Open()
	//
	prev_item, prev_index := 0, 1
	for _, v := range Config.Level2File["Level7"] {
		rf := NewRegionFile(Config.OutputPath + "/" + v)
		err_f := rf.Open()
		if err_f != nil {
			log.Fatal("Open Error :", err_f)
		}
		//
		log.Println("Total Items :", rf.ItemCount())
		//
		if len(idx_list) == 0 {
			prev_index = DumpSection(-1, prev_item, prev_index, rf, rc)
		} else {
			// Test Output Code
			for _, i := range idx_list {
				DumpSection(i, prev_item, prev_index, rf, rc)
			}
			break
		}
		//
		rf.Close()
		prev_item = rf.ItemCount()
	}
}

func PushSphinx(arg string) {
	// run script
	var (
		cmd_out []byte
		err     error
	)
	cmd_path := "../hicup/sphinx"
	cmd_name := cmd_path + "/sphinx.sh"
	cmd_args := []string{arg, Config.OutputPath, cmd_path}
	if cmd_out, err = exec.Command(cmd_name, cmd_args...).Output(); err != nil {
		log.Fatal("Error :", err, string(cmd_out))
	}
	log.Println(string(cmd_out))
}

func DumpSection(n int, p int, pd int, rf *RegionFile, rc *RegionCode) int {
	var err error
	af := NewAddressField(rc)
	//
	fd := make([]*os.File, len(Config.TSVFile))
	for i, v := range Config.TSVFile {
		fd[i], err = os.OpenFile(Config.OutputPath+"/"+v,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal("file open err", err)
		}
		defer func() {
			fd[i].Close()
		}()
	}
	//
	start, end := 0, 0
	if n == -1 {
		start, end = 0, rf.ItemCount()
	} else {
		start, end = n*SPLITCOUNT, (n+1)*SPLITCOUNT
		if end > rf.ItemCount() {
			end = rf.ItemCount()
		}
	}
	//
	back_addr := ""
	index := pd
	t0, t1 := time.Now(), time.Now()
	for i := start; i < end; i++ {
		if i != start && i%50000 == 0 {
			t1 = time.Now()
			log.Printf("[%d] %v ", i, t1.Sub(t0))
			t0 = t1
			//w.Flush()
		}
		cur_addr, _ := rf.ReadAttribute(i)
		// Check duplicated Address. Therefore an Address is not sorted
		// but the same address is to be a limited area
		// !! Sorted !!
		if cur_addr == back_addr {
			continue
		}
		af.AddCSVField(cur_addr)
		//log.Println("FFF:", cur_addr)
		fidx := func() int {
			r := 0
			for i, v := range Config.Level2File["Level7"] {
				if v == filepath.Base(rf.Filename) {
					r = i
					break
				}
			}
			return r
		}()
		records := []string{
			fmt.Sprintf("%d", index),
			fmt.Sprintf("%s", af.State()),
			fmt.Sprintf("%s", af.City()),
			fmt.Sprintf("%s", af.Town()),
			fmt.Sprintf("%s", af.TownLi()),
			fmt.Sprintf("%s %s %s %s",
				af.StreetName(),
				af.BuildingName(),
				af.BuildingNameSub(),
				af.Jibun()),
			fmt.Sprintf("{ id:%d, i:\"%d\" }\n", i+p, fidx),
		}
		//var buf bytes.Buffer
		//w := csv.NewWriter(fd)
		record := strings.Join(records, "\t")
		if _, err := fd[index%4].WriteString(record); err != nil {
			log.Fatal("Error write csv")
		}
		// Save before address string
		back_addr = cur_addr
		index += 1
	}
	return index
}
