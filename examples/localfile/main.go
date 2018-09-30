package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/tzafrirben/go-adstxt-crawler/adstxt"
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
		log.Printf("Warnings: [%d]", len(r.Warnings))
		for _, w := range r.Warnings {
			j, _ := json.Marshal(w)
			log.Println(string(j))
		}
	}

	if len(r.DataRecords) > 0 {
		log.Printf("Data Records: [%d]", len(r.DataRecords))
		for _, r := range r.DataRecords {
			j, _ := json.Marshal(r)
			log.Println(string(j))
		}
	}

	if len(r.Variables) > 0 {
		log.Printf("Variables: [%d]", len(r.Variables))
		for _, v := range r.Variables {
			j, _ := json.Marshal(v)
			log.Println(string(j))
		}
	}
}
