package middleware

import (
	"github.com/mattermost/mattermost-server/model"
	"mattermost-server-plugin/helpers"
	"net/http"
)

type AuthenticationMiddleware struct {
	Session *model.Session
}

func (a *AuthenticationMiddleware) CheckAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if a.Session == nil {
			helpers.DisplayAppErrorResponse(w, "Invalid session", http.StatusForbidden)
			return
		}
		if a.Session.UserId == "" {
			helpers.DisplayAppErrorResponse(w, "Empty user", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
