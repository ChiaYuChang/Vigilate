package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"gitlab.com/gjerry134679/vigilate/pkg/models"
)

type Service int

const (
	ServiceHTTP Service = iota
	ServiceHTTPS
	ServiceSSLCertificate
	ServiceUnknown
)

func (s Service) String() string {
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

func ParseService(s string) Service {
	var service Service
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

type jsonResp struct {
	OK            bool      `json:"ok"`
	Message       string    `json:"message"`
	ServiceID     int       `json:"service_id"`
	HostServiceID int       `json:"host_service_id"`
	HostID        int       `json:"host_id"`
	OldStatus     string    `json:"old_status"`
	NewStatus     string    `json:"new_status"`
	LastCheck     time.Time `jsom:"last_check"`
}

// Preform a schedule a schedule check on a host server by id
func (repo *DBRepo) ScheduleCheck(HostServiceID int) {

}

func (repo *DBRepo) TestCheck(w http.ResponseWriter, r *http.Request) {
	ok := true
	hostServiceId, _ := strconv.Atoi(chi.URLParam(r, "id"))
	oldStatus := chi.URLParam(r, "oldStatus")

	log.Println(hostServiceId, oldStatus)

	// get host service
	hs, err := repo.DB.GetHostServiceByID(hostServiceId)
	if err != nil {
		log.Println(err)
		ok = false
	}
	log.Println("Service name is", hs.Service.ServiceName)

	// get host
	h, err := repo.DB.GetHostByID(hs.HostID)
	if err != nil {
		log.Println(err)
		ok = false
	}

	// test the service
	newStatus, msg := repo.testServiceForHost(h, hs)

	// create json
	var resp jsonResp
	if ok {
		resp = jsonResp{
			OK:            ok,
			Message:       msg,
			ServiceID:     hs.ServiceID,
			HostServiceID: hs.ID,
			HostID:        hs.HostID,
			OldStatus:     oldStatus,
			NewStatus:     newStatus.String(),
			LastCheck:     time.Now(),
		}
	} else {
		resp.OK = ok
		resp.Message = "service internal error"
	}

	// update the host service in the database
	hs.Status = newStatus
	hs.UpdatedAt = time.Now()
	hs.LastCheck = time.Now()
	err = repo.DB.UpdateHostService(hs)
	if err != nil {
		log.Println(err)
		ok = false
	}
	// Broadcast service change event

	// send json to client
	out, _ := json.MarshalIndent(resp, "", "    ")
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (repo *DBRepo) testServiceForHost(h models.Host, hs models.HostService) (models.ServiceStatus, string) {
	var newStatus models.ServiceStatus
	var msg string

	switch ParseService(hs.Service.ServiceName) {
	case ServiceHTTP:
		newStatus, msg = testHTTPForHost(h.URL)
	}

	return newStatus, msg
}

func testHTTPForHost(url string) (models.ServiceStatus, string) {
	if strings.HasSuffix(url, "/") {
		url = strings.TrimSuffix(url, "/")
	}

	url = strings.Replace(url, "https://", "http://", -1)
	resp, err := http.Get(url)
	if err != nil {
		return models.ServiceStatusProblem, fmt.Sprintf("%s - %s", url, "error connecting")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.ServiceStatusProblem, fmt.Sprintf("%s - %s", url, resp.Status)
	}

	return models.ServiceStatusHealthy, fmt.Sprintf("%s - %s", url, resp.Status)
}
