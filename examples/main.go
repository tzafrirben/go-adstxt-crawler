package main

import (
	"encoding/json"
	"go-adstxt-crawler/adstxt"
	"io/ioutil"
	"log"
)

func main() {
	// parse local file
	body, err := ioutil.ReadFile("<path-to-local-ads.txt file>")
	if err != nil {
		log.Fatal(err)
	}
	rec, err := adstxt.ParseBody(body)
	if err != nil {
		log.Fatal(err)
	}
	showResults(rec)
}

func showResults(r *adstxt.Records) {
	if len(r.Warnings) > 0 {
		log.Println("Warnings:")
		for _, w := range r.Warnings {
			j, _ := json.Marshal(w)
			log.Println(string(j))
		}
	}

	if len(r.DataRecords) > 0 {
		log.Println("Data Records:")
		for _, r := range r.DataRecords {
			j, _ := json.Marshal(r)
			log.Println(string(j))
		}
	}

	if len(r.Variables) > 0 {
		log.Println("Variables:")
		for _, v := range r.Variables {
			j, _ := json.Marshal(v)
			log.Println(string(j))
		}
	}
}
