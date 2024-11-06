package requestx_test

import (
	"fmt"

	"github.com/yu1ec/go-pkg/requestx"
)

func ExampleNewClient() {
	cli := requestx.NewClient()

	fmt.Printf("%T", cli)
	// Output: *requestx.Request
}
