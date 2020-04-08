package server_plugin

import (
	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost-server/plugin"
	"mattermost-server-plugin/middleware"
	"net/http"
)

func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	router := mux.NewRouter()

	session, _ := p.API.GetSession(c.SessionId)
	auth := middleware.AuthenticationMiddleware{Session: session}

	// Protected router section
	protected := router.PathPrefix("/").Subrouter()
	protected.HandleFunc("/team/{team-name}/project/{name}/documents", p.GetDocumentsForProject).Methods(http.MethodGet)

	protected.HandleFunc("/team/{team-name}/project/{name}/wiki", p.CreateWiki).Methods(http.MethodPost)
	protected.HandleFunc("/team/{team-name}/project/{name}/wiki", p.UpdateWiki).Methods(http.MethodPut)
	protected.HandleFunc("/team/{team-name}/project/{name}/wiki", p.GetWikies).Methods(http.MethodGet)

	// Public router section
	router.Path("/postmessage").HandlerFunc(p.postMessage).Methods(http.MethodPost)
	router.Path("/syncuser").HandlerFunc(p.syncUserWithCoreBOS).Methods(http.MethodPost)

	protected.Use(auth.CheckAuthentication)
	router.ServeHTTP(w, r)
}
