package notifier

import (
	"fmt"

	"github.com/nikoksr/notify/service/telegram"

	"github.com/filippofinke/docker-events/internal/config"
)

func (n *notifierImpl) addTelegram(cfg config.TelegramConfig) error {
	service, err := telegram.New(cfg.Token)
	if err != nil {
		return fmt.Errorf("create telegram service: %w", err)
	}
	service.SetParseMode("")
	service.AddReceivers(cfg.ChatIDs...)
	n.client.UseServices(service)
	return nil
}
