package main

import (
	"mattermostcorebos/server"

	"github.com/mattermost/mattermost/server/public/plugin"
)

func main() {
	plugin.ClientMain(&server.Plugin{})
}
