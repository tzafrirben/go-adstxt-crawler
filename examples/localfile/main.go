package main

import (
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
	log.Println(rec)
}
