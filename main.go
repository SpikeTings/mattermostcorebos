package main

import (
	"github.com/mattermost/mattermost-server/plugin"
	"mattermost-server-plugin/server-plugin"
)

func main() {
	plugin.ClientMain(&server_plugin.Plugin{})
}
