package receivetelegramupdate

// ReceiveTelegramUpdateRequest is the subset of Telegram's Update object this
// slice reads. Provider types stay here: nothing past the endpoint sees them.
type ReceiveTelegramUpdateRequest struct {
	UpdateID int64                         `json:"update_id"`
	Message  *ReceiveTelegramUpdateMessage `json:"message"`
}

type ReceiveTelegramUpdateMessage struct {
	Chat ReceiveTelegramUpdateChat `json:"chat"`
	Text string                    `json:"text"`
}

type ReceiveTelegramUpdateChat struct {
	ID int64 `json:"id"`
}
