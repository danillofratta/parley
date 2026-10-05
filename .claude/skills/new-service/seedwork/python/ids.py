from uuid import UUID

try:
    from uuid import uuid7  # Python 3.14+
except ImportError:  # pragma: no cover - older interpreters
    from uuid6 import uuid7


def new_id() -> UUID:
    """Return a new time-ordered identity (UUID v7) for entities and events."""
    return uuid7()
