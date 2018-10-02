package main

import (
	"log"

	"github.com/tzafrirben/go-adstxt-crawler/adstxt"
)

func main() {

	// fetch and download Ads.txt file from remote host
	req, err := adstxt.NewRequest("http://fark.com")
	if err != nil {
		log.Fatal(err)
	}
	res, err := adstxt.Get(req)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(res.Records)
}
