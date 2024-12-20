package notifications

import (
	"bytes"
	"fmt"
	"text/template"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"

	"tower-defense-api/lib/models"
)

type SendGridMailer struct {
	fromEmail string
	apiKey    string
	client    *sendgrid.Client
}

func NewSendgrid(apiKey, fromEmail string) *SendGridMailer {
	client := sendgrid.NewSendClient(apiKey)

	return &SendGridMailer{
		fromEmail: fromEmail,
		apiKey:    apiKey,
		client:    client,
	}
}

func (m *SendGridMailer) Send(message *models.Message, isSandbox bool) (int, error) {
	from := mail.NewEmail(FromName, "td.mazing@gmail.com")
	to := mail.NewEmail("Test Sender", "td.mazing@gmail.com")

	// template parsing and building
	tmpl, err := template.ParseFS(FS, "templates/"+MessageTemplate)
	if err != nil {
		return -1, err
	}

	vars := struct {
		UserId  int64
		Content string
		Sender  string
	}{
		UserId:  message.UserID,
		Content: message.Content,
		Sender:  message.Sender,
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", vars)
	if err != nil {
		return -1, err
	}

	body := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(body, "body", vars)
	if err != nil {
		return -1, err
	}

	email := mail.NewSingleEmail(from, subject.String(), to, "", body.String())

	email.SetMailSettings(&mail.MailSettings{
		SandboxMode: &mail.Setting{
			Enable: &isSandbox,
		},
	})

	var retryErr error
	for i := 0; i < maxRetires; i++ {
		response, retryErr := m.client.Send(email)
		if retryErr != nil {
			// exponential backoff
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}

		return response.StatusCode, nil
	}

	return -1, fmt.Errorf("failed to send email after %d attempt, error: %v", maxRetires, retryErr)
}
