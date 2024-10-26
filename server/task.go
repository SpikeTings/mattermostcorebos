package server_plugin

import (
	"github.com/gorilla/mux"
	"github.com/thedevsaddam/govalidator"
	"github.com/tsolucio/corebosgowslib"
	"mattermost-server-plugin/configuration"
	"mattermost-server-plugin/helpers"
	"net/http"
	"time"
)

func (p *Plugin) CreateTaskForProject(w http.ResponseWriter, r *http.Request) {

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
		"projecttaskname":  []string{"required"},
		"protask_category": []string{"required"},
		"accrelated":       []string{"required"},
		"caller":           []string{"required"},
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
		helpers.DisplayAppErrorResponse(w, "There was a problem getting project! Where projectname = '"+channel.DisplayName+"' and teamname = '"+team.DisplayName+"';", http.StatusInternalServerError)
		return
	}

	projects := dq.([]interface{})
	if len(projects) == 0 {
		helpers.DisplayAppErrorResponse(w, "The project does not exist!", http.StatusNotFound)
		return
	}
	project := projects[0].(map[string]interface{})
	projectWsId := project["id"].(string)
	timeNow := time.Now()
	dateStamp := timeNow.Format("2006-01-02")
	dateTimeStamp := timeNow.Format("2006-01-02 15:04:05")
	taskData := map[string]interface{}{
		"protask_importdate": dateStamp,     //mandatory
		"protask_lastdate":   dateTimeStamp, //mandatory
		"protask_initstart":  dateStamp,     //mandatory
		"protask_initend":    dateStamp,     //mandatory
		"protask_wf23start":  dateStamp,     //mandatory
		"protask_wf23end":    dateStamp,     //mandatory
		"protask_wf22end":    dateStamp,     //mandatory
		"protask_pdsdate":    dateStamp,     //mandatory
	}
	taskData["projectid"] = projectWsId
	taskData["projecttaskname"] = input["projecttaskname"]
	taskData["protask_category"] = input["protask_category"]
	taskData["accrelated"] = input["accrelated"]
	taskData["caller"] = input["caller"]
	taskData["projecttaskpriority"] = input["projecttaskpriority"]
	taskData["projecttasktype"] = input["projecttasktype"]
	taskData["projecttaskstatus"] = input["projecttaskstatus"]
	taskData["startdate"] = input["startdate"]
	taskData["enddate"] = input["enddate"]
	taskData["description"] = input["description"]
	task, err := wsContext.DoCreate("ProjectTask", taskData)
	if err != nil {
		helpers.DisplayAppErrorResponse(w, "There was a problem Creating project!\n "+err.Error(), http.StatusInternalServerError)
		return
	}

	helpers.DisplayAppSuccessResponse(w, task, "The Task was created successfully!")
}
