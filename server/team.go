package server

import (
	"mattermostcorebos/configuration"
	"mattermostcorebos/entity"
	"mattermostcorebos/helpers"
	"strings"

	"github.com/mattermost/mattermost/server/public/model"
)

func (p *Plugin) addTeam(user model.User, userHelper entity.User) {
	Client := model.NewAPIv4Client(configuration.MatterMostHost)
	Client.Login(configuration.MatterMostAdminUsername, configuration.MatterMostAdminPassword) //admin credencials
	teams, appError := p.API.GetTeams()
	if appError != nil {
		return
	}
	teamsRequest := strings.Split(userHelper.TeamNames, ",")
	for _, team := range teams {
		if !helpers.Contains(teamsRequest, team.DisplayName) {
			continue
		}
		Client.AddTeamMember(team.Id, user.Id)
	}
}
