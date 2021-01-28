package corebos

import (
	"encoding/json"
	"github.com/tsolucio/corebosgowslib"
	"mattermost-server-plugin/configuration"
)

func DoInvoke(method, module, typeofCall string, element map[string]interface{}) (map[string]interface{}, error) {

	wsContext := corebosgowslib.GetCbContext()
	_, err := wsContext.DoLogin(configuration.CorebosUserName, configuration.CorebosUserPassword, true)
	defer wsContext.DoLogout()
	if err != nil {
		return nil, err
	}
	element["assigned_user_id"] = wsContext.GetUserId()
	byte, err := json.Marshal(element)
	if err != nil {
		return nil, err
	}

	data := map[string]interface{}{
		"element":     string(byte),
		"elementType": module,
	}

	response, err := wsContext.DoInvoke(method, data, typeofCall)
	if err != nil {
		return nil, err
	}
	return response, nil

}
