package server

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"mattermostcorebos/helpers"
	"net/http"
)

func (p *Plugin) Health(w http.ResponseWriter, r *http.Request) {
	helpers.DisplayAppSuccessResponse(w, "Health check", "The plugin status is active")
}

func GetKeyHash(key string) string {
	hash := sha256.New()
	_, _ = hash.Write([]byte(key))
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}
func (p *Plugin) DoKeyJob(w http.ResponseWriter, r *http.Request) {
	key := "name"

	value := "mmcb"

	fmt.Fprint(w, value+"\n")
	err := p.API.KVSet(key, []byte(value))

	if err != nil {
		fmt.Fprint(w, "1"+err.Error())
		return
	}

	fmt.Fprint(w, value+"\n")
	err = p.API.KVSet(GetKeyHash(key), []byte(value))

	if err != nil {
		fmt.Fprint(w, "2"+err.Error())
		return
	}

	returnNoHash, err := p.API.KVGet(key)

	if err != nil {
		fmt.Fprint(w, "3"+err.Error())
		return
	}
	returnHash, err := p.API.KVGet(GetKeyHash(key))

	if err != nil {
		fmt.Fprint(w, "4"+err.Error())
		return
	}

	hello := fmt.Sprint("Hello, world!", "No Hash \n", string(returnNoHash), "\nWith Hash\n", returnHash)
	fmt.Fprint(w, hello)
}
