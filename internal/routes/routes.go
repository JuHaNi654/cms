package routes

import (
	"log"
	"net/http"
	"time"

	"github.com/JuHaNi654/cms/internal/models"
	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/form"
	"github.com/go-playground/validator/v10"
)

var (
	decoder        *form.Decoder
	validate       *validator.Validate
	sessionManager *scs.SessionManager
)

func Routes(services *models.Services) http.Handler {
	decoder = form.NewDecoder()
	validate = validator.New()
	validate.RegisterValidation("match-passwords", ValidateMatchPasswords)

	sessionManager = scs.New()
	sessionManager.Lifetime = 6 * time.Hour
	sessionManager.IdleTimeout = 20 * time.Minute
	sessionManager.Cookie.Name = "_sid"
	sessionManager.Cookie.SameSite = http.SameSiteStrictMode
	sessionManager.Cookie.Secure = models.Environment.IsProduction()

	sessionManager.Store = sqlite3store.New(services.DB.Db)

	r := chi.NewRouter()

	// TODO: Register assets folder when production mode
	// TODO: Register file system maybe

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(isInstalled(services))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("healthy"))
	})

	r.Method("GET", "/login", Handler(login(services)))
	r.Method("POST", "/login", Handler(authentication(services)))
	r.Method("POST", "/logout", Handler(logout))
	r.Method("POST", "/install", Handler(install(services)))

	// Private routes
	r.Group(func(r chi.Router) {
		r.Use(isAuthenticated)
		r.Method("GET", "/dashboard", Handler(dashboard(services)))
    r.Method("GET", "/editor", Handler(editor(services)))	
  })

	return sessionManager.LoadAndSave(r)
}

// Check matching passwords
func ValidateMatchPasswords(fl validator.FieldLevel) bool {
	s, ok := fl.Parent().Interface().(models.Password)
	if !ok {
		log.Println("error: could not assert type in ValidateMatchPasswords function")
		return false
	}

	return s.GetPassword() == s.GetMatchingPassword()
}
