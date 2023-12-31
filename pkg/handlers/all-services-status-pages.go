package handlers

import (
	"log"
	"net/http"

	"github.com/CloudyKit/jet/v6"
	"gitlab.com/gjerry134679/vigilate/pkg/helpers"
	"gitlab.com/gjerry134679/vigilate/pkg/models"
)

// AllHealthyServices lists all healthy services
func (repo *DBRepo) AllHealthyServices(w http.ResponseWriter, r *http.Request) {
	hsNamePair, err := repo.DB.GetServiceByStatus(models.ServiceStatusHealthy)
	if err != nil {
		log.Println(err)
	}

	vars := make(jet.VarMap)
	vars.Set("services", hsNamePair)
	err = helpers.RenderPage(w, r, "healthy", vars, nil)
	if err != nil {
		printTemplateError(w, err)
	}
}

// AllWarningServices lists all warning services
func (repo *DBRepo) AllWarningServices(w http.ResponseWriter, r *http.Request) {
	hsNamePair, err := repo.DB.GetServiceByStatus(models.ServiceStatusWarning)
	if err != nil {
		log.Println(err)
	}

	vars := make(jet.VarMap)
	vars.Set("services", hsNamePair)
	err = helpers.RenderPage(w, r, "warning", vars, nil)
	if err != nil {
		printTemplateError(w, err)
	}
}

// AllProblemServices lists all problem services
func (repo *DBRepo) AllProblemServices(w http.ResponseWriter, r *http.Request) {
	hsNamePair, err := repo.DB.GetServiceByStatus(models.ServiceStatusProblem)
	if err != nil {
		log.Println(err)
	}

	vars := make(jet.VarMap)
	vars.Set("services", hsNamePair)
	err = helpers.RenderPage(w, r, "problems", vars, nil)
	if err != nil {
		printTemplateError(w, err)
	}
}

// AllPendingServices lists all pending services
func (repo *DBRepo) AllPendingServices(w http.ResponseWriter, r *http.Request) {
	hsNamePair, err := repo.DB.GetServiceByStatus(models.ServiceStatusPending)
	if err != nil {
		log.Println(err)
	}

	vars := make(jet.VarMap)
	vars.Set("services", hsNamePair)
	err = helpers.RenderPage(w, r, "pending", vars, nil)
	if err != nil {
		printTemplateError(w, err)
	}
}
