package requestx

import "net/http"

func NewClient(opts ...Options) *Request {
	req := &Request{}

	opts0 := Options{}
	if len(opts) > 0 {
		opts0 = opts[0]
	}

	req.SetOptions(opts0)

	return req
}

func Get(uri string, opts ...Options) (*Response, error) {
	r := NewClient()
	return r.Request(http.MethodGet, uri, opts...)
}

func Post(uri string, opts ...Options) (*Response, error) {
	r := NewClient()
	return r.Request(http.MethodPost, uri, opts...)
}

func Put(uri string, opts ...Options) (*Response, error) {
	r := NewClient()
	return r.Request(http.MethodPut, uri, opts...)
}

func Patch(uri string, opts ...Options) (*Response, error) {
	r := NewClient()
	return r.Request(http.MethodPatch, uri, opts...)
}

func Delete(uri string, opts ...Options) (*Response, error) {
	r := NewClient()
	return r.Request(http.MethodDelete, uri, opts...)
}
