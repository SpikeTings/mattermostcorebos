package configuration

import "mattermostcorebos/helpers"

var ChatWithMeToken string
var ChatWithMeExtensionUrl string
var ChatWithMeTriggerWords []string
var ChatWithMeTriggerWordsEphemeral []string

var MatterMostHost string
var MatterMostAdminUsername string
var MatterMostAdminPassword string

var CorebosUrl string
var CorebosUserName string
var CorebosUserPassword string

type MattermostConfig struct {
	ChatWithMeToken                 string
	ChatWithMeExtensionUrl          string
	ChatWithMeTriggerWords          string
	ChatWithMeTriggerWordsEphemeral string

	MatterMostHost          string
	MatterMostAdminUsername string
	MatterMostAdminPassword string

	CorebosUrl          string
	CorebosUserName     string
	CorebosUserPassword string
}

func (c *MattermostConfig) UpdateConfigurations() {
	ChatWithMeToken = c.ChatWithMeToken
	ChatWithMeExtensionUrl = c.ChatWithMeExtensionUrl
	ChatWithMeTriggerWords = helpers.ToArray(c.ChatWithMeTriggerWords, ",")
	ChatWithMeTriggerWordsEphemeral = helpers.ToArray(c.ChatWithMeTriggerWordsEphemeral, ",")

	MatterMostHost = helpers.RemoveIfISLast(c.MatterMostHost, "/")
	MatterMostAdminUsername = c.MatterMostAdminUsername
	MatterMostAdminPassword = c.MatterMostAdminPassword

	CorebosUrl = c.CorebosUrl
	CorebosUserName = c.CorebosUserName
	CorebosUserPassword = c.CorebosUserPassword
}
