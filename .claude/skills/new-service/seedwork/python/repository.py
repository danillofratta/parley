from typing import Protocol, TypeVar
from uuid import UUID

from .aggregate_root import AggregateRoot

T = TypeVar("T", bound=AggregateRoot)


class Repository(Protocol[T]):
    """Persistence port for one aggregate root. No query methods."""

    async def get_by_id(self, tenant_id: UUID, aggregate_id: UUID) -> T | None: ...

    async def add(self, aggregate: T) -> None:
        """Raises AggregateAlreadyExists when a unique business key already exists."""

    async def update(self, aggregate: T) -> None:
        """Raises ConcurrencyConflict when the stored version differs."""


class AggregateAlreadyExists(Exception):
    pass


class ConcurrencyConflict(Exception):
    pass
