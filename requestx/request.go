package requestx

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"net/textproto"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/clbanning/mxj/v2"
	"github.com/spf13/cast"
)

type Request struct {
	opts Options
	cli  *http.Client
	req  *http.Request
	body io.Reader
}

// FormData multipart form data
type FormData struct {
	Name     string
	Contents []byte
	Filename string
	Filepath string
	Headers  map[string]any
}

// SetOptions set request options
func (r *Request) SetOptions(opts Options) {
	r.opts = opts
}

func (r *Request) Get(uri string, opts ...Options) (*Response, error) {
	return r.Request(http.MethodGet, uri, opts...)
}

func (r *Request) Post(uri string, opts ...Options) (*Response, error) {
	return r.Request(http.MethodPost, uri, opts...)
}

func (r *Request) Put(uri string, opts ...Options) (*Response, error) {
	return r.Request(http.MethodPut, uri, opts...)
}

func (r *Request) Patch(uri string, opts ...Options) (*Response, error) {
	return r.Request(http.MethodPatch, uri, opts...)
}

func (r *Request) Delete(uri string, opts ...Options) (*Response, error) {
	return r.Request(http.MethodDelete, uri, opts...)
}

func (r *Request) Options(uri string, opts ...Options) (*Response, error) {
	return r.Request(http.MethodOptions, uri, opts...)
}

func (r *Request) Request(method, uri string, opts ...Options) (*Response, error) {
	r.opts = mergeOptions(r.opts, opts...)
	if !strings.HasPrefix(uri, "http") && strings.HasPrefix(r.opts.BaseURI, "http") {
		uri = r.opts.BaseURI + uri
	}

	if r.opts.Headers == nil {
		r.opts.Headers = make(map[string]any)
	}

	switch method {
	case http.MethodGet, http.MethodDelete:
		req, err := http.NewRequest(method, uri, nil)
		if err != nil {
			return nil, err
		}
		r.req = req
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodOptions:
		r.parseBody()

		req, err := http.NewRequest(method, uri, r.body)
		if err != nil {
			return nil, err
		}

		r.req = req
	default:
		return nil, errors.New("unsupported method")
	}

	r.parseOptions()

	r.parseClient()

	r.parseQuery()

	r.parseHeaders()

	r.parseCookies()

	if r.opts.Debug {
		dump, err := httputil.DumpRequest(r.req, true)
		if err == nil { // 修改判断条件：成功时打印
			fmt.Printf("\n%s\n\n", dump)
		} else {
			fmt.Printf("dump request error: %v\n", err)
		}
	}

	_resp, err := r.cli.Do(r.req)
	resp := &Response{
		resp: _resp,
		req:  r.req,
		err:  err,
	}

	if err != nil {
		if r.opts.Debug {
			fmt.Println(err)
		}

		return resp, err
	}

	if strings.HasPrefix(resp.GetHeaderLine("content-type"), "text/event-stream") {
		resp.parseStream()
		return resp, nil
	}

	body, err := io.ReadAll(_resp.Body)
	_resp.Body.Close()

	resp.body = body
	resp.err = err

	if r.opts.Debug {
		body, _ := resp.GetBody()
		fmt.Println(string(body))
	}
	return resp, nil
}

func (r *Request) parseOptions() {
	if r.opts.Timeout == 0 {
		r.opts.Timeout = 30
	}

	r.opts.timeout = time.Duration(r.opts.Timeout*1000) * time.Millisecond
}

func (r *Request) parseClient() {
	tlsConfig := &tls.Config{}
	if len(r.opts.Certificates) > 0 {
		tlsConfig.Certificates = r.opts.Certificates
	} else {
		tlsConfig.InsecureSkipVerify = true
	}

	tr := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	if r.opts.Proxy != "" {
		proxy, err := url.Parse(r.opts.Proxy)
		if err == nil {
			tr.Proxy = http.ProxyURL(proxy)
		}
	}

	r.cli = &http.Client{
		Timeout:   r.opts.timeout,
		Transport: tr,
	}

	// 防止在对同一主机的请求之间重复使用TCP连接
	r.req.Close = true
}

func (r *Request) parseQuery() {
	switch r.opts.Query.(type) {
	case string:
		str := r.opts.Query.(string)
		r.req.URL.RawQuery = str
	case map[string]string:
		q := r.req.URL.Query()
		for k, v := range r.opts.Query.(map[string]string) {
			q.Set(k, v)
		}
		r.req.URL.RawQuery = q.Encode()
	case map[string]any:
		q := r.req.URL.Query()
		for k, v := range r.opts.Query.(map[string]any) {
			if vv, ok := v.(string); ok {
				q.Set(k, vv)
				continue
			}

			if vv, ok := v.([]string); ok {
				for _, vvv := range vv {
					q.Add(k, vvv)
				}
				continue
			}

			if vv := cast.ToString(v); vv != "" {
				q.Set(k, vv)
			}
		}
		r.req.URL.RawQuery = q.Encode()
	}
}

