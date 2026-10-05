package receivetelegramupdate

// Command is the application input, independent of HTTP and of Telegram's format.
type Command struct {
	TenantID          string
	ProviderMessageID string
	ExternalChatID    string
	Text              string
}
