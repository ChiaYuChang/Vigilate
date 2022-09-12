package handlers

import (
	"fmt"
	"log"
	"strconv"
	"time"
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

		for key, val := range repo.App.MonitorMap {
			// remove all items from schedule
			repo.App.Scheduler.Remove(val)
			delete(repo.App.MonitorMap, key)
		}
		for _, e := range repo.App.Scheduler.Entries() {
			repo.App.Scheduler.Remove(e.ID)
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
		return sch, fmt.Errorf("unknown time unit %v", unit)
	}
	return sch, nil
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
			hn := hostName[i]
			s := servicesToMonitor[i]
			log.Printf("|  - Services Name: %s            |\n", hn)

			sch, err := repo.FormScheduleString(s.ScheduleUnit, s.ScheduleNumber)
			if err != nil {
				log.Println(err)
				continue
			}
			// create a job
			j := job{HostServiceID: s.ID}
			scheduleID, err := app.Scheduler.AddJob(sch, j)
			if err != nil {
				log.Println(err)
			}

			// save the id of the job so we can start/stop it
			app.MonitorMap[s.ID] = scheduleID

			// broadcast over websockets the fact that the service is scheduled
			payload := make(map[string]string)
			payload["message"] = "scheduling"
			payload["host_service_id"] = strconv.Itoa(s.ID)
			yearOne := time.Date(0001, 11, 17, 20, 34, 58, 65138737, time.UTC)
			if app.Scheduler.Entry(app.MonitorMap[s.ID]).Next.After(yearOne) {
				payload["next_run"] = app.Scheduler.Entry(app.MonitorMap[s.ID]).Next.Format("2006-01-02 3:04:05 PM")
			} else {
				payload["next_run"] = "Pending..."
			}
			payload["host"] = hn
			payload["service"] = s.Service.ServiceName
			if s.LastCheck.After(yearOne) {
				payload["last_run"] = s.LastCheck.Format("2006-01-02 3:04:05 PM")
			} else {
				payload["last_run"] = "Pending..."
			}
			payload["schedule"] = sch

			err = app.WsClient.Trigger("public-channel", "next-run-event", payload)
			if err != nil {
				log.Println(err)
			}

			err = app.WsClient.Trigger("public-channel", "next-changed-event", payload)
			if err != nil {
				log.Println(err)
			}

		}

	}
}
