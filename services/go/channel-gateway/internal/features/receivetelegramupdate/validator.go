package receivetelegramupdate

import (
	"errors"
	"fmt"
	"strings"
)

var ErrInvalidCommand = errors.New("invalid receive telegram update command")

// validate checks the shape of the command (transport fields present).
// Business invariants, such as "a message must have text", live in the domain.
func validate(cmd Command) error {
	var missing []string
	if strings.TrimSpace(cmd.TenantID) == "" {
		missing = append(missing, "tenant id")
	}
	if strings.TrimSpace(cmd.ProviderMessageID) == "" {
		missing = append(missing, "provider message id")
	}
	if strings.TrimSpace(cmd.ExternalChatID) == "" {
		missing = append(missing, "external chat id")
	}
	if len(missing) > 0 {
		return fmt.Errorf("%w: missing %s", ErrInvalidCommand, strings.Join(missing, ", "))
	}
	return nil
}
