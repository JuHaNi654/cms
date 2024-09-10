package routes

import (
	"net/http"

	"github.com/JuHaNi654/cms/internal/models"
	"github.com/JuHaNi654/cms/internal/vite"
)

/*
* @Access public
* @Description Render login template
 */

func login(services *models.Services) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if sessionManager.Exists(r.Context(), "id") {
			http.Redirect(w, r, "/dashboard", http.StatusTemporaryRedirect)
		}

		ctx := r.Context()
		ctx = vite.MetadataToContext(ctx, vite.Metadata{
			Title: "Login",
		})

		templates := []string{
			models.Environment.WithRoot("/templates/pages/base.tmpl"),
			models.Environment.WithRoot("/templates/pages/login.tmpl"),
		}
		err := serveTemplate(w, r.WithContext(ctx), nil, services.Vite, templates)
		if err != nil {
			return err
		}

		return nil
	}
}

/*
* @Access private
* @Description Dashboard
 */
func dashboard(services *models.Services) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		userId := sessionManager.GetInt(r.Context(), "id")
		user, err := models.GetUser(services.DB, userId)
		if err != nil {
			return err
		}

		ctx := r.Context()
		ctx = vite.MetadataToContext(ctx, vite.Metadata{
			Title: "Dashboard",
		})

		// Content
		templates := []string{
			models.Environment.WithRoot("/templates/pages/base.tmpl"),
			models.Environment.WithRoot("/templates/pages/dashboard.tmpl"),
		}
		err = serveTemplate(w, r.WithContext(ctx), user, services.Vite, templates)
		if err != nil {
			return err
		}

		return nil
	}
}

/*
* @Access private
* @Description Content page editor
 */
func editor(services *models.Services) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {	
    userId := sessionManager.GetInt(r.Context(), "id")
		user, err := models.GetUser(services.DB, userId)
		if err != nil {
			return err
		}

    ctx := r.Context()
		ctx = vite.MetadataToContext(ctx, vite.Metadata{
			Title: "Edit page",
		})

		// Content
		templates := []string{
			models.Environment.WithRoot("/templates/pages/base.tmpl"),
			models.Environment.WithRoot("/templates/pages/editor.tmpl"),
		}
		err = serveTemplate(w, r.WithContext(ctx), user, services.Vite, templates)
		if err != nil {
			return err
		}


    return nil
  }
}

/*
* @Access private
* @Description User settings
 */
func userSettings(w http.ResponseWriter, r *http.Request) error {
	return nil
}
