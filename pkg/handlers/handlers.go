package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/CloudyKit/jet/v6"
	"github.com/go-chi/chi/v5"
	"gitlab.com/gjerry134679/vigilate/pkg/checker"
	"gitlab.com/gjerry134679/vigilate/pkg/config"
	"gitlab.com/gjerry134679/vigilate/pkg/driver"
	"gitlab.com/gjerry134679/vigilate/pkg/helpers"
	"gitlab.com/gjerry134679/vigilate/pkg/models"
	"gitlab.com/gjerry134679/vigilate/pkg/repository"
	"gitlab.com/gjerry134679/vigilate/pkg/repository/dbrepo"
)

//Repo is the repository
var Repo *DBRepo
var app *config.AppConfig
var serverChecker *checker.ServerChecker

// DBRepo is the db repo
type DBRepo struct {
	App     *config.AppConfig
	DB      repository.DatabaseRepo
	Checker *checker.ServerChecker
}

// NewHandlers creates the handlers
func NewHandlers(repo *DBRepo, a *config.AppConfig, c *checker.ServerChecker) {
	Repo = repo
	app = a
	serverChecker = c
}

// NewPostgresqlHandlers creates db repo for postgres
func NewPostgresqlHandlers(db *driver.DB, a *config.AppConfig, c *checker.ServerChecker) *DBRepo {
	return &DBRepo{
		App:     a,
		DB:      dbrepo.NewPostgresRepo(db.SQL, a),
		Checker: c,
	}
}

// AdminDashboard displays the dashboard
func (repo *DBRepo) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	vars := make(jet.VarMap)

	serviceStateCount, err := repo.DB.GetAllServiceStatusCounts()
	if err != nil {
		log.Println(err)
		return
	}
	vars.Set("no_healthy", serviceStateCount[models.ServiceStatusHealthy])
	vars.Set("no_problem", serviceStateCount[models.ServiceStatusProblem])
	vars.Set("no_pending", serviceStateCount[models.ServiceStatusPending])
	vars.Set("no_warning", serviceStateCount[models.ServiceStatusWarning])

	hosts, err := repo.DB.GetAllHost()
	if err != nil {
		log.Println(err)
		return
	}
	vars.Set("hosts", hosts)

	err = helpers.RenderPage(w, r, "dashboard", vars, nil)
	if err != nil {
		printTemplateError(w, err)
	}
}

// Events displays the events page
func (repo *DBRepo) Events(w http.ResponseWriter, r *http.Request) {
	events, err := repo.DB.GetAllEvent()
	if err != nil {
		log.Println(err)
		return
	}

	data := make(jet.VarMap)
	data.Set("events", events)

	err = helpers.RenderPage(w, r, "events", data, nil)
	if err != nil {
		printTemplateError(w, err)
	}
}

// Settings displays the settings page
func (repo *DBRepo) Settings(w http.ResponseWriter, r *http.Request) {
	err := helpers.RenderPage(w, r, "settings", nil, nil)
	if err != nil {
		printTemplateError(w, err)
	}
}

