package models

import (
	"github.com/JuHaNi654/cms/internal/database"
	"github.com/JuHaNi654/cms/internal/vite"
)

type Services struct {
	Vite *vite.Handler
	DB   *database.SqlClient
}
