package receivetelegramupdate

// Result is the application output.
type Result struct {
	InboundMessageID string
	Duplicate        bool
}
