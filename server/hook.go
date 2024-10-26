package server

import (
	"errors"
	"mattermostcorebos/configuration"
	"mattermostcorebos/configuration/language"
	"mattermostcorebos/helpers"
	"regexp"
	"strings"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/tsolucio/corebosgowslib"
)

func (p *Plugin) OnConfigurationChange() error {
	var c configuration.MattermostConfig
	err := p.API.LoadPluginConfiguration(&c)
	if err != nil {
		return err
	}
	c.UpdateConfigurations()
	corebosgowslib.SetURL(configuration.CorebosUrl)
	return nil
}

func (p *Plugin) OnActivate() error {
	teams, err := p.API.GetTeams()
	if err != nil {
		return err
	}
	if len(teams) == 0 {
		return errors.New("there are no existing teams")
	}

	team := teams[0]
	channel, _ := p.API.GetChannelByNameForTeamName(team.Name, "chatwithme", false)
	if channel == nil {
		channel, err = p.API.CreateChannel(&model.Channel{
			TeamId:      team.Id,
			Type:        model.ChannelTypeOpen,
			DisplayName: "Chat With Me",
			Name:        "chatwithme",
			Header:      "The channel used by the mattermost-extend plugin.",
			Purpose:     "The channel was created by the mattermost-extend plugin to extend the server functionality.",
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Plugin) MessageWillBePosted(c *plugin.Context, post *model.Post) (*model.Post, string) {
	r, _ := regexp.Compile("^\\S+")
	triggerWord := r.FindString(post.Message)
	if helpers.Contains(configuration.ChatWithMeTriggerWordsEphemeral, triggerWord) {
		p.SendPostToChatWithMeExtension(post, triggerWord)
		p.API.SendEphemeralPost(post.UserId, post)
		post.Message = "Posted Ephemeral Trigger Word"
	}
	return post, ""
}

func (p *Plugin) MessageHasBeenPosted(c *plugin.Context, post *model.Post) {

	//Regular expression used for the replacement logic of incoming and outgoing webhooks
	r, _ := regexp.Compile("^\\S+")
	triggerWord := r.FindString(post.Message)

	if helpers.Contains(configuration.ChatWithMeTriggerWords, triggerWord) {
		p.SendPostToChatWithMeExtension(post, triggerWord)
	}

	//Regular expression user for special commands like: open, create, edit, list
	r, _ = regexp.Compile("^#(\\w+) (\\w+)(?: (\\d+))?$")
	matches := r.FindStringSubmatch(strings.TrimSpace(post.Message))

	if len(matches) > 0 {
		if action, ok := language.Command[matches[1]]; ok {
			module := matches[2]
			broadcast := &model.WebsocketBroadcast{UserId: post.UserId}
			payloadData := map[string]interface{}{
				"action": action,
				"module": module,
			}
			if matches[3] != "" {
				payloadData["id"] = matches[3]
			}
			p.API.PublishWebSocketEvent("corebos", payloadData, broadcast)
		}
	}

	// Transfer files in coreBOS
	if len(post.FileIds) > 0 {
		channel, err := p.API.GetChannel(post.ChannelId)
		if err != nil {
			return
		}
		team, err := p.API.GetTeam(channel.TeamId)
		if err != nil {
			return
		}
		p.UploadFilesToCoreBOS(post.FileIds, channel.DisplayName, team.DisplayName)

	}
}
