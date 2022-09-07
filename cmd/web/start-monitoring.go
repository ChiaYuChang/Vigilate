package main

// import (
// 	"fmt"
// 	"log"
// 	"strconv"
// 	"time"
// )

// type job struct {
// 	HostServiceID int
// }

// func (j job) Run() {
// 	repo.ScheduleCheck(j.HostServiceID)
// }

// func startMonitoring() {
// 	log.Println("====== Start Monitoring Servicse !! ======")
// 	if preferenceMap["monitoring_live"] == "1" {
// 		data := make(map[string]string)
// 		data["message"] = "Monitoring is starting..."

// 		// trigger a message to broadcast ot all clients that app is starting to monitor
// 		err := app.WsClient.Trigger("public-channel", "app-starting", data)
// 		if err != nil {
// 			log.Println(err)
// 		}

// 		// get all of the services that we want to monitor
// 		servicesToMonitor, hostName, err := repo.DB.GetServivesToMonitor()
// 		if err != nil {
// 			log.Println(err)
// 		}
// 		log.Printf("| Length of servicesToMonitor is: %d      |\n", len(servicesToMonitor))

// 		// range through the services
// 		for i := 0; i < len(hostName); i++ {
// 			hn := hostName[i]
// 			s := servicesToMonitor[i]
// 			log.Printf("|  - Services Name: %s            |\n", hn)

// 			// format schedule string
// 			var sch string
// 			switch s.ScheduleUnit {
// 			case "d":
// 				sch = fmt.Sprintf("@every %dh", s.ScheduleNumber*24)
// 			case "h":
// 				sch = fmt.Sprintf("@every %dh", s.ScheduleNumber)
// 			case "m":
// 				hr := s.ScheduleNumber / 60
// 				mn := s.ScheduleNumber % 60
// 				if hr < 1 {
// 					sch = fmt.Sprintf("@every %dm", mn)
// 				} else {
// 					sch = fmt.Sprintf("@every %dh%dm", hr, mn)
// 				}
// 			default:
// 				log.Println("unknown time unit:", s.ScheduleUnit)
// 				continue
// 			}

// 			// create a job
// 			j := job{HostServiceID: s.ID}
// 			scheduleID, err := app.Scheduler.AddJob(sch, j)
// 			if err != nil {
// 				log.Println(err)
// 			}

// 			// save the id of the job so we can start/stop it
// 			app.MonitorMap[s.ID] = scheduleID

// 			// broadcast over websockets the fact that the service is scheduled
// 			payload := make(map[string]string)
// 			payload["message"] = "scheduling"
// 			payload["host_service_id"] = strconv.Itoa(s.ID)
// 			yearOne := time.Date(0001, 11, 17, 20, 34, 58, 65138737, time.UTC)
// 			if app.Scheduler.Entry(app.MonitorMap[s.ID]).Next.After(yearOne) {
// 				payload["next_run"] = app.Scheduler.Entry(app.MonitorMap[s.ID]).Next.Format("2006-01-02 3:04:05 PM")
// 			} else {
// 				payload["next_run"] = "Pending..."
// 			}
// 			payload["host"] = hn
// 			payload["service"] = s.Service.ServiceName
// 			if s.LastCheck.After(yearOne) {
// 				payload["last_run"] = s.LastCheck.Format("2006-01-02 3:04:05 PM")
// 			} else {
// 				payload["last_run"] = "Pending..."
// 			}
// 			payload["schedule"] = sch

// 			err = app.WsClient.Trigger("public-channel", "next-run-event", payload)
// 			if err != nil {
// 				log.Println(err)
// 			}

// 			err = app.WsClient.Trigger("public-channel", "next-changed-event", payload)
// 			if err != nil {
// 				log.Println(err)
// 			}

// 		}

// 	}
// }
