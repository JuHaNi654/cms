package routes

import (
	"context"
	"html/template"
	"net/http"

	"github.com/JuHaNi654/cms/internal/models"
	"github.com/JuHaNi654/cms/internal/vite"
	"github.com/JuHaNi654/cms/templates/partials"
	"github.com/a-h/templ"
)

func serveTemplate(
	w http.ResponseWriter,
	r *http.Request,
  user *models.User, 
  vite *vite.Handler,
	templates []string,
) error {	
  tmpl, err := template.ParseFiles(templates...)
	if err != nil {
		return err
	}

	data, err := vite.GetPageData(r)
	if err != nil {
		return err
	}

  if user != nil {
    header, _ := templ.ToGoHTML(context.Background(), partials.Header(user.Firstname))
    data.Header = header  
  }

	if err := tmpl.Execute(w, data); err != nil {
		return err
	}

	return nil
}
