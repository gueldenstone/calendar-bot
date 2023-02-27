# Calendar-Bot

![test workflow](https://github.com/gueldenstone/calendar-bot/actions/workflows/go.yml/badge.svg)
![build](https://github.com/gueldenstone/calendar-bot/actions/workflows/docker-publish.yml/badge.svg)
[<img src="https://img.shields.io/badge/dockerhub-image-blue.svg?logo=Docker">](https://hub.docker.com/r/gueldenstone/calendar-bot)

This is a simple matrix bot that will monitor a calendar given by a URL and post events as `HTML` and `Plain text` according to the given templates. The bot is intended to be run in a docker imager, but feel free to build the the application like so:

`go build cmd/calendar-bot/calendar-bot.go`

## Usage

```yaml
---
version: "3.9"
services:
  calendar-bot:
      image: gueldenstone/calendar-bot:latest
      command: -config /config.yaml
      volumes:
        - "<path_to_config>:/config.yaml"
        - /etc/localtime:/etc/localtime:ro
```

## Example configuration

```yaml
# Configuration file for the calendar-bot
homeserver: matrix.org
rooms:
  - "#test-the-bot:matrix.org"

calendarURL: "<calendar_url>"

nofifyTime: "13:33"

username: "<matrix_username>"
password: "<supersecretpassword>"
```

# Tests

There are only some basic tests with the snapshot of my local [hackspaces](https://x-hain.de) calendar. You can run them like so:

`go test -v ./...`
