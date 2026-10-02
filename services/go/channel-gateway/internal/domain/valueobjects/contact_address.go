package valueobjects

import (
	"strings"

	"github.com/parley/parley/services/go/channel-gateway/internal/domain/enums"
	"github.com/parley/parley/services/go/channel-gateway/internal/domain/rules"
)

type ContactAddress struct {
	channel enums.Channel
	chat    string
}

func NewContactAddress(channel enums.Channel, chat string) (ContactAddress, error) {
	if err := rules.CheckRules(
		rules.ChannelMustBeKnown{Channel: channel},
		rules.ChatMustBeInformed{Chat: chat},
	); err != nil {
		return ContactAddress{}, err
	}

	return ContactAddress{
		channel: channel,
		chat:    chat,
	}, nil
}

func (a ContactAddress) Channel() enums.Channel {
	return a.channel
}

func (a ContactAddress) Chat() string {
	return a.chat
}

func (a ContactAddress) IsZero() bool {
	return a.channel.IsZero() && strings.TrimSpace(a.chat) == ""
}