package server_plugin

import (
	"github.com/gorilla/mux"
	"mattermost-server-plugin/corebos"
	"mattermost-server-plugin/helpers"
	"net/http"
)

func (p *Plugin) DoInvoke(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	channelName := vars["name"]
	teamName := vars["team-name"]
	method := vars["method"]
	module := vars["module"]

	team, teamErr := p.API.GetTeamByName(teamName)
	if teamErr != nil {
		helpers.DisplayAppErrorResponse(w, "The team "+teamName+" does not exist!", http.StatusNotFound)
		return
	}

	_, channelErr := p.API.GetChannelByName(team.Id, channelName, false)
	if channelErr != nil {
		helpers.DisplayAppErrorResponse(w, "The channel "+channelName+" does not exist!", http.StatusNotFound)
		return
	}

	input, err := helpers.ReadRequestBody(r)
	if err != nil {
		helpers.DisplayAppErrorResponse(w, "There were a problem parsing request", http.StatusBadRequest)
		return
	}

	response, err := corebos.DoInvoke(method, module, r.Method, input)
	if err != nil {
		helpers.DisplayAppErrorResponse(w, "There was a problem getting project!", http.StatusInternalServerError)
		return
	}

	helpers.DisplayAppSuccessResponse(w, response, "The Task was created successfully!")
}
