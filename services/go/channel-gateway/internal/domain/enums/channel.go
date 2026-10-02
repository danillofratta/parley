package enums

import "strings"

type Channel int

const (
	ChannelUnknown Channel = iota
	ChannelTelegram
)

var channelNames = map[Channel]string{
	ChannelUnknown:  "Unknown",
	ChannelTelegram:   "Telegram",
}

func (c Channel) String() string {
	if name, ok := channelNames[c]; ok {
		return name
	}
	return channelNames[ChannelUnknown]
}

func (c Channel) IsKnown() bool {
	return c != ChannelUnknown
}

func (c Channel) IsZero() bool {
	return c == ChannelUnknown
}

func ParseChannel(name string) Channel {
	clean := strings.TrimSpace(strings.ToLower(name))
	for channel, channelName := range channelNames {
		if strings.ToLower(channelName) == clean {
			return channel
		}
	}
	return ChannelUnknown
}