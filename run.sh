cd /Users/arditsarja/go/src/mattermost-server-plugin
export GOOS=linux
go build -o main main.go
tar -czvf plugin.tar.gz main plugin.yaml