package handlers

import (
	"log"
	"net/http"
	"sort"

	"github.com/CloudyKit/jet/v6"
	"gitlab.com/gjerry134679/vigilate/pkg/helpers"
	"gitlab.com/gjerry134679/vigilate/pkg/models"
)

type ByHost []models.Schedule

func (a ByHost) Len() int           { return len(a) }
func (a ByHost) Less(i, j int) bool { return a[i].Host < a[j].Host }
func (a ByHost) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// ListEntries lists schedule entries
func (repo *DBRepo) ListEntries(w http.ResponseWriter, r *http.Request) {
	var items []models.Schedule

	for k, v := range repo.App.MonitorMap {
		var item models.Schedule
		item.ID = k
		item.EntryID = v
		item.Entry = repo.App.Scheduler.Entry(v)
		hs, err := repo.DB.GetHostServiceByID(k)
		if err != nil {
			log.Println(err)
			return
		}
		item.ScheduleText, _ = repo.FormScheduleString(hs.ScheduleUnit, hs.ScheduleNumber)
		item.LastRunFromHS = hs.LastCheck

		h, err := repo.DB.GetHostByID(hs.HostID)
		if err != nil {
			log.Println(err)
			return
		}
		item.Host = h.HostName
		item.Service = hs.Service.ServiceName
		items = append(items, item)
	}

	sort.Sort(ByHost(items))
	data := make(jet.VarMap)
	data.Set("items", items)

	err := helpers.RenderPage(w, r, "schedule", data, nil)
	if err != nil {
		printTemplateError(w, err)
	}
}
