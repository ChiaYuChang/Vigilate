package dbrepo

import (
	"database/sql"

	"gitlab.com/gjerry134679/vigilate/pkg/config"
	"gitlab.com/gjerry134679/vigilate/pkg/repository"
)

var app *config.AppConfig

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

// NewPostgresRepo creates the repository
func NewPostgresRepo(Conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	app = a
	return &postgresDBRepo{
		App: a,
		DB:  Conn,
	}
}
