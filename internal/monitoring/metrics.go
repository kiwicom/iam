package monitoring

import (
	"fmt"
	"log"
	"net"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/pkg/errors"
)

// MetricSettings are settings required to set up metrics client
type MetricSettings struct {
	Host        string
	Port        string
	Namespace   string
	Environment string
}

// Metrics is the instance of metrics to be used
type Metrics struct {
	client *statsd.Client
	rate   float64
}

// CreateNewMetricService creates an instance of Metrics and returns it
func CreateNewMetricService(settings MetricSettings) (*Metrics, error) {
	if settings.Host == "" || settings.Port == "" {
		log.Println("Datadog metrics disabled.")
		return &Metrics{rate: 1.0}, nil
	}

	address := net.JoinHostPort(settings.Host, settings.Port)
	dataDogClient, ddErr := statsd.New(
		address,
		statsd.WithNamespace(settings.Namespace),
		statsd.WithTags([]string{fmt.Sprintf("env:%v", settings.Environment)}),
	)

	if ddErr != nil {
		ddErr = errors.Wrap(ddErr, "couldn't initialize Datadog client")
		return &Metrics{}, ddErr
	}

	return &Metrics{client: dataDogClient, rate: 1.0}, nil
}

// Incr increases a metric by 1
func (c *Metrics) Incr(name string, tags ...string) {
	if c == nil || c.client == nil {
		return
	}
	err := c.client.Incr(name, tags, c.rate)
	if err != nil {
		log.Printf("[ERROR] Metric Incr failed: %v", err)
	}
}
