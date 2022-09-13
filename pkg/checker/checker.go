package checker

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"gitlab.com/gjerry134679/vigilate/pkg/models"
)

type ServiceType int

const (
	ServiceHTTP ServiceType = iota
	ServiceHTTPS
	ServiceSSLCertificate
	ServiceUnknown
)

func (s ServiceType) String() string {
	switch s {
	case ServiceHTTP:
		return "HTTP"
	case ServiceHTTPS:
		return "HTTPS"
	case ServiceSSLCertificate:
		return "SSLCertificate"
	default:
		return "Unknown Service"
	}
}

func ParseService(s string) ServiceType {
	var service ServiceType
	switch s {
	case ServiceHTTP.String():
		service = ServiceHTTP
	case ServiceHTTPS.String():
		service = ServiceHTTPS
	case ServiceSSLCertificate.String():
		service = ServiceSSLCertificate
	default:
		service = ServiceUnknown
	}
	return service
}

type ServerChecker struct {
	CheckerCollection map[ServiceType]Checker
}

func NewEmptyServerChecker() *ServerChecker {
	return &ServerChecker{CheckerCollection: map[ServiceType]Checker{
		ServiceUnknown: &EmptyChecker{},
	}}
}

func NewDefaultServerChecker() *ServerChecker {
	sc := NewEmptyServerChecker()
	sc.AppendChecker(ServiceHTTP, &HTTPServiceChecker{})
	return sc
}

func (st *ServerChecker) CheckerSelector(url string, service ServiceType, args map[string]any) (models.ServiceStatus, string, time.Time) {
	checker, ok := st.CheckerCollection[service]
	if ok {
		return checker.Check(url, args)
	}
	return st.CheckerCollection[ServiceUnknown].Check(url, args)
}

func (st *ServerChecker) AppendChecker(stype ServiceType, checker Checker) {
	st.CheckerCollection[stype] = checker
}

type Checker interface {
	Check(url string, args map[string]any) (models.ServiceStatus, string, time.Time)
}

type EmptyChecker struct{}

func (e *EmptyChecker) Check(url string, args map[string]any) (models.ServiceStatus, string, time.Time) {
	return models.ServiceStatusUnknown, "service not support", time.Now()
}

type HTTPServiceChecker struct{}

func (t *HTTPServiceChecker) Check(url string, args map[string]any) (models.ServiceStatus, string, time.Time) {
	url = strings.TrimSuffix(url, "/")

	url = strings.Replace(url, "https://", "http://", -1)
	checkTime := time.Now()

	resp, err := http.Get(url)
	if err != nil {
		return models.ServiceStatusProblem, fmt.Sprintf("%s - %s", url, "error connecting"), checkTime
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.ServiceStatusProblem, fmt.Sprintf("%s - %s", url, resp.Status), checkTime
	}

	return models.ServiceStatusHealthy, fmt.Sprintf("%s - %s", url, resp.Status), checkTime
}

type HTTPSServiceChecker struct {
	EmptyChecker
}

type SSLCertificateServiceChecker struct {
	EmptyChecker
}
