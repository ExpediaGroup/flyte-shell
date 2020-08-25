# flyte-shell

![Build Status](https://travis-ci.org/ExpediaGroup/flyte-shell.svg?branch=master)
[![Docker Stars](https://img.shields.io/docker/stars/hotelsdotcom/flyte-shell.svg)](https://hub.docker.com/r/hotelsdotcom/flyte-shell)
[![Docker Pulls](https://img.shields.io/docker/pulls/hotelsdotcom/flyte-shell.svg)](https://hub.docker.com/r/hotelsdotcom/flyte-shell)

## Overview
The Shell pack can run an arbitrary bash script. If you can run it in a bash terminal then it'll run on this pack.
The pack has bash, curl, jq and openssh installed on its container. If you need something more, feel free
to make a pull request.

## Build & Run
### Command Line
To build and run from the command line:
* Clone this repo
* Run `dep ensure` (must have [dep](https://github.com/golang/dep) installed )
* Run `go build`
* Run `FLYTE_API_URL=http://.../ ./flyte-shell`, FLYTE_API_URL environment variable is required
* User can set labels (optional) using `FLYTE_LABELS` env. variable. Example: `FLYTE_LABELS='env=dev,user=root'`

### Docker
To build and run in docker:
* Run `docker build -t flyte-shell .`
* Run `docker run -e FLYTE_API_URL=http://.../ flyte-shell`
* The The FLYTE_API_URL environment variable needs to be set
* Pack labels (optional) can be set using FLYTE_LABELS env. variable

## Commands
### ```shell``` Command
#### Input
This command's input is the script to be run:
```
"input" : "echo -e 'world\nhello' | sort"
```
The script needs to be a one line bash command. Distinct commands should be separated by `;`. While newlines are
perfectly fine in the shell, they do make for invalid JSON and no escape expansion is performed, so stick with
single lines. The whole script argument is passed to `exec.Command()`. To test your command
on the command line run it using this: `bash -c '$script'`. **Single quotes are only required when running this command on
the command line, do not put single quotes around your script in the flow.**

#### Output Events
There is one type of output event, it contains the stdout, stderr and exit code.

For example, the script:
`"echo -e 'world\nhello' | sort"` would return the output:
```
"payload": {
    "stdout": "hello\nworld",
    "stderr": "",
    "exit_code": 0
}
```
While `echo oh ; >&2 echo dear` results in:
```
"payload: {
    "stdout": "oh",
    "stderr": "dear",
    "exit_code": 0
}
```
### Flows
An example flow, listening for a slack message and executing a command when one is received:
```
{
  "name": "slackShell",
  "description": "Run script when message received",
  "steps": [
      "event": {
        "name": "ReceivedMessage",
        "packName": "Slack"
      }, 
      "command": {
        "name": "Shell",
        "packName": "Shell",
        "input": "echo 'received an event from the Slack pack'"
      }
  ]
}
```

An example flow listening to shell output:
```
{
  "name": "shellOutput",
  "description": "Send message to slack when script has been run",
  "steps": [
      "event": {
        "name": "Output",
        "packName": "Shell"
      }, 
      "command": {
        "name": "SendMessage",
        "packName": "Slack",
        "input": {
          "channelId": "12345",
          "message": "Shell pack executed command. Output is: {{ Event.Payload.stdout }}"
        }
      }
  ]
}
```
