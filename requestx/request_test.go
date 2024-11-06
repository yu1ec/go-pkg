package requestx_test

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/yu1ec/go-pkg/requestx"
)

func ExampleGet() {
	resp, err := requestx.Get("http://127.0.0.1:8091/get")
	if err != nil {
		log.Fatalln(err)
	}

	body, _ := resp.GetBody()
	fmt.Printf("%s", body)
	// Output: http get
}

func ExampleRequest_Get() {
	cli := requestx.NewClient(requestx.Options{
		BaseURI: "http://127.0.0.1:8091",
	})

	resp, err := cli.Get("/get")
	if err != nil {
		log.Fatalln(err)
	}

	body, _ := resp.GetBody()
	fmt.Printf("%s", body)
	// Output: http get
}

func ExampleRequest_Get_withQuery_arr() {
	cli := requestx.NewClient()

	resp, err := cli.Get("http://127.0.0.1:8091/get-with-query", requestx.Options{
		Query: map[string]interface{}{
			"key1": "value1",
			"key2": []string{"value21", "value22"},
			"key3": "abc",
		},
	})
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%s", resp.GetRequest().URL.RawQuery)
	// Output: key1=value1&key2=value21&key2=value22&key3=abc
}

func ExampleRequest_Get_withQuery_str() {
	cli := requestx.NewClient()

	resp, err := cli.Get("http://127.0.0.1:8091/get-with-query?key0=value0", requestx.Options{
		Query: "key1=value1&key2=value21&key2=value22&key3=333",
	})
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%s", resp.GetRequest().URL.RawQuery)
	// Output: key1=value1&key2=value21&key2=value22&key3=333
}

func ExampleRequest_Get_withProxy() {
	cli := requestx.NewClient()

	resp, err := cli.Get("https://www.test.com/test.php", requestx.Options{
		Timeout: 5.0,
		Proxy:   "http://127.0.0.1:1080",
	})
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(resp.GetStatusCode())
	}

	// Output: Get "https://www.test.com/test.php": proxyconnect tcp: dial tcp 127.0.0.1:1080: connect: connection refused
}

func ExampleRequest_Post() {
	cli := requestx.NewClient()

	resp, err := cli.Post("http://127.0.0.1:8091/post")
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%T", resp)
	// Output: *requestx.Response
}

func ExampleRequest_Post_withStreamResponse() {
	cli := requestx.NewClient()

	resp, err := cli.Post("http://127.0.0.1:8091/post-with-stream-response", requestx.Options{
		Headers: map[string]interface{}{
			"Accept": "text/event-stream",
		},
		JSON: map[string]interface{}{
			"foo": "bar",
		},
	})
	if err != nil {
		log.Fatalln(err)
	}

	if !strings.HasPrefix(resp.GetHeaderLine("content-type"), "text/event-stream") {
		body, _ := resp.GetBody()
		log.Fatalf("get stream failed: %s\n", body)
	}

	var message []byte

	for data := range resp.Stream() {
		log.Printf("stream data: %s\n", data)
		message = append(message, data...)
	}

	if err := resp.Err(); err != nil {
		log.Fatalf("stream closed with error: %v\n", err)
	}

	// log.Printf("%s", message)
	fmt.Printf("%s", message)
	// Output: this message will response with stream
}

func ExampleRequest_Post_withHeaders() {
	cli := requestx.NewClient()

	resp, err := cli.Post("http://127.0.0.1:8091/post-with-headers", requestx.Options{
		Headers: map[string]interface{}{
			"User-Agent": "testing/1.0",
			"Accept":     "application/json",
			"X-Foo":      []string{"Bar", "Baz"},
		},
	})
	if err != nil {
		log.Fatalln(err)
	}

	headers := resp.GetRequest().Header["X-Foo"]
	fmt.Println(headers)
	// Output: [Bar Baz]
}

func ExampleRequest_Post_withCookies_str() {
	cli := requestx.NewClient()

	resp, err := cli.Post("http://127.0.0.1:8091/post-with-cookies", requestx.Options{
		Cookies: "cookie1=value1;cookie2=value2",
	})
	if err != nil {
		log.Fatalln(err)
	}

	body, _ := resp.GetBody()
	fmt.Printf("%T", body)
	// Output: requestx.ResponseBody
}

