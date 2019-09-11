package okta

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	cfg "github.com/kiwicom/iam/configs"
)

// Mock body for HTTP response. It implements the io.ReadWriter interface
// by exposing the Read and Close methods.
type BodyMock struct {
	mock.Mock
	Value string
}

func (b *BodyMock) Read(p []byte) (int, error) {
	copy(p, b.Value)

	return len(b.Value), io.EOF
}

func (b *BodyMock) Close() error {
	b.Called()
	return nil
}

func TestString(t *testing.T) {
	var body = BodyMock{Value: `{ "message": "this is a test" }`}
	body.On("Close").Return()

	res := Response{
		&http.Response{Body: &body},
	}
	expected := `{ "message": "this is a test" }`
	actual, err := res.String()

	assert.NoError(t, err)
	body.AssertNumberOfCalls(t, "Close", 1)
	assert.Equal(t, expected, actual)
}

func TestJSON(t *testing.T) {
	var body = BodyMock{Value: `{ "message": "this is a test" }`}
	body.On("Close").Return()

	type Data struct{ Message string }
	var expectedData Data

	var res = Response{
		&http.Response{Body: &body},
	}
	err := res.JSON(&expectedData)

	if err != nil {
		panic(err)
	}

	body.AssertNumberOfCalls(t, "Close", 1)
	assert.Equal(t, expectedData, Data{Message: "this is a test"})
}

func TestJoinURL(t *testing.T) {
	tests := map[string]struct {
		args     []string
		expected string
	}{
		"no trailing slashes": {
			args:     []string{"http://example.com", "/api", "/path"},
			expected: "http://example.com/api/path",
		},
		"with trailing slashes": {
			args:     []string{"ws://example.com/", "/api/"},
			expected: "ws://example.com/api",
		},
		"no leading slashes": {
			args:     []string{"https://example.com", "api", "path"},
			expected: "https://example.com/api/path",
		},
		"no URL scheme": {
			args:     []string{"example.com", "api", "path"},
			expected: "example.com/api/path",
		},
	}

	for name, test := range tests {
		// scopelint has some issues here. https://github.com/kyoh86/scopelint/issues/4
		test := test

		t.Run(name, func(t *testing.T) {
			result, err := joinURL(test.args[0], test.args[1:]...)
			require.NoError(t, err)
			assert.Equal(t, test.expected, result)
		})
	}
}

func mockHandler(t *testing.T) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "token", r.Header.Get("Authorization"))
		assert.Equal(t, "kiwi-iam/version (Kiwi.com test)", r.Header.Get("User-Agent"))

		_, _ = w.Write([]byte("Okay"))
	})
}

func TestFetch(t *testing.T) {
	ts := httptest.NewServer(mockHandler(t))
	defer ts.Close()

	c := NewClient(&ClientOpts{
		IAMConfig: &cfg.ServiceConfig{
			Environment: "test",
			Release:     "version",
		},
	})

	_, err := c.fetch(Request{
		Method: "GET",
		URL:    ts.URL,
		Token:  "token",
	})

	assert.NoError(t, err, "HTTP request should be sent without errors")
}
