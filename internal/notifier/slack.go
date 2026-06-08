package notifier

import (
	"fmt"
	"strings"

	"github.com/nikoksr/notify/service/slack"

	"github.com/filippofinke/docker-events/internal/config"
)

func (n *notifierImpl) addSlack(cfg config.SlackConfig) error {
	if strings.TrimSpace(cfg.Token) == "" {
		return fmt.Errorf("empty slack token")
	}
	if len(cfg.Channels) == 0 {
		return fmt.Errorf("no slack channels configured")
	}
	service := slack.New(cfg.Token)
	service.AddReceivers(cfg.Channels...)
	n.client.UseServices(service)
	return nil
}
