from datetime import datetime
from uuid import UUID

from .ids import new_id


class Entity:
    """Base of every entity: identity, audit trail and version. Equality by id."""

    def __init__(self, entity_id: UUID | None = None) -> None:
        self._id = entity_id or new_id()
        self._created_at: datetime | None = None
        self._created_by: str | None = None
        self._updated_at: datetime | None = None
        self._updated_by: str | None = None
        self._version = 0

    @property
    def id(self) -> UUID:
        return self._id

    @property
    def created_at(self) -> datetime | None:
        return self._created_at

    @property
    def created_by(self) -> str | None:
        return self._created_by

    @property
    def updated_at(self) -> datetime | None:
        return self._updated_at

    @property
    def updated_by(self) -> str | None:
        return self._updated_by

    @property
    def version(self) -> int:
        return self._version

    def stamp_created(self, at: datetime, by: str) -> None:
        """Called by persistence (the Auditor) on the first save."""
        self._created_at = at
        self._created_by = by
        self._version = 1

    def stamp_updated(self, at: datetime, by: str) -> None:
        """Called by persistence (the Auditor) on every update."""
        self._updated_at = at
        self._updated_by = by
        self._version += 1

    def __eq__(self, other: object) -> bool:
        return isinstance(other, Entity) and type(self) is type(other) and self._id == other._id

    def __hash__(self) -> int:
        return hash(self._id)
