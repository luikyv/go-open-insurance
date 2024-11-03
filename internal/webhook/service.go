package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/luikyv/go-oidc/pkg/goidc"
	"github.com/luikyv/go-oidc/pkg/provider"
	"github.com/luikyv/go-open-insurance/internal/api"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const webhookBasePath string = "/open-insurance/webhook/v1"

type Service struct {
	op             provider.Provider
	httpClientFunc goidc.HTTPClientFunc
}

func NewService(op provider.Provider, httpClientFunc goidc.HTTPClientFunc) Service {
	return Service{
		op:             op,
		httpClientFunc: httpClientFunc,
	}
}

func (s Service) Notify(ctx context.Context, clientID, endpointPath string) {
	baseWebhookURL, err := s.clientBaseWebhookURL(ctx, clientID)
	if err != nil {
		return
	}

	webhookURL, _ := url.JoinPath(
		baseWebhookURL,
		webhookBasePath,
		endpointPath,
	)
	go func() {
		time.Sleep(10 * time.Second)
		s.notify(ctx, webhookURL)
	}()
}

func (s Service) clientBaseWebhookURL(
	ctx context.Context,
	clientID string,
) (
	string,
	error,
) {
	client, err := s.op.Client(ctx, clientID)
	if err != nil {
		api.Logger(ctx).Error("could not fetch the client to post webhook",
			slog.String("client_id", clientID))
		return "", err
	}

	rawWebhookURIs := client.Attribute("webhook_uris")
	if rawWebhookURIs == nil {
		api.Logger(ctx).Info("client does not have webhook uris defined")
		return "", errors.New("the client has no webhook uris defined")
	}

	webhookURIs := rawWebhookURIs.(primitive.A)
	if len(webhookURIs) == 0 {
		api.Logger(ctx).Info("client has 0 webhook uris defined")
		return "", errors.New("client has 0 webhook uris defined")
	}

	return webhookURIs[0].(string), nil
}

func (s Service) notify(ctx context.Context, webhookURL string) {
	body, _ := json.Marshal(api.WebhookRequest{
		Data: struct {
			Timestamp api.DateTime "json:\"timestamp\""
		}{
			Timestamp: api.DateTimeNow(),
		},
	})
	req, _ := http.NewRequest(http.MethodPost, webhookURL, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Webhook-Interaction-ID", uuid.NewString())

	api.Logger(ctx).Info("sending webhook notification", slog.String("url", webhookURL))
	httpClient := s.httpClientFunc(ctx)
	resp, err := httpClient.Do(req)
	if err != nil {
		api.Logger(ctx).Info("error sending webhook notification",
			slog.String("error", err.Error()))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		api.Logger(ctx).Info("error sending the webhook notification",
			slog.Int("status_code", resp.StatusCode))
	}
}
