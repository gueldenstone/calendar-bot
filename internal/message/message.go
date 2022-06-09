package message

import (
	"bytes"
	"fmt"
	"text/template"
	"time"

	"github.com/emersion/go-ical"
	"maunium.net/go/mautrix/event"
)

const (
	timeLayout = "15:04"
)

type Event struct {
	Summary     string
	StartTime   string
	EndTime     string
	Description string
}

type TemplatedMessage struct {
	Events       []Event
	htmlTemplate *template.Template
	txtTemplate  *template.Template
}

func NewTemplatedMessage(htmlTemplate, txtTemplate string, events []ical.Event, tz *time.Location) (TemplatedMessage, error) {
	msg := TemplatedMessage{}
	for _, evt := range events {
		prop := evt.Props.Get(ical.PropSummary)
		var summary string
		if prop != nil {
			summary = prop.Value
		} else {
			return msg, fmt.Errorf("no summary for event %s", evt.Name)
		}
		startTime, err := evt.Props.DateTime(ical.PropDateTimeStart, tz)
		if err != nil {
			return msg, err
		}
		endTime, err := evt.Props.DateTime(ical.PropDateTimeEnd, tz)
		if err != nil {
			return msg, err
		}
		var description string
		prop = evt.Props.Get(ical.PropDescription)
		if prop != nil {
			description = prop.Value
		} else {
			description = ""
		}
		evt := Event{
			Summary:     summary,
			StartTime:   startTime.Format(timeLayout),
			EndTime:     endTime.Format(timeLayout),
			Description: description,
		}
		if startTime.Hour() == 0 && startTime.Minute() == 0 {
			evt.StartTime = ""
		}
		if endTime.Hour() == 0 && endTime.Minute() == 0 {
			evt.EndTime = ""
		}

		msg.Events = append(msg.Events, evt)
	}
	htmlTmpl, err := template.ParseFiles(htmlTemplate)
	if err != nil {
		return msg, err
	}
	txtTmpl, err := template.ParseFiles(txtTemplate)
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
		MsgType:       event.MsgText,
		Body:          txt,
		Format:        event.FormatHTML,
		FormattedBody: html,
	}, nil
}
