package rules

import (
	"strings"

	"github.com/danillofratta/parley/services/go/channel-gateway/internal/domain/enums"
)

// ChannelMustBeKnown: a contact can only be reached through a channel Parley supports.
type ChannelMustBeKnown struct{ Channel enums.Channel }

func (r ChannelMustBeKnown) IsBroken() bool { return !r.Channel.IsKnown() }
func (ChannelMustBeKnown) Code() string     { return "contact_address.unknown_channel" }
func (ChannelMustBeKnown) Message() string  { return "the channel is not supported" }

// ChatMustBeInformed: the contact's chat on the channel is needed to reply.
type ChatMustBeInformed struct{ Chat string }

func (r ChatMustBeInformed) IsBroken() bool { return strings.TrimSpace(r.Chat) == "" }
func (ChatMustBeInformed) Code() string     { return "contact_address.chat_required" }
func (ChatMustBeInformed) Message() string  { return "the contact's chat is required" }
