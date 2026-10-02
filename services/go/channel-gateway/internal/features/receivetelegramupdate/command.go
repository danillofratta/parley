package receivetelegramupdate

type ReceiveTelegramUpdateCommand struct {
	TenantID          string
	ProviderMessageID string
	ExternalChatID    string
	Text              string
}
