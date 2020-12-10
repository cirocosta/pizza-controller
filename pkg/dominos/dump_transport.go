package dominos

import (
	"fmt"
	"net/http"
	"net/http/httputil"
)

type DumpTransport struct {
	r http.RoundTripper
}

func (d *DumpTransport) RoundTrip(h *http.Request) (*http.Response, error) {
	dump, _ := httputil.DumpRequestOut(h, true)
	fmt.Println(string(dump))

	resp, err := d.r.RoundTrip(h)
	fmt.Println(resp.StatusCode)

	dump, _ = httputil.DumpResponse(resp, true)
	fmt.Println(string(dump))

	return resp, err
}
