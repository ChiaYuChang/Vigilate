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

var YearOne time.Time = time.Date(0001, 1, 1, 0, 0, 0, 1, time.UTC)

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
func (repo *DBRepo) ScheduleCheck(HostServiceId int) {
	log.Println("Running check for host service id:", HostServiceId)
	hs, err := repo.DB.GetHostServiceByID(HostServiceId)
	if err != nil {
		log.Println(err)
		return
	}

	h, err := repo.DB.GetHostByID(hs.HostID)
	if err != nil {
		log.Println(err)
		return
	}

	s := ParseService(hs.Service.ServiceName)
	log.Println("Host:", h.URL)
	log.Println("Service name is", s.String())

	var newServiceStatus models.ServiceStatus
	var message string
	var updateTime time.Time

	switch s {
	case ServiceHTTP:
		newServiceStatus, message, updateTime = testHTTPForHost(h.URL)
		// newServiceStatus, message, updateTime = repo.testServiceForHost(h, hs)
	default:
		log.Println("Currently not support.")
		return
	}
	log.Printf("New Status: %s; Message: %s", newServiceStatus, message)

	// update host service record in db with status and update the last check time
	serviceStatusHasChange := hs.Status != newServiceStatus
	hs.Status = newServiceStatus
	hs.LastCheck = updateTime
	message = fmt.Sprintf(
		"host service %s on %s has change to %s",
		hs.Service.ServiceName, h.HostName, newServiceStatus.String(),
	)

	err = repo.DB.UpdateHostService(hs)
	if err != nil {
		log.Println(err)
		return
	}

	// broadcast the update info
	if serviceStatusHasChange {
		repo.broadcastMessage(
			"public-channel",
			"host-service-count-change",
			repo.updateHostServiceCount(message),
		)
	}
}

func (repo *DBRepo) updateHostServiceCount(message string) map[string]string {
	serviceStateCount, err := repo.DB.GetAllServiceStatusCounts()
	data := make(map[string]string)
	if err != nil {
		log.Println(err)
		return data
	}

	data["message"] = message
	for k, v := range serviceStateCount {
		data[k.String()+"_count"] = strconv.Itoa(v)
	}
	return data
}

func (repo *DBRepo) broadcastMessage(channel, messageType string, data map[string]string) {
	err := repo.App.WsClient.Trigger(channel, messageType, data)
	if err != nil {
		log.Println(err)
	}
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
	newServiceStatus, msg, checkedTime := repo.testServiceForHost(h, hs)

	// Broadcast service change event
	if hs.Status != newServiceStatus {
		repo.pushScheduleChangeEvent(hs.ID, hs.HostID, h.HostName, hs.ServiceID, hs.Service.ServiceName,
			hs.Service.Icon, hs.ScheduleUnit, hs.ScheduleNumber, checkedTime, newServiceStatus)
		repo.broadcastMessage(
			"public-channel",
			"host-service-count-change",
			repo.updateHostServiceCount("from TestCheck function"),
		)
	}

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
			NewStatus:     newServiceStatus.String(),
			LastCheck:     time.Now(),
		}
	} else {
		resp.OK = ok
		resp.Message = "service internal error"
	}

	// update the host service in the database
	hs.Status = newServiceStatus
	hs.UpdatedAt = checkedTime
	hs.LastCheck = checkedTime
	err = repo.DB.UpdateHostService(hs)
	if err != nil {
		log.Println(err)
		ok = false
	}

	// send json to client
	out, _ := json.MarshalIndent(resp, "", "    ")
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (repo *DBRepo) testServiceForHost(h models.Host, hs models.HostService) (models.ServiceStatus, string, time.Time) {
	var newServiceStatus models.ServiceStatus
	var msg string
	var t time.Time

	log.Printf("Start testing %s on %s", hs.Service.ServiceName, h.URL)
	switch ParseService(hs.Service.ServiceName) {
	case ServiceHTTP:
		newServiceStatus, msg, t = testHTTPForHost(h.URL)
	}
	log.Printf("Old status: %s; New status: %s", hs.Status, newServiceStatus)

	// broadcast to all clients
	if hs.Status != newServiceStatus {
		repo.pushServerStatusChangeEvent(hs.ID, hs.HostID, h.HostName, hs.ServiceID, hs.Service.ServiceName,
			hs.Service.Icon, hs.ScheduleUnit, hs.ScheduleNumber, t, hs.Status, newServiceStatus)
		repo.broadcastMessage(
			"public-channel",
			"host-service-count-change",
			repo.updateHostServiceCount("from testServiceForHost function"),
		)
	}

	// Broadcast schedule-change-event
	repo.pushScheduleChangeEvent(hs.ID, hs.HostID, h.HostName, hs.ServiceID, hs.Service.ServiceName,
		hs.Service.Icon, hs.ScheduleUnit, hs.ScheduleNumber, t, newServiceStatus)

	// TODO - send email if necessary

	return newServiceStatus, msg, t
}

