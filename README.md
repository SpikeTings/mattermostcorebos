# Mattermost-coreBOS plugin

Used to establish a bidirectional communication and synchronization between Mattermost and coreBOS. This plugin will listen on a set of trigger words that start with a hash sign, send requests to coreBOS and wait for the response which will be forwarded to the user. This plugin also listens for messages initiated from the related coreBOS installation.

## Installation

Steps to install the coreBOS plugin in Mattermost

1. Get the zip or git clone the repo from https://github.com/SpikeTings/mattermostcorebos

2. Install the dependencies: `go mod tidy`

3. Inside the project directory, build the plugin by running: `go build -o main main.go`

4. Create the tar.gz plugin file: `tar -czvf plugin.tar.gz main plugin.yaml`

5. Upload the plugin to Mattermost in System Console Plugin Management

## API

### health

plugin status check endpoint. Will answer with a JSON object like this

```json
{
  "success": true,
  "data": {
    "information": "Health check",
    "message": "The plugin status is active",
    "status": 200
  }
}
```

### key

another health check endpoint which executes some internal writes and reads to the mattermos key-value store (database). Will responde with:

```text
mmcb
mmcb
Hello, world!No Hash 
mmcb
With Hash
[109 109 99 98]
```

### syncuser

A POST request to this endpoint with information of the user we want to synchronize with that looks like this

```json
{
  "Username": "someusername",
  "Password": "thatuserspassword",
  "Email": "thatusersemail",
  "FirstName": "thatusersfirstname",
  "LastName": "thatuserslastname",
  "Position": "",
  "Roles": "system_user",
  "TeamNames": "TEAMNAME"
}
```

this endpoint will search for a user with the given username or email. If found, the mattermost user ID is returned. If not found, a new user will be created using the given information and the mattermost ID of the new user will be returned.

There are some error responses that may be returned.

### postmessage

this endpoint is to receive messages from the coreBOS installation

### team

For project management functionality.

- `/team/{team-name}/project/{name}/documents` get a list of documents related to the project
- `/team/{team-name}/project/{name}/task` create a task for the project
- `/team/{team-name}/project/{name}/method/{method}/module/{module}/invoke`
- `/team/{team-name}/project/{name}/wiki` POST: create a Conversation record
- `/team/{team-name}/project/{name}/wiki` PUT: update a Conversation record
- `/team/{team-name}/project/{name}/wiki` GET: retrieve Conversation records related to the project

## Code

- configuration directory holds the structure of the configuration settings of the plugin
- helpers directory holds two scripts that define a set of auxiliary functions to help with the functionality
- corebos/invoke is an abstraction proxy on top of the coreBOS golang library to login and invoke
- entity holds the POST and USER objects
- middleware/authentication is an authentication validator
- server directory holds the functionality of the plugin
  - api is the router of the actions
  - hook configures the different mattermost listeners we use
  - plugin sends and recieves messages between the applications
