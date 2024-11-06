package requestx

import (
	"crypto/tls"
	"time"
)

type Options struct {
	Debug        bool
	BaseURI      string
	Timeout      float32
	timeout      time.Duration
	Query        any
	Headers      map[string]any
	Cookies      any
	FormParams   map[string]any
	JSON         any
	XML          any
	Multipart    []FormData
	Proxy        string
	Certificates []tls.Certificate
}

func mergeOptions(opts0 Options, opts ...Options) Options {
	for _, opt := range opts {
		if opt.Debug {
			opts0.Debug = true
		}
		if opt.BaseURI != "" {
			opts0.BaseURI = opt.BaseURI
		}
		if opt.Timeout > 0 {
			opts0.Timeout = opt.Timeout
		}
		if opt.timeout > 0 {
			opts0.timeout = opt.timeout
		}
		if opt.Query != nil {
			opts0.Query = opt.Query
		}
		if opt.Headers != nil {
			opts0.Headers = opt.Headers
		}
		if opt.Cookies != nil {
			opts0.Cookies = opt.Cookies
		}
		if opt.FormParams != nil {
			opts0.FormParams = opt.FormParams
		}
		if opt.JSON != nil {
			opts0.JSON = opt.JSON
		}
		if opt.XML != nil {
			opts0.XML = opt.XML
		}
		if opt.Multipart != nil {
			opts0.Multipart = opt.Multipart
		}
		if opt.Proxy != "" {
			opts0.Proxy = opt.Proxy
		}
		if opt.Certificates != nil {
			opts0.Certificates = opt.Certificates
		}
	}
	return opts0
}
