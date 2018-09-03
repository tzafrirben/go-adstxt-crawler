package main

import (
	"encoding/json"
	"log"

	"github.com/tzafrirben/go-adstxt-crawler/adstxt"
)

func main() {
	domains := []string{
		"http://example.com",
		"http://test.com",
	}

	requests := make([]*adstxt.Request, len(domains))
	for index, d := range domains {
		r, _ := adstxt.NewRequest(d)
		requests[index] = r
	}

	adstxt.GetMultiple(requests, adstxt.HandlerFunc(handler))
}

func handler(req *adstxt.Request, res *adstxt.Response, err error) {
	if err != nil {
		log.Println(err)
	} else {
		showResults(res.Records)
	}
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
