package rules

import "fmt"

type Rule interface {
	IsBroken() bool
	Code() string
	Message() string
}

type BrokenRuleError struct {
	code    string
	message string
}

func NewBrokenRuleError(code, message string) BrokenRuleError {
	return BrokenRuleError{code: code, message: message}
}

func (e BrokenRuleError) Error() string {
	if e.code == "" {
		return e.message
	}
	return fmt.Sprintf("%s: %s", e.code, e.message)
}

func (e BrokenRuleError) Code() string {
	return e.code
}

func CheckRules(rules ...Rule) error {
	for _, rule := range rules {
		if rule == nil {
			continue
		}
		if rule.IsBroken() {
			return NewBrokenRuleError(rule.Code(), rule.Message())
		}
	}
	return nil
}