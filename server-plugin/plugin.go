package server_plugin

import (
	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
	"github.com/tsolucio/corebosgowslib"
	"mattermost-server-plugin/configuration"
)

type Plugin struct {
	plugin.MattermostPlugin
}

func (p *Plugin) MessageHasBeenPosted(c *plugin.Context, post *model.Post) {
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
