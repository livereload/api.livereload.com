package main

import (
	"strings"

	"github.com/keighl/postmark"
)

func sendEmail(from, to, replyTo, subject, body, tag string, track bool) error {
	email := postmark.Email{
		From:    from,
		To:      to,
		Subject: subject,
		// HtmlBody:   "...",
		TextBody:   body,
		Tag:        tag,
		TrackOpens: true,
	}
	if replyTo != "" {
		email.ReplyTo = replyTo
	}

	_, err := postmarkClient.SendEmail(email)
	return err
}

func applyEmailTemplate(template string, params map[string]string) (subject string, body string) {
	lines := strings.SplitN(template, "\n", 3)
	if len(lines) < 3 || !strings.HasPrefix(lines[0], "Subject: ") || lines[1] != "" {
		panic("Invalid email template format")
	}

	subject = replaceStrings(strings.Replace(lines[0], "Subject: ", "", 1), params)
	body = replaceStrings(lines[2], params)
	return
}

func replaceStrings(s string, params map[string]string) string {
	for k, v := range params {
		s = strings.Replace(s, k, v, -1)
	}
	return s
}
