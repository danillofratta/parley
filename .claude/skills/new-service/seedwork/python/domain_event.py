from dataclasses import dataclass, field
from datetime import datetime
from uuid import UUID

from .ids import new_id


@dataclass(frozen=True, kw_only=True)
class DomainEvent:
    """A business fact that already happened, named in past tense."""

    aggregate_id: UUID
    tenant_id: UUID
    occurred_at: datetime
    event_id: UUID = field(default_factory=new_id)
