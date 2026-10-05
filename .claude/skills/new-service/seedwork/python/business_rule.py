from abc import ABC, abstractmethod


class BusinessRule(ABC):
    """An invariant named in the ubiquitous language (e.g. ReplyMustHaveText)."""

    code: str
    """Stable code for clients, e.g. "reply_proposal.text_required"."""

    message: str

    @abstractmethod
    def is_broken(self) -> bool: ...


class BusinessRuleViolation(Exception):
    """The single exception type the domain raises for a broken rule."""

    def __init__(self, rule: BusinessRule) -> None:
        super().__init__(f"{rule.code}: {rule.message}")
        self.rule = rule


def check_rule(rule: BusinessRule) -> None:
    if rule.is_broken():
        raise BusinessRuleViolation(rule)


def check_rules(*rules: BusinessRule) -> None:
    for rule in rules:
        check_rule(rule)
