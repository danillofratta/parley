package events

import (
	"strings"
	"time"

	"github.com/parley/parley/services/go/channel-gateway/internal/domain/valueobjects"
)

type MessageReceived struct {
	EventBase
	inboundMessageID  string
	tenantID          string
	sender            valueobjects.ContactAddress
	providerMessageID string
	text              string
	receivedAtUnix    int64
}

func NewMessageReceived(
	inboundMessageID string,
	tenantID string,
	sender valueobjects.ContactAddress,
	providerMessageID string,
	text string,
	receivedAt time.Time,
) MessageReceived {
	if receivedAt.IsZero() {
		receivedAt = time.Now().UTC()
	}

	return MessageReceived{
		EventBase:          NewEventBase("MessageReceived.v1", time.Now().UTC().Unix()),
		inboundMessageID:   strings.TrimSpace(inboundMessageID),
		tenantID:           strings.TrimSpace(tenantID),
		sender:             sender,
		providerMessageID:  strings.TrimSpace(providerMessageID),
		text:               text,
		receivedAtUnix:     receivedAt.UTC().Unix(),
	}
}

func (e MessageReceived) InboundMessageID() string { return e.inboundMessageID }
func (e MessageReceived) TenantID() string         { return e.tenantID }
func (e MessageReceived) Sender() valueobjects.ContactAddress {
	return e.sender
}
func (e MessageReceived) ProviderMessageID() string { return e.providerMessageID }
func (e MessageReceived) Text() string              { return e.text }
func (e MessageReceived) ReceivedAtUnix() int64     { return e.receivedAtUnix }