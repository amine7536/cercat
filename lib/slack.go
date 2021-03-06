package lib

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

// slackAttachmentField
type slackAttachmentField struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

// slackAttachment
type slackAttachment struct {
	Color  string                 `json:"color"`
	Text   string                 `json:"text,omitempty"`
	Fields []slackAttachmentField `json:"fields"`
	// Footer     string                 `json:"footer,omitempty"`
	// FooterIcon string                 `json:"footer_icon,omitempty"`
}

// SlackPayload represents a message to send to Slack
type SlackPayload struct {
	Text        string            `json:"text,omitempty"`
	Username    string            `json:"username,omitempty"`
	IconURL     string            `json:"icon_url,omitempty"`
	Attachments []slackAttachment `json:"attachments,omitempty"`
}

// NewSlackPayload generates a new Slack Payload
func NewSlackPayload(config *Configuration, r *Result) SlackPayload {
	var attachments []slackAttachment
	var attachment slackAttachment
	var fields []slackAttachmentField
	var field slackAttachmentField

	field.Title = "Domain"
	field.Value = r.Domain
	field.Short = true
	fields = append(fields, field)

	field.Title = "Issuer"
	field.Value = r.Issuer
	field.Short = true
	fields = append(fields, field)

	if r.IDN != "" {
		field.Title = "IDN"
		field.Value = r.IDN
		field.Short = true
		fields = append(fields, field)
	}

	field.Title = "SAN"
	field.Short = false
	field.Value = strings.Join(r.SAN, ", ")
	fields = append(fields, field)

	field.Title = "Addresses"
	field.Short = false
	field.Value = strings.Join(r.Addresses, ", ")
	fields = append(fields, field)

	attachment.Fields = fields

	attachment.Color = "#ff5400"

	attachments = append(attachments, attachment)

	domain := r.Domain
	if r.IDN != "" {
		domain += " (" + r.IDN + ")"
	}

	return SlackPayload{
		Text:        "A certificate for " + domain + " has been issued",
		Username:    config.SlackUsername,
		IconURL:     config.SlackIconURL,
		Attachments: attachments,
	}
}

// post posts to Slack a Payload
func (s SlackPayload) post(config *Configuration) {
	body, _ := json.Marshal(s)
	req, _ := http.NewRequest(http.MethodPost, config.SlackWebHookURL, bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	_, err := client.Do(req)
	if err != nil {
		log.Warn("Slack Post error")
	}
}
