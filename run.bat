@echo
cd /Users/arditsarja/go/src/mattermost-server-plugin
set GOOS=linux
go build -o main.exe main.go
tar -czvf plugin.tar.gz main plugin.yaml