func ExampleRequest_Post_withCookies_map() {
	cli := requestx.NewClient()

	resp, err := cli.Post("http://127.0.0.1:8091/post-with-cookies", requestx.Options{
		Cookies: map[string]interface{}{
			"cookie1": "value1",
			"cookie2": "value2",
		},
	})
	if err != nil {
		log.Fatalln(err)
	}

	body, _ := resp.GetBody()
	fmt.Printf("%T", body)
	// Output: requestx.ResponseBody
}

func ExampleRequest_Post_withCookies_obj() {
	cli := requestx.NewClient()

	cookies := make([]*http.Cookie, 0, 2)
	cookies = append(cookies, &http.Cookie{
		Name:     "cookie133",
		Value:    "value1",
		Domain:   "httpbin.org",
		Path:     "/cookies",
		HttpOnly: true,
	})
	cookies = append(cookies, &http.Cookie{
		Name:   "cookie2",
		Value:  "value2",
		Domain: "httpbin.org",
		Path:   "/cookies",
	})

	resp, err := cli.Post("http://127.0.0.1:8091/post-with-cookies", requestx.Options{
		Cookies: cookies,
	})
	if err != nil {
		log.Fatalln(err)
	}

	body, _ := resp.GetBody()
	fmt.Printf("%T", body)
	// Output: requestx.ResponseBody
}

func ExampleRequest_Post_withFormParams() {
	cli := requestx.NewClient()

	resp, err := cli.Post("http://127.0.0.1:8091/post-with-form-params", requestx.Options{
		Headers: map[string]interface{}{
			"Content-Type": "application/x-www-form-urlencoded",
		},
		FormParams: map[string]interface{}{
			"key1": "value1",
			"key2": []string{"value21", "value22"},
			"key3": "333",
		},
	})
	if err != nil {
		log.Fatalln(err)
	}

	body, _ := resp.GetBody()
	fmt.Println(body)
	// Output: form params:{"key1":["value1"],"key2":["value21","value22"],"key3":["333"]}
}

func ExampleRequest_Post_withJSON() {
	cli := requestx.NewClient()

	resp, err := cli.Post("http://127.0.0.1:8091/post-with-json", requestx.Options{
		Headers: map[string]any{
			"Content-Type": "application/json",
		},
		JSON: struct {
			Key1 string   `json:"key1"`
			Key2 []string `json:"key2"`
			Key3 int      `json:"key3"`
		}{"value1", []string{"value21", "value22"}, 333},
	})
	if err != nil {
		log.Fatalln(err)
	}

	body, _ := resp.GetBody()
	fmt.Println(body)
	// Output: json:{"key1":"value1","key2":["value21","value22"],"key3":333}
}

func ExampleRequest_Post_withXML() {
	cli := requestx.NewClient()

	resp, err := cli.Post("http://127.0.0.1:8091/post-with-xml", requestx.Options{
		XML: map[string]interface{}{
			"out_trade_no": "xxx",
			"total_fee":    333,
		},
	})
	if err != nil {
		log.Fatalln(err)
	}

	body, _ := resp.GetBody()
	fmt.Println(string(body.Read(9)))
	// Output: xml:<xml>
}

func ExampleRequest_Post_withMultipart() {
	cli := requestx.NewClient(requestx.Options{
		Debug: false,
	})

	resp, err := cli.Post("http://127.0.0.1:8091/post-with-multipart", requestx.Options{
		Multipart: []requestx.FormData{
			{
				Name:     "foo",
				Contents: []byte("bar"),
			},
			{
				Name:     "json",
				Contents: []byte(`{"title":"title","intro":"introduction"}`),
			},
			{
				Name:     "media",
				Filepath: "./image.png",
			},
		},
	})
	if err != nil {
		log.Fatalln(err)
	}

	body, _ := resp.GetBody()
	fmt.Println(string(body.Read(10)))
	// Output: body:
}

func ExampleRequest_Put() {
	cli := requestx.NewClient()

	resp, err := cli.Put("http://127.0.0.1:8091/put")
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%T", resp)
	// Output: *requestx.Response
}

func ExampleRequest_Patch() {
	cli := requestx.NewClient()

	resp, err := cli.Patch("http://127.0.0.1:8091/patch")
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%T", resp)
	// Output: *requestx.Response
}

func ExampleRequest_Delete() {
	cli := requestx.NewClient()

	resp, err := cli.Delete("http://127.0.0.1:8091/delete")
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%T", resp)
	// Output: *requestx.Response
}

func ExampleRequest_Options() {
	cli := requestx.NewClient()

	resp, err := cli.Options("http://127.0.0.1:8091/options")
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%T", resp)
	// Output: *requestx.Response
}
