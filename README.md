# Calendar-Bot

![test workflow](https://github.com/gueldenstone/calendar-bot/actions/workflows/go.yml/badge.svg)
![build](https://github.com/gueldenstone/calendar-bot/actions/workflows/docker-publish.yml/badge.svg)
[<img src="https://img.shields.io/badge/dockerhub-image-blue.svg?logo=Docker">](https://hub.docker.com/r/gueldenstone/calendar-bot)

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
