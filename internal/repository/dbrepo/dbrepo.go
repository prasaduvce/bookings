package dbrepo

import (
	"database/sql"

	"github.com/prasaduvce/bookings/internal/config"
	"github.com/prasaduvce/bookings/internal/repository"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB *sql.DB
}

func NewPostGresRepo(conn *sql.DB, a *config.AppConfig) (repository.DataBaseRepo) {
	return &postgresDBRepo {
		App: a,
		DB: conn,
	}
}