// PostSettings saves site settings
func (repo *DBRepo) PostSettings(w http.ResponseWriter, r *http.Request) {
	prefMap := make(map[string]string)

	prefMap["site_url"] = r.Form.Get("site_url")
	prefMap["notify_name"] = r.Form.Get("notify_name")
	prefMap["notify_email"] = r.Form.Get("notify_email")
	prefMap["smtp_server"] = r.Form.Get("smtp_server")
	prefMap["smtp_port"] = r.Form.Get("smtp_port")
	prefMap["smtp_user"] = r.Form.Get("smtp_user")
	prefMap["smtp_password"] = r.Form.Get("smtp_password")
	prefMap["sms_enabled"] = r.Form.Get("sms_enabled")
	prefMap["sms_provider"] = r.Form.Get("sms_provider")
	prefMap["twilio_phone_number"] = r.Form.Get("twilio_phone_number")
	prefMap["twilio_sid"] = r.Form.Get("twilio_sid")
	prefMap["twilio_auth_token"] = r.Form.Get("twilio_auth_token")
	prefMap["smtp_from_email"] = r.Form.Get("smtp_from_email")
	prefMap["smtp_from_name"] = r.Form.Get("smtp_from_name")
	prefMap["notify_via_sms"] = r.Form.Get("notify_via_sms")
	prefMap["notify_via_email"] = r.Form.Get("notify_via_email")
	prefMap["sms_notify_number"] = r.Form.Get("sms_notify_number")

	if r.Form.Get("sms_enabled") == "0" {
		prefMap["notify_via_sms"] = "0"
	}

	err := repo.DB.InsertOrUpdateSitePreferences(prefMap)
	if err != nil {
		log.Println(err)
		ClientError(w, r, http.StatusBadRequest)
		return
	}

	// update app config
	for k, v := range prefMap {
		app.PreferenceMap[k] = v
	}

	app.Session.Put(r.Context(), "flash", "Changes saved")

	if r.Form.Get("action") == "1" {
		http.Redirect(w, r, "/admin/overview", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/admin/settings", http.StatusSeeOther)
	}
}

// AllHosts displays list of all hosts
func (repo *DBRepo) AllHosts(w http.ResponseWriter, r *http.Request) {
	hosts, err := repo.DB.GetAllHost()
	if err != nil {
		log.Println(err)
		return
	}

	vars := make(jet.VarMap)
	vars.Set("hosts", hosts)

	err = helpers.RenderPage(w, r, "hosts", vars, nil)
	if err != nil {
		printTemplateError(w, err)
	}
}

// Host shows the host add/edit form
func (repo *DBRepo) Host(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	var h models.Host
	if id > 0 {
		// get the host from the database
		host, err := repo.DB.GetHostByID(id)
		if err != nil {
			log.Println(err)
			return
		}
		h = host
	}

	// h.HostName = "Localhost"
	// h.CanonicalName = "localhost"
	// h.URL = "http://localhost:8080"
	vars := make(jet.VarMap)
	vars.Set("host", h)
	vars.Set("form_id", "host-form")

	err := helpers.RenderPage(w, r, "host", vars, nil)
	if err != nil {
		printTemplateError(w, err)
	}
}

// PostHost handles posting of host form
func (repo *DBRepo) PostHost(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintln(w, "Post Success")
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	var h models.Host

	if id > 0 {
		// get the host from the database
		host, err := repo.DB.GetHostByID(id)
		if err != nil {
			log.Println(err)
			return
		}
		h = host
	}
	h.HostName = r.FormValue("host_name")
	h.CanonicalName = r.FormValue("canonical_name")
	h.URL = r.FormValue("url")
	h.IP = r.FormValue("ip")
	h.IPv6 = r.FormValue("ipv6")
	h.Location = r.FormValue("location")
	h.OS = r.FormValue("os")
	h.Active, _ = strconv.Atoi(r.FormValue("active"))

	if id > 0 {
		err := repo.DB.UpdateHost(h)
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		newID, err := repo.DB.InsertHost(h)
		if err != nil {
			log.Println("error when inserting host")
			log.Println(err)
			helpers.ServerError(w, r, err)
			return
		}
		h.ID = newID
	}

	repo.App.Session.Put(r.Context(), "flash", "Change saved")
	http.Redirect(w, r, fmt.Sprintf("/admin/host/%d", h.ID), http.StatusSeeOther)
}

// AllUsers lists all admin users
func (repo *DBRepo) AllUsers(w http.ResponseWriter, r *http.Request) {
	vars := make(jet.VarMap)

	u, err := repo.DB.AllUsers()
	if err != nil {
		ClientError(w, r, http.StatusBadRequest)
		return
	}

	vars.Set("users", u)

	err = helpers.RenderPage(w, r, "users", vars, nil)
	if err != nil {
		printTemplateError(w, err)
	}
}

// OneUser displays the add/edit user page
func (repo *DBRepo) OneUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Println(err)
	}

	vars := make(jet.VarMap)

	if id > 0 {

		u, err := repo.DB.GetUserById(id)
		if err != nil {
			ClientError(w, r, http.StatusBadRequest)
			return
		}

		vars.Set("user", u)
	} else {
		var u models.User
		vars.Set("user", u)
	}

	err = helpers.RenderPage(w, r, "user", vars, nil)
	if err != nil {
		printTemplateError(w, err)
	}
}

// PostOneUser adds/edits a user
func (repo *DBRepo) PostOneUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Println(err)
	}

	var u models.User

	if id > 0 {
		u, _ = repo.DB.GetUserById(id)
		u.FirstName = r.Form.Get("first_name")
		u.LastName = r.Form.Get("last_name")
		u.Email = r.Form.Get("email")
		u.UserActive, _ = strconv.Atoi(r.Form.Get("user_active"))
		err := repo.DB.UpdateUser(u)
		if err != nil {
			log.Println(err)
			ClientError(w, r, http.StatusBadRequest)
			return
		}

		if len(r.Form.Get("password")) > 0 {
			// changing password
			err := repo.DB.UpdatePassword(id, r.Form.Get("password"))
			if err != nil {
				log.Println(err)
				ClientError(w, r, http.StatusBadRequest)
				return
			}
		}
	} else {
		u.FirstName = r.Form.Get("first_name")
		u.LastName = r.Form.Get("last_name")
		u.Email = r.Form.Get("email")
		u.UserActive, _ = strconv.Atoi(r.Form.Get("user_active"))
		u.Password = []byte(r.Form.Get("password"))
		u.AccessLevel = 3

		_, err := repo.DB.InsertUser(u)
		if err != nil {
			log.Println(err)
			ClientError(w, r, http.StatusBadRequest)
			return
		}
	}

	repo.App.Session.Put(r.Context(), "flash", "Changes saved")
	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}

// DeleteUser soft deletes a user
func (repo *DBRepo) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	_ = repo.DB.DeleteUser(id)
	repo.App.Session.Put(r.Context(), "flash", "User deleted")
	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}

