package rules

import "strings"

// zeroable is satisfied by value objects that can be empty, without importing them.
type zeroable interface {
	IsZero() bool
}

// SenderMustBeInformed: every inbound message comes from a contact address.
type SenderMustBeInformed struct{ Sender zeroable }

func (r SenderMustBeInformed) IsBroken() bool { return r.Sender == nil || r.Sender.IsZero() }
func (SenderMustBeInformed) Code() string     { return "inbound_message.sender_required" }
func (SenderMustBeInformed) Message() string  { return "an inbound message must have a sender" }

// ProviderMessageMustBeIdentified: the provider's id is what lets Parley accept a message only once.
type ProviderMessageMustBeIdentified struct{ ProviderMessageID string }

func (r ProviderMessageMustBeIdentified) IsBroken() bool {
	return strings.TrimSpace(r.ProviderMessageID) == ""
}
func (ProviderMessageMustBeIdentified) Code() string {
	return "inbound_message.provider_message_id_required"
}
func (ProviderMessageMustBeIdentified) Message() string {
	return "an inbound message must carry the provider's message id"
}

// MessageMustHaveText: in Release 1 Parley accepts text messages only.
type MessageMustHaveText struct{ Text string }

func (r MessageMustHaveText) IsBroken() bool { return strings.TrimSpace(r.Text) == "" }
func (MessageMustHaveText) Code() string     { return "inbound_message.text_required" }
func (MessageMustHaveText) Message() string  { return "an inbound message must have text" }
