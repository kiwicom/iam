package okta

import (
	"log"
	"time"

	"github.com/getsentry/raven-go"
	cfg "github.com/iam/config"
	"github.com/iam/monitoring"
	"github.com/iam/storage"
	"github.com/kiwicom/go-useragent"
	"golang.org/x/sync/singleflight"
)

// Cacher contains methods needed from a cache
type Cacher interface {
	Get(key string, value interface{}) error
	Set(key string, value interface{}, ttl time.Duration) error
	Del(key string) error
	MSet(pairs map[string]interface{}, ttl time.Duration) error
}

// Fetcher is a function used to send HTTP requests.
type Fetcher func(req Request) (*Response, error)

// ClientOpts contains options to create an Okta client
type ClientOpts struct {
	Cache         Cacher
	LockManager   *storage.LockManager
	BaseURL       string
	AuthToken     string
	IAMConfig     *cfg.ServiceConfig
	Metrics       *monitoring.Metrics
	CustomFetcher func(userAgent string, metrics *monitoring.Metrics) Fetcher
}

// Client represent an Okta client
type Client struct {
	group     singleflight.Group
	cache     Cacher
	lock      *storage.LockManager
	baseURL   string
	authToken string
	iamConfig *cfg.ServiceConfig
	metrics   *monitoring.Metrics
	fetch     Fetcher
}

func getUserAgent(iamConfig *cfg.ServiceConfig) (string, error) {
	if iamConfig == nil {
		// We can get to this point only when running tests. If IAM config is not
		// defined there will be a failure much earlier in a real scenario.
		return "", nil
	}

	ua := useragent.UserAgent{
		Name:        "kiwi-iam",
		Environment: iamConfig.Environment,
		Version:     iamConfig.Release,
	}

	uaString, uaErr := ua.Format()
	return uaString, uaErr
}

// NewClient creates an Okta client based on the given options
func NewClient(opts *ClientOpts) *Client {
	uaString, uaErr := getUserAgent(opts.IAMConfig)
	if uaErr != nil {
		log.Println("[ERR]", uaErr)
		raven.CaptureError(uaErr, nil)
	}

	fetch := defaultFetcher(uaString, opts.Metrics)
	if opts.CustomFetcher != nil {
		fetch = opts.CustomFetcher(uaString, opts.Metrics)
	}

	return &Client{
		cache:     opts.Cache,
		lock:      opts.LockManager,
		baseURL:   opts.BaseURL,
		authToken: opts.AuthToken,
		iamConfig: opts.IAMConfig,
		metrics:   opts.Metrics,
		fetch:     fetch,
	}
}
