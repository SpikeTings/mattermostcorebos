package server_plugin

import (
	"encoding/json"
	"io/ioutil"
	"mattermost-server-plugin/entity"
	"mattermost-server-plugin/helpers"
	"net/http"
)

func (p *Plugin) syncUserWithCoreBOS(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	rawBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		helpers.DisplayAppErrorResponse(w, "There were a problem parsing request", http.StatusBadRequest)
		return
	}

	userRequest := entity.User{}
	err = json.Unmarshal(rawBody, &userRequest)
	if err != nil {
		helpers.DisplayAppErrorResponse(w, "There were a problem decoding json user", http.StatusBadRequest)
		return
	}

	userCreate := userRequest.GetMMUser()
	userExist, appError := p.API.GetUserByUsername(userCreate.Username)
	if appError == nil {
		p.addTeam(*userExist, userRequest)
		userReturn := entity.User{}.GetUser(userExist)
		jsonValue, _ := json.Marshal(userReturn)
		w.Write(jsonValue)
		return
	}
	userExist, appError = p.API.GetUserByEmail(userCreate.Email)
	if appError == nil {
		p.addTeam(*userExist, userRequest)
		userReturn := entity.User{}.GetUser(userExist)
		jsonValue, _ := json.Marshal(userReturn)
		w.Write(jsonValue)
		return
	}

	userCreated, appError := p.API.CreateUser(&userCreate)
	if appError != nil && appError.StatusCode != http.StatusOK {
		helpers.DisplayAppErrorResponse(w, appError.Error(), http.StatusInternalServerError)
		return
	}

	p.addTeam(*userCreated, userRequest)
	userReturn := entity.User{}.GetUser(userCreated)
	jsonValue, _ := json.Marshal(userReturn)
	w.Write(jsonValue)
}
