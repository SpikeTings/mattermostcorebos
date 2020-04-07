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
	// Router section
	router.HandleFunc("/team/{team-name}/project/{name}/documents", p.GetDocumentsForProject).Methods(http.MethodGet)

	router.HandleFunc("/team/{team-name}/project/{name}/wiki", p.CreateWiki).Methods(http.MethodPost)
	router.HandleFunc("/team/{team-name}/project/{name}/wiki", p.UpdateWiki).Methods(http.MethodPut)
	router.HandleFunc("/team/{team-name}/project/{name}/wiki", p.GetWikies).Methods(http.MethodGet)

	router.Use(auth.CheckAuthentication)
	router.ServeHTTP(w, r)
}
