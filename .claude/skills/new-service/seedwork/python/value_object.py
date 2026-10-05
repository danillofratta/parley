from dataclasses import dataclass


@dataclass(frozen=True)
class ValueObject:
    """Base of value objects: immutable and compared by value.

    Subclasses are also declared with @dataclass(frozen=True) and validate
    their rules in __post_init__ with check_rule.
    """
