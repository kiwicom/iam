package okta

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"time"

	jsoniter "github.com/json-iterator/go"

	"github.com/kiwicom/iam/internal/monitoring"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

var httpClient = &http.Client{
	Timeout: time.Second * 10,
}

// Request contains options for an HTTP request created using shared.Fetch
type Request struct {
	Method string
	URL    string
	Token  string
	Body   io.Reader
}

// Response for an HTTP request, exposes a JSON method to get the
// retrieved data
type Response struct {
	*http.Response
}

// String reads data from HTTP response and returns it in the form of a string.
func (res Response) String() (string, error) {
	defer res.Body.Close()

	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(res.Body)
	return buf.String(), err
}

// JSON retrieves data from HTTP response and store it in the struct pointed by
// `body` (note: `body` should be a pointer to the struct you expect the HTTP
// call to return)
func (res Response) JSON(body interface{}) error {
	// Right before returning, close the stream used for reading the
	// response's body.
	defer res.Body.Close()

	return json.NewDecoder(res.Body).Decode(&body)
}

// defaultFetcher returns a function to send HTTP requests response, this
// wrapper is used to set the userAgent, and metrics client package wide.
func defaultFetcher(userAgent string, metrics *monitoring.Metrics) func(req Request) (*Response, error) {
	return func(req Request) (*Response, error) {
		log.Println(req.Method, req.URL)

		httpReq, err := http.NewRequest(req.Method, req.URL, req.Body)
		if err != nil {
			return nil, err
		}

		httpReq.Header.Set("User-Agent", userAgent)
		httpReq.Header.Set("Authorization", req.Token)

		metrics.Incr("outgoing.requests", monitoring.Tag("url", req.URL))
		httpRes, err := httpClient.Do(httpReq) //nolint:bodyclose // body is closed on reading either to string or JSON
		if err != nil {
			return nil, err
		}
		return &Response{httpRes}, nil
	}
}

// joinURL parses and joins a base URL to a path safely
func joinURL(baseURL string, pathname ...string) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	// prepend u.Path to pathname slice
	elems := append([]string{u.Path}, pathname...)
	u.Path = path.Join(elems...)
	return u.String(), nil
}
