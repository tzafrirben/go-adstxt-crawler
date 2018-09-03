package adstxt

import (
	"bufio"
	"bytes"
	"fmt"
	"sync"
	"time"
)

// Get crawl and parse Ads.txt file from remote host based on Ads.txt Specification Version 1.0.1
// https://iabtechlab.com/wp-content/uploads/2017/09/IABOpenRTB_Ads.txt_Public_Spec_V1-0-1.pdf
func Get(req *Request) (*Response, error) {
	c := newCrawler()

	// send Ads.txt request to remote server unit and parse response. In case of redirect, follow redirect URL and read Ads.txt
	// from the source of redirect
	for {
		// Send request to remote server
		res, err := c.sendRequest(req)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		// handle Ads.txt response
		switch {
		// the server response indicates redirect (301, 302, 307 status codes), follow redirect and read Ads.txt
		// file from the   source of the redirect
		case 300 <= res.StatusCode && res.StatusCode < 400:
			redirect, err := c.handleRedirect(req, res)
			if err != nil {
				return nil, err
			}
			req.URL = redirect
		// client error in remote server response
		case 400 <= res.StatusCode && res.StatusCode < 500:
			return nil, fmt.Errorf(errHTTPClientError, res.Status, req.Domain, req.URL)
		// the server response indicates Success (HTTP 2xx Status Code,) read and parse the content of the Ads.txt file
		case res.StatusCode == 200:
			body, err := c.parseBody(req, res)
			if err != nil {
				return nil, err
			}

			// return new resposne
			records, err := ParseBody(body)
			if err != nil {
				return nil, err
			}

			// Ads.txt response
			r := &Response{
				Request: req,
				Records: records,
				// Ads.txt file default expiration date is set to 7 days (secion 3.6 EXPIRATION of IAB Ads.txt specification)
				Expires: time.Now().UTC().AddDate(0, 0, 7),
			}

			// parse Ads.txt expiration date from response (else default expiration time is used)
			expires, err := c.parseExpires(res)
			if err == nil {
				r.Expires = expires
			}

			return r, nil
		// un known HTTP status
		default:
			return nil, fmt.Errorf(errHTTPGeneralError, res.Status, req.Domain, req.URL)
		}
	}
}

// GetMultiple crawl and parse multiple Ads.txt files from remote hosts based on Ads.txt Specification Version 1.0.1
// https://iabtechlab.com/wp-content/uploads/2017/09/IABOpenRTB_Ads.txt_Public_Spec_V1-0-1.pdf
func GetMultiple(req []*Request, h Handler) {
	var wg sync.WaitGroup
	wg.Add(len(req))

	// buffer of channels to handle response
	for _, r := range req {
		go func(r *Request) {
			res, err := Get(r)
			h.Handle(r, res, err)
			// Decrement the counter when the goroutine completes.
			defer wg.Done()
		}(r)
		// Sleep is not mandatory, but since crawling remote hosts might take time we allow program to "sleep" before
		// proccessing next request
		time.Sleep(100 * time.Millisecond)
	}

	// Wait for all Requests to complete
	wg.Wait()
}

// ParseBody parse Ads.txt file based on Ads.txt Specification Version 1.0.1
// https://iabtechlab.com/wp-content/uploads/2017/09/IABOpenRTB_Ads.txt_Public_Spec_V1-0-1.pdf
func ParseBody(b []byte) (*Records, error) {
	// use custom split function to sunpport different end-of-line marker (CR, CRLF etc)
	split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}
		if i := bytes.IndexAny(data, "\r\n"); i >= 0 {
			if data[i] == '\n' {
				// We have a line terminated by single newline.
				return i + 1, data[0:i], nil
			}
			advance = i + 1
			if len(data) > i+1 && data[i+1] == '\n' {
				advance++
			}
			return advance, data[0:i], nil
		}
		// If we're at EOF, we have a final, non-terminated line. Return it.
		if atEOF {
			return len(data), data, nil
		}
		// Request more data.
		return 0, nil, nil
	}

	scanner := bufio.NewScanner(bytes.NewReader(b))
	scanner.Split(split)

	// loop over Ads.txt file lines and parse each line
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return parseRecords(lines), nil
}
