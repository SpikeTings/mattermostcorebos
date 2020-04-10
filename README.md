# Mattermost server plugin used to communicate and synchronize with coreBOS

Steps to install corebos plugin in Mattermost

1. Get the zip or git clone the repo from https://github.com/SpikeTings/mattermost-server-plugin

2. Unzip it in the $GOPATH/src directory

3. Install following dependencies:
   - go get -u github.com/mattermost/mattermost-server
   - go get -u github.com/gorilla/mux
   - go get -u github.com/tsolucio/corebosgowslib
   - go get -u github.com/thedevsaddam/govalidator

4. Inside the project directory, build the plugin by running the following command (linux example): go build -o main main.go

5. Create the tar.gz plugin (linux example): tar -czvf plugin.tar.gz main plugin.yaml

6. Upload the plugin to Mattermost in System Console Plugin Management