package message

import (
	"bytes"
	"strings"
	"text/template"
	"time"

	"github.com/gueldenstone/calendar-bot/pkg/calendar"

	"maunium.net/go/mautrix/event"
)

const (
	timeLayout = "15:04"
)

type Event struct {
	Summary         string
	StartTime       string
	EndTime         string
	HtmlDescription string
	TxtDescription  string
}

type TemplatedMessage struct {
	Events       []Event
	htmlTemplate *template.Template
	txtTemplate  *template.Template
}

func NewTemplatedMessage(htmlTemplate, txtTemplate string, events []calendar.EventData, tz *time.Location) (TemplatedMessage, error) {
	msg := TemplatedMessage{}
	for _, evt := range events {
		evt := Event{
			Summary:         evt.Summary,
			StartTime:       evt.Start.In(tz).Format(timeLayout),
			EndTime:         evt.End.In(tz).Format(timeLayout),
			HtmlDescription: strings.ReplaceAll(strings.ReplaceAll(evt.Description, "\\n", "<br>"), "\\", ""),
			TxtDescription:  evt.Description,
		}

		if evt.Start.Hour() == 0 && evt.Start.Minute() == 0 {
			evt.StartTime = ""
		}
		if evt.End.Hour() == 0 && evt.End.Minute() == 0 {
			evt.EndTime = ""
		}

		msg.Events = append(msg.Events, evt)
	}

	funcMap := template.FuncMap{
		"today": func() string {
			return time.Now().Format("Monday 02.01.2006")
		},
	}
	htmlTmpl, err := template.New("event.html").Funcs(funcMap).ParseFiles(htmlTemplate)
	if err != nil {
		return msg, err
	}
	txtTmpl, err := template.New("event.txt").Funcs(funcMap).ParseFiles(txtTemplate)
	if err != nil {
		return msg, err
	}
	msg.htmlTemplate = htmlTmpl
	msg.txtTemplate = txtTmpl
	return msg, nil
}

func (t TemplatedMessage) RenderHtml() (string, error) {
	buf := bytes.Buffer{}
	err := t.htmlTemplate.Execute(&buf, t)
	return buf.String(), err
}

func (t TemplatedMessage) RenderTxt() (string, error) {
	buf := bytes.Buffer{}
	err := t.txtTemplate.Execute(&buf, t)
	return buf.String(), err
}
func (t TemplatedMessage) Render() (html string, txt string, err error) {
	html, err = t.RenderHtml()
	if err != nil {
		return "", "", err
	}
	txt, err = t.RenderTxt()
	if err != nil {
		return "", "", err
	}
	return
}

func (t TemplatedMessage) MatrixMessage() (event.MessageEventContent, error) {
	html, err := t.RenderHtml()
	if err != nil {
		return event.MessageEventContent{}, err
	}
	txt, err := t.RenderTxt()
	if err != nil {
		return event.MessageEventContent{}, err
	}
	return event.MessageEventContent{
		MsgType:       event.MsgNotice,
		Body:          txt,
		Format:        event.FormatHTML,
		FormattedBody: html,
	}, nil
}
