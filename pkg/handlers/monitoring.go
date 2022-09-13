package handlers

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"gitlab.com/gjerry134679/vigilate/pkg/models"
)

type job struct {
	HostServiceID int
}

func (j job) Run() {
	Repo.ScheduleCheck(j.HostServiceID)
}

func (repo *DBRepo) StopMonitoring() {
	log.Println("====== Stop Monitoring Servicse !!! ======")
	if app.PreferenceMap["monitoring_live"] == "0" {
		data := make(map[string]string)
		data["message"] = "Monitoring is off"

		// trigger a message to broadcast ot all clients that app is starting to monitor
		err := app.WsClient.Trigger("public-channel", "app-ending", data)
		if err != nil {
			log.Println(err)
		}

		for key := range repo.App.MonitorMap {
			// remove all items from schedule
			_, _ = repo.RemoveFromMonitorMap(key)
			// repo.App.Scheduler.Remove(val)
			// delete(repo.App.MonitorMap, key)
		}
		for _, e := range repo.App.Scheduler.Entries() {
			repo.App.Scheduler.Remove(e.ID)
		}
	}
}

func (repo *DBRepo) StartMonitoring() {
	log.Println("====== Start Monitoring Servicse !! ======")
	if app.PreferenceMap["monitoring_live"] == "1" {
		data := make(map[string]string)
		data["message"] = "Monitoring is starting..."

		// trigger a message to broadcast ot all clients that app is starting to monitor
		err := app.WsClient.Trigger("public-channel", "app-starting", data)
		if err != nil {
			log.Println(err)
		}

		// get all of the services that we want to monitor
		servicesToMonitor, hostName, err := repo.DB.GetServivesToMonitor()
		if err != nil {
			log.Println(err)
		}
		log.Printf("| Length of servicesToMonitor is: %d      |\n", len(servicesToMonitor))

		// range through the services
		for i := 0; i < len(hostName); i++ {
			payload, err := repo.AddtoMonitorMap(servicesToMonitor[i], hostName[i])
			if err != nil {
				log.Println(err)
				continue
			}
			repo.broadcastMessage("public-channel", "next-run-event", payload)
			repo.broadcastMessage("public-channel", "next-changed-event", payload)
		}

	}
}

func (repo *DBRepo) FormScheduleString(unit string, number int) (string, error) {
	// format schedule string
	var sch string
	switch unit {
	case "s":
		// should not be used in production
		sch = fmt.Sprintf("@every %ds", number)
	case "d":
		sch = fmt.Sprintf("@every %dh", number*24)
	case "h":
		sch = fmt.Sprintf("@every %dh", number)
	case "m":
		hr := number / 60
		mn := number % 60
		if hr < 1 {
			sch = fmt.Sprintf("@every %dm", mn)
		} else {
			sch = fmt.Sprintf("@every %dh%dm", hr, mn)
		}
	default:
		// log.Println("unknown time unit:", unit)
		if len(unit) < 1 {
			return sch, errors.New("time unit is empty")
		}
		return sch, fmt.Errorf("unknown time unit %v", unit)
	}
	return sch, nil
}

func (repo *DBRepo) AddtoMonitorMap(hs models.HostService, hn string) (map[string]string, error) {
	payload := make(map[string]string)
	if repo.App.PreferenceMap["monitoring_live"] == "1" {
		// log.Printf("|  - Services Name: %s            |\n", hn)
		sch, err := repo.FormScheduleString(hs.ScheduleUnit, hs.ScheduleNumber)
		if err != nil {
			return payload, err
		}

		// create a job
		j := job{HostServiceID: hs.ID}
		scheduleID, err := repo.App.Scheduler.AddJob(sch, j)
		if err != nil {
			return payload, err
		}

		// save the id of the job so we can start/stop it
		app.MonitorMap[hs.ID] = scheduleID

		// broadcast over websockets the fact that the service is scheduled
		payload["message"] = "scheduling"
		payload["host_service_id"] = strconv.Itoa(hs.ID)

		if app.Scheduler.Entry(app.MonitorMap[hs.ID]).Next.After(YearOne) {
			payload["next_run"] = app.Scheduler.Entry(app.MonitorMap[hs.ID]).Next.Format("2006-01-02 3:04:05 PM")
		} else {
			payload["next_run"] = "Pending..."
		}
		payload["host"] = hn
		payload["service"] = hs.Service.ServiceName
		if hs.LastCheck.After(YearOne) {
			payload["last_run"] = hs.LastCheck.Format("2006-01-02 3:04:05 PM")
		} else {
			payload["last_run"] = "Pending..."
		}
		payload["schedule"] = sch

		return payload, nil
	}
	return payload, fmt.Errorf("monitoring_live equals to %q", repo.App.PreferenceMap["monitoring_live"])
}

func (repo *DBRepo) RemoveFromMonitorMap(hostServiceID int) (map[string]string, error) {
	payload := make(map[string]string)
	if repo.App.PreferenceMap["monitoring_live"] == "1" {
		entryID := repo.App.MonitorMap[hostServiceID]
		repo.App.Scheduler.Remove(entryID)
		delete(repo.App.MonitorMap, hostServiceID)
		payload["host_service_id"] = strconv.Itoa(hostServiceID)
		return payload, nil
	}
	return payload, fmt.Errorf("monitoring_live equals to %q", repo.App.PreferenceMap["monitoring_live"])
}
