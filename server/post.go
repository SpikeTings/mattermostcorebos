package server

import (
	"encoding/json"
	"io"
	"mattermostcorebos/entity"
	"mattermostcorebos/helpers"
	"net/http"

	"github.com/mattermost/mattermost/server/public/model"
)

func (p *Plugin) postMessage(w http.ResponseWriter, r *http.Request) {
	rawBody, err := io.ReadAll(r.Body)
	if err != nil {
		helpers.DisplayAppErrorResponse(w, "There were a problem parsing request", http.StatusBadRequest)
		return
	}

	incomingWebhookRequest := model.IncomingWebhookRequest{}
	incomingWebhook := model.IncomingWebhook{}
	post := model.Post{}
	postHelper := entity.PostHelper{}
	err = json.Unmarshal(rawBody, &incomingWebhookRequest)
	if err != nil {
		helpers.DisplayAppErrorResponse(w, "There was an error decoding json user", http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(rawBody, &incomingWebhook)
	if err != nil {
		helpers.DisplayAppErrorResponse(w, "There was an error decoding json user", http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(rawBody, &post)
	if err != nil {
		helpers.DisplayAppErrorResponse(w, "There was an error decoding json user", http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(rawBody, &postHelper)
	if err != nil {
		helpers.DisplayAppErrorResponse(w, "There was an error decoding json user", http.StatusBadRequest)
		return
	}
	post.Message = incomingWebhookRequest.Text
	if incomingWebhookRequest.Props != nil {
		post.Props = incomingWebhookRequest.Props
	}
	post.AddProp("attachments", incomingWebhookRequest.Attachments)
	if post.Message == "" && postHelper.EphemeralText != "" {
		post.Message = postHelper.EphemeralText
		p.API.SendEphemeralPost(post.UserId, &post)
		return
	}
	p.API.CreatePost(&post)
}
