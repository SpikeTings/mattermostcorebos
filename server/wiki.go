package server

import (
	"mattermostcorebos/configuration"
	"mattermostcorebos/helpers"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/thedevsaddam/govalidator"
	"github.com/tsolucio/corebosgowslib"
)

func (p *Plugin) CreateWiki(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	channelName := vars["name"]
	teamName := vars["team-name"]

	team, teamErr := p.API.GetTeamByName(teamName)
	if teamErr != nil {
		helpers.DisplayAppErrorResponse(w, "The team does not exist!", http.StatusNotFound)
		return
	}

	channel, channelErr := p.API.GetChannelByName(team.Id, channelName, false)
	if channelErr != nil {
		helpers.DisplayAppErrorResponse(w, "The channel does not exist!", http.StatusNotFound)
		return
	}

	input, err := helpers.ReadRequestBody(r)
	if err != nil {
		helpers.DisplayAppErrorResponse(w, "There were a problem parsing request", http.StatusBadRequest)
		return
	}

	rules := govalidator.MapData{
		"content": []string{"required"},
		"title":   []string{"required"},
	}

	data := make(map[string]interface{}, 0)
	opts := govalidator.Options{
		Request: r,
		Rules:   rules,
		Data:    &data,
	}
	v := govalidator.New(opts)
	e := v.ValidateJSON()
	if len(e) > 0 {
		validationErrors := map[string]interface{}{"validation_error": e}
		helpers.DisplayAppErrorResponse(w, validationErrors, http.StatusBadRequest)
		return
	}

	wsContext := corebosgowslib.GetCbContext()
	_, err = wsContext.DoLogin(configuration.CorebosUserName, configuration.CorebosUserPassword, true)
	defer wsContext.DoLogout()
	if err != nil {
		helpers.DisplayAppErrorResponse(w, "There were a problem accessing coreBOS!", http.StatusInternalServerError)
		return
	}

	dq, err := wsContext.DoQuery("select id from Project where projectname = '" + channel.DisplayName + "' and teamname = '" + team.DisplayName + "';")
	if err != nil {
		helpers.DisplayAppErrorResponse(w, "There was a problem getting project!", http.StatusInternalServerError)
		return
	}

	projects := dq.([]interface{})
	if len(projects) == 0 {
		helpers.DisplayAppErrorResponse(w, "The project does not exist!", http.StatusNotFound)
		return
	}
	project := projects[0].(map[string]interface{})
	projectWsId := project["id"].(string)

	wikiData := map[string]interface{}{
		"conversationtitle":   input["title"],
		"dopublic":            "0",
		"conversationcontent": input["content"],
		"relations":           []string{projectWsId},
	}

	wiki, err := wsContext.DoCreate("Conversation", wikiData)
	if err != nil {
		helpers.DisplayAppErrorResponse(w, "There was a problem creating wiki!", http.StatusInternalServerError)
	}

	helpers.DisplayAppSuccessResponse(w, wiki, "The wiki was created successfully!")
}

func (p *Plugin) UpdateWiki(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	channelName := vars["name"]
	teamName := vars["team-name"]

	team, teamErr := p.API.GetTeamByName(teamName)
	if teamErr != nil {
		helpers.DisplayAppErrorResponse(w, "The team does not exist!", http.StatusNotFound)
		return
	}

	_, channelErr := p.API.GetChannelByName(team.Id, channelName, false)
	if channelErr != nil {
		helpers.DisplayAppErrorResponse(w, "The channel does not exist!", http.StatusNotFound)
		return
	}

	input, err := helpers.ReadRequestBody(r)
	if err != nil {
		helpers.DisplayAppErrorResponse(w, "There were a problem parsing request", http.StatusBadRequest)
		return
	}

	rules := govalidator.MapData{
		"content": []string{"required"},
		"title":   []string{"required"},
		"id":      []string{"required"},
	}

	data := make(map[string]interface{}, 0)
	opts := govalidator.Options{
		Request: r,
		Rules:   rules,
		Data:    &data,
	}
	v := govalidator.New(opts)
	e := v.ValidateJSON()
	if len(e) > 0 {
		validationErrors := map[string]interface{}{"validation_error": e}
		helpers.DisplayAppErrorResponse(w, validationErrors, http.StatusBadRequest)
		return
	}

	wikiData := map[string]interface{}{
		"conversationtitle":   input["title"],
		"conversationcontent": input["content"],
		"id":                  input["id"],
	}

	wsContext := corebosgowslib.GetCbContext()
	_, err = wsContext.DoLogin(configuration.CorebosUserName, configuration.CorebosUserPassword, true)
	defer wsContext.DoLogout()
	if err != nil {
		helpers.DisplayAppErrorResponse(w, "There were a problem accessing coreBOS!", http.StatusInternalServerError)
		return
	}

	wiki, err := wsContext.DoRevise("Conversation", wikiData)
	if err != nil {
		helpers.DisplayAppErrorResponse(w, "There was a problem updating wiki!", http.StatusInternalServerError)
	}

	helpers.DisplayAppSuccessResponse(w, wiki, "The wiki was updated successfully!")

}

func (p *Plugin) GetWikies(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	channelName := vars["name"]
	teamName := vars["team-name"]

	team, teamErr := p.API.GetTeamByName(teamName)
	if teamErr != nil {
		helpers.DisplayAppErrorResponse(w, "The team does not exist!", http.StatusNotFound)
		return
	}

	channel, channelErr := p.API.GetChannelByName(team.Id, channelName, false)
	if channelErr != nil {
		helpers.DisplayAppErrorResponse(w, "The channel does not exist!", http.StatusNotFound)
		return
	}

	wsContext := corebosgowslib.GetCbContext()
	_, err := wsContext.DoLogin(configuration.CorebosUserName, configuration.CorebosUserPassword, true)
	defer wsContext.DoLogout()
	if err != nil {
		helpers.DisplayAppErrorResponse(w, "There were a problem accessing coreBOS!", http.StatusInternalServerError)
		return
	}

	dq, err := wsContext.DoQuery("select id from Project where projectname = '" + channel.DisplayName + "' and teamname = '" + team.DisplayName + "';")
	if err != nil {
		helpers.DisplayAppErrorResponse(w, "There was a problem getting project!", http.StatusInternalServerError)
		return
	}

	projects := dq.([]interface{})
	if len(projects) == 0 {
		helpers.DisplayAppErrorResponse(w, "The project does not exist!", http.StatusNotFound)
		return
	}
	project := projects[0].(map[string]interface{})
	projectWsId := project["id"].(string)

	wikies, err := wsContext.DoQuery("select * from Conversation where related.Project='" + projectWsId + "';")
	if err != nil {
		helpers.DisplayAppErrorResponse(w, "There was a problem getting wikies!", http.StatusInternalServerError)
		return
	}

	helpers.DisplayAppSuccessResponse(w, wikies, "The wikies were returned successfully!")
}
