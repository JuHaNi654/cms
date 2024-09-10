package routes

import (
	"fmt"
	"log"
	"net/http"

	"github.com/JuHaNi654/cms/internal/models"
	"github.com/JuHaNi654/cms/templates/partials"
	"github.com/go-playground/validator/v10"
)

/*
* Method:       POST
* Route:        /login
* Access:       Public
* Description:  Authenticate user to the system
 */
func authentication(s *models.Services) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if err := r.ParseForm(); err != nil {
			return err
		}

		var login models.Login

		if err := decoder.Decode(&login, r.Form); err != nil {
			return fmt.Errorf("(/login) error while decoding form values: %v", err)
		}

		if err := validate.Struct(login); err != nil {
			errMsgs := []string{}
			for _, err := range err.(validator.ValidationErrors) {
				msg := models.GetErrorMessage(err.Field(), err.Tag(), err.Param())
				errMsgs = append(errMsgs, msg)
			}

			if len(errMsgs) > 0 {
				w.WriteHeader(http.StatusBadRequest)
				partials.FormErrors(errMsgs).Render(r.Context(), w)
				return nil
			}
		}

		match, id, err := login.Authenticate(s.DB)
		if err != nil {
			return err
		}

		if !match {
			msg := []string{"Invalid email or password"}
			w.WriteHeader(http.StatusUnauthorized)
			partials.FormErrors(msg).Render(r.Context(), w)
			return nil
		}

		sessionManager.Put(r.Context(), "id", id)

		log.Println("New session created", id)
		// On success redirect to the dashboard page
		w.Header().Add("Hx-Redirect", "/dashboard")
		http.Redirect(w, r, "/dashboard", http.StatusNoContent)

		return nil
	}
}

/*
* Method:       POST
* Route:        /logout
* Access:       Public
* Description:  -
 */
func logout(w http.ResponseWriter, r *http.Request) error {
	sessionManager.Destroy(r.Context())
	w.Header().Add("Hx-Redirect", "/login")
	http.Redirect(w, r, "/login", http.StatusNoContent)
	return nil
}

/*
* Method:       POST
* Route:        /install
* Access:       Public
* Description:  Accessible until first user is registered, then
*   this should return 404
 */
func install(s *models.Services) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		isInstalled, _ := models.IsServceInitialized(s.DB)
		if isInstalled {
			notFound(w, r)
			return nil
		}

		var install models.Install
		r.ParseForm()
		if err := decoder.Decode(&install, r.Form); err != nil {
			return fmt.Errorf("(/install) error while decoding form values: %v", err)
		}

		if err := validate.Struct(install); err != nil {
			errMsgs := []string{}
			for _, err := range err.(validator.ValidationErrors) {
				msg := models.GetErrorMessage(err.Field(), err.Tag(), err.Param())
				errMsgs = append(errMsgs, msg)
			}

			if len(errMsgs) > 0 {
				w.WriteHeader(http.StatusBadRequest)
				partials.FormErrors(errMsgs).Render(r.Context(), w)
				return nil
			}
		}

		if err := install.SaveUser(s.DB); err != nil {
			return fmt.Errorf("(/install) error while saving user to the database: %v", err)
		}

		w.Header().Add("Hx-Redirect", "/login")
		http.Redirect(w, r, "/login", http.StatusNoContent)
		return nil
	}
}
