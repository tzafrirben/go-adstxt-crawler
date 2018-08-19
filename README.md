# go-adstxt-crawler
[Ads.txt](https://iabtechlab.com/ads-txt-about/) crawler and parser based on [IAB Ads.txt Specification Version 1.0.1](https://iabtechlab.com/wp-content/uploads/2017/09/IABOpenRTB_Ads.txt_Public_Spec_V1-0-1.pdf) implemented in Go

This library provides mechanism for obtaining and parsing Ads.txt file from websites, or parse your local copy of Ads.txt file

# Examples
You can see main.go file for a short example of the library 2 main methods: adstxt.Get to fetch and parse Ads.txt file from remote host, or adstxt.ParseBody to parse the content of a local Ads.txt file

```
req, err := adstxt.NewRequest("example.com")
if err != nil {
  log.Fatal(err)
}
res, err := adstxt.Get(req)
if err != nil {
  log.Fatal(err)
}
// res now holds Ads.txt file DataRecords, Variables and Warnings for Ads.txt parse warnings
for _, r := range res.DataRecords { ... }
for _, v := range res.Variables { ... }
for _, w := range res.Warnings { ... }
```

You can also parse local Ads.txt file in a similar way
```
body, err := ioutil.ReadFile("/<path>/ads.txt")
if err != nil {
  log.Fatal(err)
}
rec, err := adstxt.ParseBody(body)
if err != nil {
  log.Fatal(err)
}
// rec now holds Ads.txt file DataRecords, Variables and Warnings for Ads.txt parse warnings
for _, r := range rec.DataRecords { ... }
for _, v := range rec.Variables { ... }
for _, w := range rec.Warnings { ... } 
```

# import as a Library
import "github.com/tzafrirben/go-adstxt-crawler/adstxt"

## LICENSE

MIT

## Author
Tzafrir Ben Ami (a.k.a. tzafrirben)
