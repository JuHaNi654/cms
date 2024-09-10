package routes

import (
	"html/template"
	"log"
	"net/http"

	"github.com/JuHaNi654/cms/internal/models"
)

type Handler func(w http.ResponseWriter, r *http.Request) error

// TODO render error html templates/pages instead of the text
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h(w, r); err != nil {
		log.Printf("Error: %s\n", err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		path := models.Environment.WithRoot("/ui/login.tmpl")
		tmpl, err := template.ParseFiles(path)
		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}

		if err := tmpl.Execute(w, tmpl); err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
	}
}

func notFound(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "page is not found", http.StatusNotFound)
	return
}
