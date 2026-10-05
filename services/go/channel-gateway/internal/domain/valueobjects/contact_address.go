package valueobjects

import (
	"github.com/danillofratta/parley/services/go/channel-gateway/internal/domain/enums"
	"github.com/danillofratta/parley/services/go/channel-gateway/internal/domain/rules"
	"github.com/danillofratta/parley/services/go/channel-gateway/internal/domain/seedwork"
)

// ContactAddress is where a contact can be reached: a channel and the contact's chat on it.
// Replies go back to the same address. Immutable and compared by value.
type ContactAddress struct {
	channel enums.Channel
	chat    string
}

func NewContactAddress(channel enums.Channel, chat string) (ContactAddress, error) {
	if err := seedwork.CheckRules(
		rules.ChannelMustBeKnown{Channel: channel},
		rules.ChatMustBeInformed{Chat: chat},
	); err != nil {
		return ContactAddress{}, err
	}
	return ContactAddress{channel: channel, chat: chat}, nil
}

func (a ContactAddress) Channel() enums.Channel { return a.channel }
func (a ContactAddress) Chat() string           { return a.chat }

// IsZero reports whether the address was never created through NewContactAddress.
func (a ContactAddress) IsZero() bool { return a == (ContactAddress{}) }
