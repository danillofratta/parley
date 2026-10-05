from dataclasses import dataclass
from uuid import UUID

from .business_rule import BusinessRule, check_rule
from .domain_event import DomainEvent
from .entity import Entity


@dataclass(frozen=True)
class TenantMustBeInformed(BusinessRule):
    tenant_id: UUID | None
    code = "aggregate.tenant_required"
    message = "every aggregate must belong to a tenant"

    def is_broken(self) -> bool:
        return self.tenant_id is None or self.tenant_id.int == 0


class AggregateRoot(Entity):
    """Entity that owns a consistency boundary, belongs to a tenant and raises domain events."""

    def __init__(self, tenant_id: UUID, entity_id: UUID | None = None) -> None:
        check_rule(TenantMustBeInformed(tenant_id))
        super().__init__(entity_id)
        self._tenant_id = tenant_id
        self._domain_events: list[DomainEvent] = []

    @property
    def tenant_id(self) -> UUID:
        return self._tenant_id

    @property
    def domain_events(self) -> tuple[DomainEvent, ...]:
        return tuple(self._domain_events)

    def _raise(self, event: DomainEvent) -> None:
        self._domain_events.append(event)

    def clear_domain_events(self) -> None:
        """Called by persistence after the outbox rows are committed."""
        self._domain_events.clear()
