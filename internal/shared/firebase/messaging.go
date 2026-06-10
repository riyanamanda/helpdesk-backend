package firebase

import (
	"context"
	"log/slog"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

type FCMSender interface {
	SendMulticast(ctx context.Context, tokens []string, data map[string]string) error
}

type fcmSender struct {
	client *messaging.Client
}

type noopFCMSender struct{}

func (n *noopFCMSender) SendMulticast(_ context.Context, _ []string, _ map[string]string) error {
	return nil
}

func NewFCMSender(ctx context.Context, projectID, credentialsJSON string) (FCMSender, error) {
	var opts []option.ClientOption
	if credentialsJSON != "" {
		opts = append(opts, option.WithCredentialsJSON([]byte(credentialsJSON)))
	}

	app, err := firebase.NewApp(ctx, &firebase.Config{ProjectID: projectID}, opts...)
	if err != nil {
		return nil, err
	}

	client, err := app.Messaging(ctx)
	if err != nil {
		return nil, err
	}

	return &fcmSender{client: client}, nil
}

func NewNoopFCMSender() FCMSender {
	return &noopFCMSender{}
}

func (f *fcmSender) SendMulticast(ctx context.Context, tokens []string, data map[string]string) error {
	if len(tokens) == 0 {
		return nil
	}

	msg := &messaging.MulticastMessage{
		Tokens: tokens,
		Data:   data,
	}

	br, err := f.client.SendEachForMulticast(ctx, msg)
	if err != nil {
		return err
	}

	if br.FailureCount > 0 {
		for i, result := range br.Responses {
			if !result.Success {
				slog.WarnContext(ctx, "fcm: failed to send to token", "token_index", i, "error", result.Error)
			}
		}
	}

	return nil
}
