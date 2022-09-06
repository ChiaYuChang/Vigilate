package main

import (
	"log"
)

type job struct {
	HostServiceID int
}

func (j job) Run() {
	repo.ScheduleCheck(j.HostServiceID)
}

func startMonitoring() {
	log.Println("====== Start Monitoring Servicse !! ======")
	if preferenceMap["monitoring_live"] == "1" {
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
		for _, hn := range hostName {
			log.Printf("|  - Services Name: %s            |\n", hn)
		}
		// range through the services

		// get the schedule unit and number

		// create a job

		// save the id of the job so we can start/stop it

		// broadcast over websockets the fact that the service is scheduled
	}
}
