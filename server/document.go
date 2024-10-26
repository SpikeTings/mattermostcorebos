package server

import (
	b64 "encoding/base64"
	"mattermostcorebos/configuration"
	"mattermostcorebos/helpers"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/tsolucio/corebosgowslib"
)

func (p *Plugin) GetDocumentsForProject(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	channelName := vars["name"]
	teamName := vars["team-name"]

	team, teamErr := p.API.GetTeamByName(teamName)
	if teamErr != nil {
		helpers.DisplayAppErrorResponse(w, "The team does not exist!", http.StatusNotFound)
		return
	}

	channel, channelErr := p.API.GetChannelByName(team.Id, channelName, false)
	if channelErr != nil {
		helpers.DisplayAppErrorResponse(w, "The channel does not exist!", http.StatusNotFound)
		return
	}

	wsContext := corebosgowslib.GetCbContext()
	_, err := wsContext.DoLogin(configuration.CorebosUserName, configuration.CorebosUserPassword, true)
	defer wsContext.DoLogout()
	if err != nil {
		helpers.DisplayAppErrorResponse(w, "There were a problem accessing coreBOS!", http.StatusInternalServerError)
		return
	}

	dq, err := wsContext.DoQuery("select id from Project where projectname = '" + channel.DisplayName + "' and teamname = '" + team.DisplayName + "';")
	if err != nil {
		helpers.DisplayAppErrorResponse(w, "There was a problem getting project!", http.StatusInternalServerError)
		return
	}

	projects := dq.([]interface{})
	if len(projects) == 0 {
		helpers.DisplayAppErrorResponse(w, "The project does not exist!", http.StatusNotFound)
		return
	}
	project := projects[0].(map[string]interface{})
	projectWsId := project["id"].(string)

	documents, err := wsContext.DoQuery("select filelocationtype, filename, filetype, note_no, notecontent, notes_title, _downloadurl from Documents where Project.id='" + projectWsId + "';")
	if err != nil {
		helpers.DisplayAppErrorResponse(w, "There was a problem getting documents!", http.StatusInternalServerError)
		return
	}

	helpers.DisplayAppSuccessResponse(w, documents, "The documents were returned successfully!")
}

func (p *Plugin) UploadFilesToCoreBOS(fileIds []string, projectName string, teamName string) {
	cbContext := corebosgowslib.GetCbContext()
	_, err := cbContext.DoLogin(configuration.CorebosUserName, configuration.CorebosUserPassword, true)
	if err != nil {
		return
	}
	defer cbContext.DoLogout()

	projectsResult, err := cbContext.DoQuery("select id from Project where projectname = '" + projectName + "' and teamname = '" + teamName + "';")
	if err != nil {
		return
	}

	projects := projectsResult.([]interface{})
	if len(projects) == 0 {
		return
	}
	project := projects[0].(map[string]interface{})
	projectWsId := project["id"].(string)

	var wg sync.WaitGroup
	for _, fileId := range fileIds {
		// Upload files in parallel
		wg.Add(1)
		go func(fileId string, wg *sync.WaitGroup) {
			defer wg.Done()
			fileInfo, fileErr := p.API.GetFileInfo(fileId)
			if fileErr != nil {
				return
			}
			fileBody, fileErr := p.API.GetFile(fileId)
			base64Data := b64.RawStdEncoding.EncodeToString(fileBody)
			if fileErr != nil {
				return
			}
			documentData := map[string]interface{}{
				"notes_title":       "MM Uplaod",
				"reltoproject":      projectWsId,
				"filelocationtype":  "I",
				"filestatus":        1,
				"filedownloadcount": 0,
				"filename": map[string]interface{}{
					"size":    fileInfo.Size,
					"name":    fileInfo.Name,
					"type":    fileInfo.MimeType,
					"content": base64Data,
				},
			}

			cbContext.DoCreate("Documents", documentData)
		}(fileId, &wg)
	}
	// Wait for all uploads to finish
	wg.Wait()
}
