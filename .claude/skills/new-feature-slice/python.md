# Python slice template (Agent: propose a reply when a conversation is updated)

```
src/agent_orchestrator/features/propose_reply/
  __init__.py
  message.py      # ConversationUpdated.v1 payload (consumed event)
  command.py      # ProposeReplyCommand
  result.py       # ProposeReplyResult
  validator.py    # rules beyond types
  handler.py      # use case orchestration
  graph.py        # LangGraph workflow for this use case
  consumer.py     # Kafka entry point: inbox -> map -> handler
tests/features/propose_reply/test_handler.py
```

```python
# command.py
class ProposeReplyCommand(BaseModel):
    model_config = ConfigDict(frozen=True)
    message_id: UUID
    tenant_id: UUID
    conversation_key: str
    recent_messages: list[str]

# result.py
class ProposeReplyResult(BaseModel):
    proposal_id: UUID
    confidence: float
```

```python
# validator.py
class ProposeReplyValidator:
    def validate(self, command: ProposeReplyCommand) -> None:
        if not command.recent_messages:
            raise InvalidCommand("recent_messages must not be empty")
```

```python
# handler.py
class ProposeReplyHandler:
    def __init__(
        self,
        validator: ProposeReplyValidator,
        graph: ReplyGraph,
        proposals: ReplyProposalRepository,   # port declared in domain
        uow: UnitOfWork,
    ) -> None:
        self._validator = validator
        self._graph = graph
        self._proposals = proposals
        self._uow = uow

    async def handle(self, command: ProposeReplyCommand) -> ProposeReplyResult:
        self._validator.validate(command)
        draft = await self._graph.run(command.recent_messages)
        proposal = ReplyProposal.create(
            tenant_id=TenantId(command.tenant_id),
            conversation_key=ConversationKey(command.conversation_key),
            text=draft.text,
            confidence=Confidence(draft.confidence),   # value object validates 0..1
        )
        async with self._uow:                          # proposal + outbox, one transaction
            await self._proposals.add(proposal)
            self._uow.outbox.add(proposal.to_reply_proposed_event())
        return ProposeReplyResult(proposal_id=proposal.id, confidence=proposal.confidence.value)
```

```python
# consumer.py
class ProposeReplyConsumer:
    def __init__(self, handler: ProposeReplyHandler, inbox: Inbox) -> None:
        self._handler = handler     # called directly: no mediator
        self._inbox = inbox

    async def consume(self, envelope: Envelope) -> None:
        message = ConversationUpdatedMessage.model_validate(envelope.payload)
        if not await self._inbox.try_register(envelope.message_id, consumer="propose_reply"):
            return
        await self._handler.handle(ProposeReplyCommand(
            message_id=envelope.message_id,
            tenant_id=envelope.tenant_id,
            conversation_key=envelope.conversation_key,
            recent_messages=message.recent_messages,
        ))
```

```python
# bootstrap/container.py (composition root, the only place that knows concrete classes)
from dependency_injector import containers, providers

class Container(containers.DeclarativeContainer):
    config = providers.Configuration()

    db_pool = providers.Resource(create_pool, dsn=config.database_url)
    chat_model = providers.Singleton(build_chat_model, provider=config.llm_provider, model=config.llm_model)

    unit_of_work = providers.Factory(PostgresUnitOfWork, pool=db_pool)
    inbox = providers.Factory(PostgresInbox, pool=db_pool)
    reply_proposals = providers.Factory(PostgresReplyProposalRepository, pool=db_pool)

    propose_reply_handler = providers.Factory(
        ProposeReplyHandler,
        validator=providers.Factory(ProposeReplyValidator),
        graph=providers.Singleton(ReplyGraph, chat_model=chat_model),
        proposals=reply_proposals,
        uow=unit_of_work,
    )
    propose_reply_consumer = providers.Factory(
        ProposeReplyConsumer, handler=propose_reply_handler, inbox=inbox)
```

Simplification: in this sketch the repository and the unit of work would
need to share one database connection per message; the real implementation
makes the unit of work own the connection and hand it to the repository.
