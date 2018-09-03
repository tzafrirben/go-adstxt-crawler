# go-adstxt-crawler
[Ads.txt](https://iabtechlab.com/ads-txt-about/) crawler and parser based on [IAB Ads.txt Specification Version 1.0.1](https://iabtechlab.com/wp-content/uploads/2017/09/IABOpenRTB_Ads.txt_Public_Spec_V1-0-1.pdf) implemented in Go

This library provides a mechanism for obtaining and parsing Ads.txt file from websites, or parse your local copy of Ads.txt file

# Motivation
There are some nice online tools for crawling and validating Ads.txt files (for example [Ads.txt validator](https://adstxt.adnxs.com) from AppNexus or another [Ads.txt Validator](https://www.adstxtvalidator.com) by AdReform) that follows [IAB Ads.txt Specification Version 1.0.1](https://iabtechlab.com/wp-content/uploads/2017/09/IABOpenRTB_Ads.txt_Public_Spec_V1-0-1.pdf). 
However, you cannot easily use those tools for massive site scanning since they do not provide free API.

There are also few open source projects I've found Github for scanning Ads.txt files, but at least the ones that I've tried were not fully competible with latest [Ads.txt Spec](https://iabtechlab.com/wp-content/uploads/2017/09/IABOpenRTB_Ads.txt_Public_Spec_V1-0-1.pdf). You can use them, of course, and they would do a decent job, but they are lacking a good validation mechanism to ensure that Ads.txt format is correct and follows the official spec.

This Ads.txt library allows massive sites crawling, and follows [IAB Ads.txt Specification Version 1.0.1](https://iabtechlab.com/wp-content/uploads/2017/09/IABOpenRTB_Ads.txt_Public_Spec_V1-0-1.pdf) to validate that the Ads.txt file is valid

# Examples
You can see [examples](https://github.com/tzafrirben/go-adstxt-crawler/tree/master/examples) folder for a short example of adstxt library 3 main methods: adstxt.Get to fetch and parse single Ads.txt file from a remote host, adstxt.GetMultiple to fetch and parse multiple Ads.txt files form different hosts or adstxt.ParseBody that can be used to parse the content of a local Ads.txt file

```go
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

Or get Ads.txt files for multiple hosts simultaneously
```go
// define handler function to handle Ads.txt response
h := func(req *Request, res *Response, err error) {
  for _, r := range res.DataRecords { ... }
  for _, v := range res.Variables { ... }
  for _, w := range res.Warnings { ... }
}

// collection of domains to validate
domains := []string{
  "http://example.com",
  "http://test.com",
}

requests := make([]*adstxt.Request, len(domains))
for index, d := range domains {
  r, _ := adstxt.NewRequest(d)
  requests[index] = r
}

adstxt.GetMultiple(requests, adstxt.HandlerFunc(h))
```

You can also parse local Ads.txt file in a similar way
```go
body, err := ioutil.ReadFile("/<path_to>/ads.txt")
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

# Import as a Library
import "github.com/tzafrirben/go-adstxt-crawler/adstxt" and you can use adstxt library in your code

# ToDo
- robots.txt file on remote host is ignored by crawler, a good practice will be to scan this file first (as specified in Ads.txt specification)

## LICENSE

MIT

## Author
[Tzafrir Ben Ami](https://github.com/tzafrirben) (a.k.a. tzafrirben)
