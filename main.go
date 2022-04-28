//Rikki Naldo Napitupulu (GLNG026ONL013)

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"text/template"
	"time"
)

type CheckData struct {
	Status struct {
		Water int `json:"water"`
		Wind  int `json:"wind"`
	} `json:"status"`
}

func main() {
	go ReloadOnJson()
	http.HandleFunc("/", ReloadOnWeb)
	http.Handle("/asset/", http.StripPrefix("/asset/", http.FileServer(http.Dir("asset"))))
	fmt.Println("listening on PORT:", ":9090")
	http.ListenAndServe(":9090", nil)
}

func ReloadOnJson() {
	for {
		minimum := 1
		maximum := 100
		wind := rand.Intn(maximum-minimum) + minimum
		water := rand.Intn(maximum-minimum) + minimum

		data := CheckData{}
		data.Status.Wind = wind
		data.Status.Water = water

		jsonData, err := json.Marshal(data)

		if err != nil {
			log.Fatal("error occured while marshalling status data:", err.Error())
		}
		err = ioutil.WriteFile("data.json", jsonData, 0644)

		if err != nil {
			log.Fatal("error occured while writing data to data.json file", err.Error())
		}
		time.Sleep(15 * time.Second)
	}
}

func ReloadOnWeb(w http.ResponseWriter, r *http.Request) {
	fileData, err := ioutil.ReadFile("data.json")

	if err != nil {
		log.Fatal("error occured while reading data from data.json file", err.Error())
	}

	var checkData CheckData

	err = json.Unmarshal(fileData, &checkData)
	if err != nil {
		log.Fatal("error occured while unMarshalling from data.json file", err.Error())
	}

	waterVal := checkData.Status.Water
	windVal := checkData.Status.Wind

	var (
		waterStat string
		windStat  string
	)

	waterValue := strconv.Itoa(waterVal)
	windValue := strconv.Itoa(windVal)

	switch {
	case waterVal <= 5:
		waterStat = "Aman"
	case waterVal >= 6 && waterVal <= 8:
		waterStat = "Siaga"
	case waterVal > 8:
		waterStat = "Bahaya"
	default:
		waterStat = "Water Value not defined"
	}

	switch {
	case windVal <= 6:
		windStat = "Aman"
	case windVal >= 7 && windVal <= 15:
		windStat = "Siaga"
	case windVal > 15:
		windStat = "Bahaya"
	default:
		windStat = "Wind Value not defined"
	}

	data := map[string]string{
		"waterStat":  waterStat,
		"windStat":   windStat,
		"waterValue": waterValue,
		"windValue":  windValue,
	}

	templ, err := template.ParseFiles("index.html")

	if err != nil {
		log.Fatal("error parsing html:", err.Error())
	}

	templ.Execute(w, data)

}
