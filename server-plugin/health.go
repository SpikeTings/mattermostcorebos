package server_plugin

import (
	"mattermost-server-plugin/helpers"
	"net/http"
)

func (p *Plugin) Health(w http.ResponseWriter, r *http.Request) {
	helpers.DisplayAppSuccessResponse(w, "Health check", "The plugin status is active")
}
