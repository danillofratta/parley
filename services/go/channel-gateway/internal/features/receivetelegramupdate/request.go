package receivetelegramupdate

// Request is the subset of Telegram's Update object this slice reads.
// Provider types stay here: nothing past the endpoint sees them (anticorruption layer).
type Request struct {
	UpdateID int64           `json:"update_id"`
	Message  *RequestMessage `json:"message"`
}

type RequestMessage struct {
	Chat RequestChat `json:"chat"`
	Text string      `json:"text"`
}

type RequestChat struct {
	ID int64 `json:"id"`
}
