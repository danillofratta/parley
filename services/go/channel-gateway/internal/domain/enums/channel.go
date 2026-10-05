package enums

// Channel is the medium a contact uses to reach a tenant.
type Channel int

const (
	ChannelUnknown Channel = iota // zero value: never a valid channel
	ChannelTelegram
)

// channelNames holds the canonical lowercase names used in storage, events and keys.
var channelNames = map[Channel]string{
	ChannelTelegram: "telegram",
}

func (c Channel) String() string {
	if name, ok := channelNames[c]; ok {
		return name
	}
	return "unknown"
}

// IsKnown reports whether Parley supports the channel.
func (c Channel) IsKnown() bool {
	_, ok := channelNames[c]
	return ok
}

// ParseChannel converts a stored or transported name into a Channel.
func ParseChannel(name string) (Channel, bool) {
	for channel, n := range channelNames {
		if n == name {
			return channel, true
		}
	}
	return ChannelUnknown, false
}
