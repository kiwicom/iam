package security

import (
	"errors"
	"regexp"

	"github.com/kiwicom/go-useragent"
)

const (
	serviceNamePattern string = `^[\w\s-]+`
)

// It's important to add $ in order to match the whole string
var checkServiceNameRe = regexp.MustCompile(serviceNamePattern + "$")

// CheckServiceName returns if the given service name contains expected characters only
func CheckServiceName(service string) error {
	safe := checkServiceNameRe.MatchString(service)
	if !safe {
		return errors.New("service name has to match " + serviceNamePattern + "$")
	}

	return nil
}

// Service holds the requesting service's name and environment
type Service struct {
	Name        string
	Environment string
}

// GetService returns the service name and environment based on the given user agent
func GetService(incomingUserAgent string) (Service, error) {
	ua, err := useragent.Parse(incomingUserAgent)
	if err != nil {
		return Service{}, err
	}
	return Service{ua.Name, ua.Environment}, nil
}
