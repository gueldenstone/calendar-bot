##
## Build
##
FROM golang:1.23-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY cmd ./cmd/
COPY internal ./internal
COPY pkg ./pkg

RUN go build -o /calendar-bot ./cmd/calendar-bot

##
## Deploy
##
FROM alpine

WORKDIR /
RUN apk add --no-cache tzdata
RUN ln -s /usr/share/zoneinfo/Europe/Berlin /etc/localtime
RUN adduser -D calendar-bot
USER calendar-bot
COPY --from=build /calendar-bot /calendar-bot
COPY ./templates /templates

ENTRYPOINT ["/calendar-bot", "-html", "/templates/event.html", "-txt", "/templates/event.txt"]
