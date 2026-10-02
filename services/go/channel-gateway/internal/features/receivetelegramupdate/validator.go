package receivetelegramupdate

import (
	"errors"
	"fmt"
	"strings"
)

var ErrInvalidCommand = errors.New("invalid receive telegram update command")

// validate checks the boundary (required fields). Business invariants live in
// the aggregate and value objects.
func validate(cmd ReceiveTelegramUpdateCommand) error {
	var missing []string

	if strings.TrimSpace(cmd.TenantID) == "" {
		missing = append(missing, "TenantID")
	}
	if strings.TrimSpace(cmd.ProviderMessageID) == "" {
		missing = append(missing, "ProviderMessageID")
	}
	if strings.TrimSpace(cmd.ExternalChatID) == "" {
		missing = append(missing, "ExternalChatID")
	}
	if strings.TrimSpace(cmd.Text) == "" {
		missing = append(missing, "Text")
	}
	if len(missing) > 0 {
		return fmt.Errorf("%w: missing %s", ErrInvalidCommand, strings.Join(missing, ", "))
	}

	return nil
}
