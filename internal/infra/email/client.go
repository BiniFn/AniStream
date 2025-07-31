package email

import (
	"context"
	"errors"
	"fmt"

	"github.com/resend/resend-go/v2"
)

type SendSimpleEmailParams struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Html    string   `json:"html,omitempty"`
}

func (p SendSimpleEmailParams) Validate() error {
	if p.From == "" {
		return errors.New("from email is required")
	}
	if len(p.To) == 0 {
		return errors.New("at least one recipient is required")
	}
	if p.Subject == "" {
		return errors.New("subject is required")
	}
	return nil
}

type EmailClient interface {
	SendSimpleEmail(ctx context.Context, params SendSimpleEmailParams) error
}

type Client struct {
	resend           *resend.Client
	defaultFromEmail string
}

func NewClient(apiKey string, defaultFromEmail string) EmailClient {
	return &Client{
		resend:           resend.NewClient(apiKey),
		defaultFromEmail: defaultFromEmail,
	}
}

func (c *Client) SendSimpleEmail(ctx context.Context, params SendSimpleEmailParams) error {
	if params.From == "" {
		params.From = c.defaultFromEmail
	}

	if err := params.Validate(); err != nil {
		return fmt.Errorf("invalid parameters: %v", err)
	}

	_, err := c.resend.Emails.SendWithContext(ctx, &resend.SendEmailRequest{
		From:    params.From,
		To:      params.To,
		Subject: params.Subject,
		Html:    params.Html,
	})

	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}