func (repo *DBRepo) pushServerStatusChangeEvent(hostServiceID int, hostID int, hostName string,
	serviceID int, serviceName string, serviceIcon string, scheduleUnit string, scheduleNumber int,
	lastCheckTime time.Time, oldServiceStatus models.ServiceStatus, newServiceStatus models.ServiceStatus) {
	data := make(map[string]string)
	data["host_service_id"] = strconv.Itoa(hostServiceID)
	data["host_id"] = strconv.Itoa(hostID)
	data["host_name"] = hostName
	data["service_id"] = strconv.Itoa(serviceID)
	data["service_name"] = serviceName
	data["icon"] = serviceIcon
	data["status"] = newServiceStatus.String()
	data["old_status"] = oldServiceStatus.String()

	data["message"] = fmt.Sprintf(
		"%s on %s reports %s",
		serviceName, hostName, newServiceStatus.String(),
	)

	data["last_check"] = lastCheckTime.Format("2006-01-02 03:04:05 PM")
	repo.broadcastMessage("public-channel", "host-service-status-change", data)
}

func (repo *DBRepo) pushScheduleChangeEvent(hostServiceID int, hostID int, hostName string,
	serviceID int, serviceName string, serviceIcon string, scheduleUnit string, scheduleNumber int,
	lastCheckTime time.Time, newServiceStatus models.ServiceStatus) {
	data := make(map[string]string)
	data["host_service_id"] = strconv.Itoa(hostServiceID)
	data["host_id"] = strconv.Itoa(hostID)
	data["host_name"] = hostName
	data["service_id"] = strconv.Itoa(serviceID)
	data["service_name"] = serviceName
	data["schedule"], _ = repo.FormScheduleString(scheduleUnit, scheduleNumber)
	if app.Scheduler.Entry(repo.App.MonitorMap[hostServiceID]).Next.After(YearOne) {
		data["next_run"] = app.Scheduler.
			Entry(repo.App.MonitorMap[hostServiceID]).
			Next.
			Format("2006-01-02 03:04:05 PM")
	} else {
		data["next_run"] = "Pending"
	}
	data["last_run"] = lastCheckTime.Format("2006-01-02 03:04:05 PM")
	data["status"] = newServiceStatus.String()
	data["icon"] = serviceIcon
	repo.broadcastMessage("public-channel", "schedule-changed-event", data)
}

func testHTTPForHost(url string) (models.ServiceStatus, string, time.Time) {
	url = strings.TrimSuffix(url, "/")

	url = strings.Replace(url, "https://", "http://", -1)
	t := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		return models.ServiceStatusProblem, fmt.Sprintf("%s - %s", url, "error connecting"), t
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.ServiceStatusProblem, fmt.Sprintf("%s - %s", url, resp.Status), t
	}

	return models.ServiceStatusHealthy, fmt.Sprintf("%s - %s", url, resp.Status), t
}

// type ServiceTester struct {
// 	repo       *DBRepo
// 	strategies map[Service]Strategy
// }

// func (st *ServiceTester) TestHandleFunc(w http.ResponseWriter, r *http.Request) {
// 	ok := true
// 	hostServiceId, _ := strconv.Atoi(chi.URLParam(r, "id"))
// 	oldServiceStatus := chi.URLParam(r, "oldStatus")

// 	newServiceStatus, message, tm, err := st.ScheduleCheck(hostServiceId)
// }

// func (st *ServiceTester) ScheduleCheck(hostServiceId int) (models.ServiceStatus, string, time.Time, err) {
// 	log.Println("Running check for host service id:", hostServiceId)
// 	hs, err := st.repo.DB.GetHostServiceByID(hostServiceId)
// 	if err != nil {
// 		return models.ServiceStatusUnknown, "", time.Now(), err
// 	}

// 	h, err := st.repo.DB.GetHostByID(hs.HostID)
// 	if err != nil {
// 		return models.ServiceStatusUnknown, "", time.Now(), err
// 	}

// 	s := ParseService(hs.Service.ServiceName)
// 	return st.strategies[s].Test(h.URL)
// }

// func (st *ServiceTester) Test(service Service, URL string) (models.ServiceStatus, string, time.Time) {
// 	return st.strategies[service].Test(service, URL)
// }

// func NewServiceTester(repo *DBRepo) *ServiceTester {
// 	return &ServiceTester{
// 		repo: repo,
// 		strategies: map[Service]Strategy{
// 			ServiceHTTP: &HTTPServiceTester{},
// 		},
// 	}
// }

// type Strategy interface {
// 	Test(URL string) (models.ServiceStatus, string, time.Time, error)
// }

// type HTTPServiceTester struct{}

// func (st *HTTPServiceTester) Test(s Service, URL string) (models.ServiceStatus, string, time.Time) {
// 	if s == ServiceHTTP {
// 		if strings.HasSuffix(URL, "/") {
// 			URL = strings.TrimSuffix(URL, "/")
// 		}

// 		URL = strings.Replace(URL, "https://", "http://", -1)
// 		t := time.Now()
// 		resp, err := http.Get(URL)
// 		if err != nil {
// 			return models.ServiceStatusProblem, fmt.Sprintf("%s - %s", URL, "error connecting"), t
// 		}
// 		defer resp.Body.Close()

// 		if resp.StatusCode != http.StatusOK {
// 			return models.ServiceStatusProblem, fmt.Sprintf("%s - %s", URL, resp.Status), t
// 		}

// 		return models.ServiceStatusHealthy, fmt.Sprintf("%s - %s", URL, resp.Status), t
// 	}
// 	return models.ServiceStatusUnknown, "error service type", time.Now()
// }
