package config

import (
	"html/template"

	"github.com/alexedwards/scs/v2"
)

type AppConfig struct {
	Templates map[string]*template.Template
	UseCache bool
	InProduction bool
	Session *scs.SessionManager
}