package adstxt

// The Handler interface is used to process Ads.txt requests. It is similar to the
// net/http.Handler interface.
type Handler interface {
	Handle(*Request, *Response, error)
}

// A HandlerFunc is a function signature that implements the Handler interface. A function
// with this signature can thus be used as a Handler.
type HandlerFunc func(*Request, *Response, error)

// Handle is the Handler interface implementation for the HandlerFunc type.
func (h HandlerFunc) Handle(req *Request, res *Response, err error) {
	h(req, res, err)
}
