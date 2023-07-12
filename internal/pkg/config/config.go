package config

import (
	"booking/internal/pkg/models"
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
)

type AppConfig struct {
	UseCashe      bool
	TemplateCashe map[string]*template.Template
	Session       *scs.SessionManager
	IsProduction  bool
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	EmailChan     chan models.EmailData
}
