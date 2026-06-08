package notifier

import (
	"fmt"
	stdhttp "net/http"
	"strings"

	"github.com/nikoksr/notify/service/discord"
	notifyhttp "github.com/nikoksr/notify/service/http"

	"github.com/filippofinke/docker-events/internal/config"
)

func (n *notifierImpl) addDiscord(cfg config.DiscordConfig) error {
	// Setup Discord bot if token is provided
	if strings.TrimSpace(cfg.Token) != "" {
		if len(cfg.ChannelIDs) == 0 {
			return fmt.Errorf("discord bot configured but no channels specified")
		}
		service := discord.New()
		if err := service.AuthenticateWithBotToken(cfg.Token); err != nil {
			return fmt.Errorf("authenticate discord bot: %w", err)
		}
		service.AddReceivers(cfg.ChannelIDs...)
		n.client.UseServices(service)
	}

	// Setup Discord webhooks if URLs are provided using notify's http service
	if len(cfg.WebhookURLs) > 0 {
		httpService := notifyhttp.New()

		for _, url := range cfg.WebhookURLs {
			httpService.AddReceivers(&notifyhttp.Webhook{
				URL:         url,
				Header:      stdhttp.Header{},
				ContentType: "application/json",
				Method:      stdhttp.MethodPost,
				BuildPayload: func(subject, message string) (payload any) {
					return map[string]any{
						"content": fmt.Sprintf("**%s**\n%s", subject, message),
						"embeds":  map[string]any{},
					}
				},
			})
		}

		n.client.UseServices(httpService)
	}

	if strings.TrimSpace(cfg.Token) == "" && len(cfg.WebhookURLs) == 0 {
		return fmt.Errorf("discord enabled but no bot token or webhook URLs configured")
	}

	return nil
}