func (r *Request) parseCookies() {
	switch r.opts.Cookies.(type) {
	case string:
		cookies := r.opts.Cookies.(string)
		r.req.Header.Add("Cookie", cookies)
	case map[string]string:
		cookies := r.opts.Cookies.(map[string]string)
		for k, v := range cookies {
			r.req.AddCookie(&http.Cookie{
				Name:  k,
				Value: v,
			})
		}
	case map[string]interface{}:
		cookies := r.opts.Cookies.(map[string]interface{})
		for k, v := range cookies {
			r.req.AddCookie(&http.Cookie{
				Name:  k,
				Value: cast.ToString(v),
			})
		}
	case []*http.Cookie:
		cookies := r.opts.Cookies.([]*http.Cookie)
		for _, cookie := range cookies {
			r.req.AddCookie(cookie)
		}
	}
}

func (r *Request) parseHeaders() {
	if r.opts.Headers != nil {
		for k, v := range r.opts.Headers {
			if vv, ok := v.(string); ok {
				r.req.Header.Set(k, vv)
				continue
			}
			if vv, ok := v.([]string); ok {
				for _, vvv := range vv {
					r.req.Header.Add(k, vvv)
				}
			}
			if vv := cast.ToString(v); vv != "" {
				r.req.Header.Set(k, vv)
			}
		}
	}
}

func (r *Request) parseBody() {
	// application/x-www-form-urlencoded
	if r.opts.FormParams != nil {
		if _, ok := r.opts.Headers["Content-Type"]; !ok {
			r.opts.Headers["Content-Type"] = "application/x-www-form-urlencoded"
		}

		values := url.Values{}
		for k, v := range r.opts.FormParams {
			if vv, ok := v.(string); ok {
				values.Set(k, vv)
			}

			if vv, ok := v.([]string); ok {
				for _, vvv := range vv {
					values.Add(k, vvv)
				}
			}

			if vv := cast.ToString(v); vv != "" {
				values.Set(k, vv)
			}
		}
		r.body = strings.NewReader(values.Encode())

		return
	}

	// application/json
	if r.opts.JSON != nil {
		if _, ok := r.opts.Headers["Content-Type"]; !ok {
			r.opts.Headers["Content-Type"] = "application/json"
		}

		b, err := json.Marshal(r.opts.JSON)
		if err == nil {
			r.body = bytes.NewReader(b)
			return
		}
	}

	// application/xml
	if r.opts.XML != nil {
		if _, ok := r.opts.Headers["Content-Type"]; !ok {
			r.opts.Headers["Content-Type"] = "application/xml"
		}

		switch r.opts.XML.(type) {
		case map[string]any:
			mv := mxj.Map(r.opts.XML.(map[string]any))
			b, err := mv.Xml("xml")
			if err == nil {
				r.body = bytes.NewReader(b)
				return
			}
		case map[string]string:
			mv := mxj.Map(r.opts.XML.(map[string]any))
			b, err := mv.Xml("xml")
			if err == nil {
				r.body = bytes.NewReader(b)
				return
			}
		default:
			b, err := xml.Marshal(r.opts.XML)
			if err == nil {
				r.body = bytes.NewReader(b)
				return
			}
		}
	}

	// multipart/form-data
	if r.opts.Multipart != nil {
		if _, ok := r.opts.Headers["Content-Type"]; !ok {
			r.opts.Headers["Content-Type"] = "multipart/form-data"
		}

		buf := new(bytes.Buffer)
		bw := multipart.NewWriter(buf)
		for _, v := range r.opts.Multipart {
			if v.Headers == nil {
				v.Headers = map[string]any{}
			}

			if v.Filepath != "" {
				if v.Filename == "" {
					v.Filename = filepath.Base(v.Filepath)
				}

				if v.Contents == nil {
					f, err := os.Open(v.Filepath)
					if err == nil {
						defer f.Close()

						fi, err := f.Stat()
						if err == nil {
							size := fi.Size()
							fd := make([]byte, size)
							f.Read(fd)
							v.Contents = fd

							if _, ok := v.Headers["Content-Type"]; !ok {
								v.Headers["Content-Type"] = http.DetectContentType(fd)
							}
						}
					}
				}
			}

			arr := []string{
				"form-data",
				fmt.Sprintf("name=%q", v.Name),
			}
			if v.Filename != "" {
				arr = append(arr, fmt.Sprintf("filename=%q", v.Filename))
			}

			h := make(textproto.MIMEHeader)

			h.Set("Content-Disposition", strings.Join(arr, "; "))

			for key, value := range v.Headers {
				if header, ok := value.(string); ok {
					h.Set(key, header)
				}
			}

			p, _ := bw.CreatePart(h)
			io.Copy(p, bytes.NewReader(v.Contents))
		}

		bw.Close()

		r.body = buf
		r.opts.Headers["Content-Type"] = bw.FormDataContentType()
	}
}
