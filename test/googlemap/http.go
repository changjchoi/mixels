package main

import (
	"fmt"
	_ "io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"../../../hicup/hishape"
	"../../shared"
)

type Level7File struct {
	File0 shared.RegionFile
	File1 shared.RegionFile
}

func (l *Level7File) Open() {
	filename0 := shared.Config.OutputPath + "/" +
		shared.Config.Level2File["Level7"][0]
	filename1 := shared.Config.OutputPath + "/" +
		shared.Config.Level2File["Level7"][1]

	log.Print(filename0)
	log.Print(filename1)

	l.File0.Filename = filename0
	l.File1.Filename = filename1
	err := l.File0.Open()
	if err != nil {
		log.Fatal("level71 file open error")
	}
	err = l.File1.Open()
	if err != nil {
		log.Fatal("level72 file open error")
	}
	log.Println("file0 count =", l.File0.ItemCount())
	log.Println("file1 count =", l.File1.ItemCount())
}

func (l *Level7File) ReadPoints(id int) []hishape.PointInt32 {
	count := l.File0.ItemCount()
	// id = [0, n), [n, k)
	if id >= count {
		points, _ := l.File1.ReadPoints(id - count)
		return points
	}
	points, _ := l.File0.ReadPoints(id)
	return points
}

var l7f Level7File

func init() {
	if shared.Config.Init("../../config.json") == false {
		log.Fatal("map config.json loading error")
	}
	l7f.Open()
}

func fileAccess(w http.ResponseWriter, r *http.Request) {
	fmt.Println("access : " + r.URL.Path[1:])
	http.ServeFile(w, r, "../googlemap/"+r.URL.Path[1:])
}

func jsonAccess(w http.ResponseWriter, r *http.Request) {
	service_type := r.URL.Query()
	if value, ok := service_type["id"]; ok {
		w.Header().Set("Content-Type", "application/json")

		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers",
				"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		}
		// the value is id number which type is []string
		//log.Println("raw id =", value)
		id, _ := strconv.Atoi(value[0])

		log.Println("Get id =", id)
		// read a map
		pts := l7f.ReadPoints(id)
		// send points
		var first_object string
		w.Write([]byte("["))
		for i, v := range pts {
			stream := fmt.Sprintf("{\"lat\":%f,\"lng\":%f},",
				float64(v.Y)/hishape.LLCONV, float64(v.X)/hishape.LLCONV)

			w.Write([]byte(stream))
			if i == 0 {
				first_object = stream
			}
		}
		first_object = strings.TrimRight(first_object, ",")
		w.Write([]byte(first_object))
		w.Write([]byte("]"))
	}

}

func main() {
	go func() {
		web_serve_mux := http.NewServeMux()
		web_serve_mux.HandleFunc("/", fileAccess)
		http.ListenAndServe(":8001", web_serve_mux)
	}()

	json_serve_mux := http.NewServeMux()
	json_serve_mux.HandleFunc("/getpoint", jsonAccess)
	http.ListenAndServe(":8002", json_serve_mux)
}
