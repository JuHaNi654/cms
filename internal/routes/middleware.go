package routes

import (
	"log"
	"net/http"

	"github.com/JuHaNi654/cms/internal/models"
	"github.com/JuHaNi654/cms/internal/vite"
)

type Middleware func(next http.Handler) http.Handler

func isInstalled(services *models.Services) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/" && r.Method == "GET" {
				isInstalled, err := models.IsServceInitialized(services.DB)
				if err != nil {
					log.Println(err)
					http.Error(w, "Something went wrong", http.StatusInternalServerError)
					return
				}

				if !isInstalled {
					ctx := r.Context()
					ctx = vite.MetadataToContext(ctx, vite.Metadata{
						Title: "Install",
					})

					templates := []string{
						models.Environment.WithRoot("/templates/pages/base.tmpl"),
						models.Environment.WithRoot("/templates/pages/install.tmpl"),
					}
					err := serveTemplate(w, r.WithContext(ctx), nil, services.Vite, templates)
					if err != nil {
						log.Println(err)
						http.Error(w, "Something went wrong", http.StatusInternalServerError)
					}

					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

func isAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !sessionManager.Exists(r.Context(), "id") {
			http.Error(w, "Page not found", http.StatusNotFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}
