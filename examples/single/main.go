package main

import (
	"encoding/json"
	"log"

	"github.com/tzafrirben/go-adstxt-crawler/adstxt"
)

func main() {

	// fetch and download Ads.txt file from remote host
	req, err := adstxt.NewRequest("http://example.com")
	if err != nil {
		log.Fatal(err)
	}
	res, err := adstxt.Get(req)
	if err != nil {
		log.Fatal(err)
	}
	showResults(res.Records)
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
