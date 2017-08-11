package handlers

import (
	"net/http"

	log "github.com/Sirupsen/logrus"

	"github.com/aarjan/blog/app/api"
	"github.com/aarjan/blog/app/shared"
)

type AppHandler struct {
	*api.AppService
	HanlderFunc func(w http.ResponseWriter, r *http.Request, app *api.AppService) error
}

func (app AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err := app.HanlderFunc(w, r, app.AppService)

	// Todo	:	Proper Error Handling
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err.Error(),
		}).Warn()
		return
	}

}

func AuthMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		log.WithFields(log.Fields{
			"Referrer":   r.Referer(),
			"CurrentURL": r.URL.Path,
			"Method":     r.Method,
		}).Info("Auth Handler Logging!!")

		// Redirect unauthorized personnel
		if shared.GetCookie(r).AccessToken == nil {
			log.Warn("Auth Handler Redirecting!!")
			data := shared.Meta{false, 401, "Unathorized Personnel", nil}
			shared.SetCookie(w, data, 1000)
			http.Redirect(w, r, "/api/v1", http.StatusFound)
			return
		}

		h.ServeHTTP(w, r)
	})
}
