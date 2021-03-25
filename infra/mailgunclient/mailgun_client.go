package mailgunclient

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

// MailgunClient interface
type MailgunClient interface {
	Welcome(subject, text, to, htmlStr string) error
	ResetPassword(subject, text, to, htmlStr, token string) error
}

type mailgunClient struct {
	//conf   configs.Config
	client *mailgun.MailgunImpl
}

// NewMailgunClient creates new Mailgun client given config
func NewMailgunClient() MailgunClient {
	return &mailgunClient{
		//conf:   c, // viper ile doğrudan alınabilir
		client: mailgun.NewMailgun("example-domain", "api-key"),
	}
}

func (mg *mailgunClient) Welcome(subject, text, to, htmlStr string) error {
	message := mg.createNewMessage(
		"email",
		subject,
		text,
		to,
		htmlStr,
	)

	ctx, cancel := mg.setContext(10)
	defer cancel()
	return mg.send(ctx, message)
}

func (mg *mailgunClient) ResetPassword(subject, text, to, htmlStr, token string) error {
	v := url.Values{}
	v.Set("token", token)

	resetURL := mg.getURL() + "/api/update_password?" + v.Encode()
	resetText := fmt.Sprintf(text, resetURL, token)
	resetHTML := fmt.Sprintf(htmlStr, resetURL, token)
	message := mg.createNewMessage(
		"email",
		subject,
		resetText,
		to,
		resetHTML,
	)

	ctx, cancel := mg.setContext(10)
	defer cancel()
	return mg.send(ctx, message)
}

// ========= Private methods =========

func (mg *mailgunClient) getURL() string {
	url := "host" + ":" + "port"
	return url
}

func (mg *mailgunClient) setContext(seconds time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Second*seconds)
}

func (mg *mailgunClient) createNewMessage(from, subject, text, to, htmlStr string) *mailgun.Message {
	message := mg.client.NewMessage(
		from,
		subject,
		text,
		to,
	)
	message.SetHtml(htmlStr)
	return message
}

func (mg *mailgunClient) send(ctx context.Context, message *mailgun.Message) error {
	_, _, err := mg.client.Send(ctx, message)
	if err != nil {
		return err
	}
	return nil
}
