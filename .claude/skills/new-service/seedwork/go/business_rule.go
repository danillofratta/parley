package seedwork

// BusinessRule is an invariant named in the ubiquitous language (e.g. MessageMustHaveText).
type BusinessRule interface {
	IsBroken() bool
	Code() string    // stable code for clients, e.g. "inbound_message.text_required"
	Message() string // human-readable explanation
}

// BusinessRuleViolation is the single error type the domain returns for a broken rule.
type BusinessRuleViolation struct {
	Rule BusinessRule
}

func (v *BusinessRuleViolation) Error() string {
	return v.Rule.Code() + ": " + v.Rule.Message()
}

// CheckRule returns a *BusinessRuleViolation when the rule is broken.
func CheckRule(rule BusinessRule) error {
	if rule.IsBroken() {
		return &BusinessRuleViolation{Rule: rule}
	}
	return nil
}

// CheckRules checks rules in order and returns the first violation.
func CheckRules(rules ...BusinessRule) error {
	for _, rule := range rules {
		if err := CheckRule(rule); err != nil {
			return err
		}
	}
	return nil
}
