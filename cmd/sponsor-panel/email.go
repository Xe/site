package main

import (
	"context"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

// EmailSender sends emails.
type EmailSender interface {
	SendEmail(ctx context.Context, to, subject, htmlBody string) error
}

// sesEmailSender sends emails via AWS SES.
type sesEmailSender struct {
	client *ses.Client
	from   string
}

func newSESEmailSender(ctx context.Context, region, from string) (*sesEmailSender, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, err
	}
	return &sesEmailSender{
		client: ses.NewFromConfig(cfg),
		from:   from,
	}, nil
}

func (s *sesEmailSender) SendEmail(ctx context.Context, to, subject, htmlBody string) error {
	_, err := s.client.SendEmail(ctx, &ses.SendEmailInput{
		Source: &s.from,
		Destination: &types.Destination{
			ToAddresses: []string{to},
		},
		Message: &types.Message{
			Subject: &types.Content{Data: &subject},
			Body: &types.Body{
				Html: &types.Content{Data: &htmlBody},
			},
		},
	})
	return err
}

// logEmailSender logs emails instead of sending them (for development).
type logEmailSender struct{}

func (l *logEmailSender) SendEmail(ctx context.Context, to, subject, htmlBody string) error {
	slog.Info("logEmailSender: email",
		"to", to,
		"subject", subject,
		"body", htmlBody,
	)
	return nil
}
