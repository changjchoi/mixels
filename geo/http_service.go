package geo

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"../../hicup/hishape"
	"../shared"
)

// Copied from hicup
func Service(sa *SearchArea) {
	http.HandleFunc("/geocoding", func(w http.ResponseWriter, r *http.Request) {
		t0 := time.Now()
		service_type := r.URL.Query()
		if _, ok := service_type["latlon"]; ok {
			pt, err := parselatlon(service_type)
			if err != nil {
				respondErrorJson(err, w)
				return
			}
			result := sa.SearchAddress(pt)
			io.WriteString(w, string(result))
		} else if _, ok := service_type["address"]; ok {
			param, _ := parseAddress(service_type)
			result := sa.SearchPoint(param)
			io.WriteString(w, string(result))
		} else {
			respondErrorJson(fmt.Errorf("parse error"), w)
		}
		log.Println("Service Time :", time.Since(t0))
	})
	http.ListenAndServe(shared.Config.ServicePort, http.DefaultServeMux)
}

// /geocoding?address=seoul+gangnamgu+yeoksamdong+yeoksamro+306
func parseAddress(p map[string][]string) (string, error) {
	addr := p["address"]
	search_param := strings.Replace(addr[0], "+", " ", -1)
	//search_param := strings.Split(addr[0], ",")
	//last := len(search_param)
	// Search Query Make - last word scale up by twice
	//search_param[last-1] = search_param[last-1] + "^2"
	//log.Println("search param:", strings.Join(search_param, " "))
	//return strings.Join(search_param, " "), nil
	return search_param, nil
}

// /geocoding?latlon=34,127
func parselatlon(p map[string][]string) (hishape.Point, error) {
	lonlat, _ := p["latlon"]
	ll_str := strings.Split(lonlat[0], ",")
	if len(ll_str) != 2 {
		return hishape.Point{},
			fmt.Errorf("put in lon and lat like ?latlon=lat,lon")
	}
	lat, err := strconv.ParseFloat(ll_str[0], 64)
	if err != nil {
		return hishape.Point{}, fmt.Errorf("lat parse error")
	}
	lon, err := strconv.ParseFloat(ll_str[1], 64)
	if err != nil {
		return hishape.Point{}, fmt.Errorf("lon parse error")
	}
	return hishape.Point{lon, lat}, nil
}

var jsonMarshalErr string = "{\"error\": \"json marshal error\"}"

type ResError struct {
	Err string `json:"error"`
}

func respondErrorJson(err error, w http.ResponseWriter) {
	res := ResError{err.Error()}
	b, e := json.Marshal(res)
	if e != nil {
		io.WriteString(w, jsonMarshalErr)
		return
	}
	io.WriteString(w, string(b[:]))
}
