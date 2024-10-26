package server_plugin

import (
	"errors"
	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
	"mattermost-server-plugin/configuration"
	"net/http"
	"net/url"
)

type Plugin struct {
	plugin.MattermostPlugin
}

func (p *Plugin) SendPostToChatWithMeExtension(post *model.Post, triggerWord string) error {

	cnl, _ := p.API.GetChannel(post.ChannelId)
	var tname = ""
	var tdname = ""
	var cdname = ""
	if cnl.Type == "D" {
		user, errr := p.API.GetUser(post.UserId)
		if errr != nil {
			return errr
		}
		cdname = user.FirstName + user.LastName
		tname = user.FirstName + "_" + user.LastName
		tdname = user.FirstName + "_" + user.LastName
	} else {
		team, _ := p.API.GetTeam(cnl.TeamId)
		tname = team.Name
		tdname = team.DisplayName
		cdname = cnl.DisplayName
	}
	formData := url.Values{
		"text":         {post.Message},
		"token":        {configuration.ChatWithMeToken},
		"trigger_word": {triggerWord},
		"user_id":      {post.UserId},
		"channel_id":   {post.ChannelId},
		"team_id":      {cnl.TeamId},
		"team_name":    {tname},
		"team_dname":   {tdname},
		"chnl_name":    {cnl.Name},
		"chnl_dname":   {cdname},
	}

	newPost := &model.Post{
		UserId:    post.UserId,
		ChannelId: post.ChannelId,
		Type:      model.POST_SLACK_ATTACHMENT,
	}
	resp, err := http.PostForm(configuration.ChatWithMeExtensionUrl, formData)
	defer resp.Body.Close()

	if err != nil {
		return err
	}

	incomingWebhookPayload, decodeError := model.IncomingWebhookRequestFromJson(resp.Body)
	if decodeError != nil {
		return decodeError
	}

	if len(incomingWebhookPayload.Text) == 0 && incomingWebhookPayload.Attachments == nil {
		return errors.New("Wrong response format")
	}

	if incomingWebhookPayload.Props != nil {
		newPost.Props = incomingWebhookPayload.Props
	}
	newPost.Message = incomingWebhookPayload.Text
	newPost.AddProp("attachments", incomingWebhookPayload.Attachments)

	p.API.SendEphemeralPost(newPost.UserId, newPost)
	return nil
}
