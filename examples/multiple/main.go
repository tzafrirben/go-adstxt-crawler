package main

import (
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
		log.Println(res.Records)
	}
}