type serviceJSON struct {
	OK bool `json:"ok"`
}

//
func (repo *DBRepo) ToggleServiceForHost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
		return
	}

	hostID, _ := strconv.Atoi(r.FormValue("host_id"))
	serviceID, _ := strconv.Atoi(r.FormValue("service_id"))
	active, _ := strconv.Atoi(r.FormValue("active"))

	log.Printf("Host ID: %d Service ID: %d Active: %d", hostID, serviceID, active)
	err = repo.DB.UpdateHostServiceStatusByID(hostID, serviceID, active)

	var out []byte
	if err != nil {
		log.Println(err)
		out, _ = json.MarshalIndent(serviceJSON{OK: false}, "", "    ")
	} else {
		out, _ = json.MarshalIndent(serviceJSON{OK: true}, "", "    ")
	}

	// broadcast
	hs, err := repo.DB.GetHostByHostIDServiceID(hostID, serviceID)
	if err != nil {
		log.Println(err)
	}

	h, err := repo.DB.GetHostByID(hostID)
	if err != nil {
		log.Println(err)
	}

	// add or remove host service from schedule
	if active == 1 {
		// add to schedule
		log.Printf("add to schedule")
		data, err := repo.AddtoMonitorMap(hs, h.HostName)
		if err != nil {
			log.Println(err)
		} else {
			log.Println("broadcast message")
			repo.broadcastMessage("public-channel", "schedule-changed-event", data)
			repo.pushScheduleChangeEvent(hs.ID, hs.HostID, h.HostName, hs.ServiceID, hs.Service.ServiceName,
				hs.Service.Icon, hs.ScheduleUnit, hs.ScheduleNumber, time.Now(), models.ServiceStatusPending)
			repo.pushServerStatusChangeEvent(hs.ID, hs.HostID, h.HostName, hs.ServiceID, hs.Service.ServiceName,
				hs.Service.Icon, hs.ScheduleUnit, hs.ScheduleNumber, time.Now(), hs.Status, models.ServiceStatusPending)
		}
	} else {
		// remove from schedule
		log.Println("remove from schedule")
		data, err := repo.RemoveFromMonitorMap(hs.ID)
		if err != nil {
			log.Println(err)
		} else {
			log.Println("broadcast message to schedule-item-removed-event")
			repo.broadcastMessage("public-channel", "schedule-item-removed-event", data)
		}
	}
	repo.updateHostServiceCount("schedule has changed")

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (repo *DBRepo) SetSystemPref(w http.ResponseWriter, r *http.Request) {
	prefName := models.ParseSystemPreference(r.PostForm.Get("pref_name"))
	prefValue := r.PostForm.Get("pref_value")

	resp := jsonResp{OK: true, Message: ""}

	err := repo.DB.SetSystemPref(prefName, prefValue)
	if err != nil {
		resp.OK = false
		resp.Message = err.Error()
	}

	// log.Printf("Change %s from %s to %s\n", prefName.String(), repo.App.PreferenceMap[prefName.String()], prefValue)
	repo.App.PreferenceMap[prefName.String()] = prefValue

	out, _ := json.MarshalIndent(resp, "", "    ")

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// turns monitoring on and off
func (repo *DBRepo) ToggleMonitoring(w http.ResponseWriter, r *http.Request) {
	enabled := r.PostForm.Get("enabled")
	// log.Println("enabled:", enabled)

	resp := jsonResp{OK: true, Message: ""}

	switch enabled {
	case "1":
		log.Println("Turning monitoring on")
		repo.App.PreferenceMap[models.MonitoringLive.String()] = enabled
		repo.StartMonitoring()
		repo.App.Scheduler.Start()
	case "0":
		log.Println("Turning monitoring off")
		repo.App.PreferenceMap[models.MonitoringLive.String()] = enabled
		repo.StopMonitoring()
		repo.App.Scheduler.Stop()
	default:
		log.Println("unknown enable value", enabled)
		resp.OK = false
		resp.Message = fmt.Sprintf("unknown enable value: %q", enabled)
	}
	// if err != nil {
	// 	resp.OK = false
	// 	resp.Message = err.Error()
	// }

	out, _ := json.MarshalIndent(resp, "", "    ")

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// ClientError will display error page for client error i.e. bad request
func ClientError(w http.ResponseWriter, r *http.Request, status int) {
	switch status {
	case http.StatusNotFound:
		show404(w, r)
	case http.StatusInternalServerError:
		show500(w, r)
	default:
		http.Error(w, http.StatusText(status), status)
	}
}

// ServerError will display error page for internal server error
func ServerError(w http.ResponseWriter, r *http.Request, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	_ = log.Output(2, trace)
	show500(w, r)
}

func show404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
	http.ServeFile(w, r, "./ui/static/404.html")
}

func show500(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
	http.ServeFile(w, r, "./ui/static/500.html")
}

func printTemplateError(w http.ResponseWriter, err error) {
	_, _ = fmt.Fprintf(w, "<small><span class='text-danger'>Error executing template: %s</span></small>", err)
}
