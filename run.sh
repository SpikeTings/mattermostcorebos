cd /Users/arditsarja/go/src/mattermost-server-plugin
rm plugin.exe
rm plugin.tar.gz
export GOOS=linux
go build -o main main.go
tar -czvf plugin.tar.gz main plugin.yaml