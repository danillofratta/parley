-- One schema and one user per service (ADR 0006).
-- Each service can only touch its own schema.
REVOKE ALL ON SCHEMA public FROM PUBLIC;

CREATE USER gateway WITH PASSWORD 'gateway';
CREATE SCHEMA gateway AUTHORIZATION gateway;

CREATE USER dispatcher WITH PASSWORD 'dispatcher';
CREATE SCHEMA dispatcher AUTHORIZATION dispatcher;

CREATE USER conversations WITH PASSWORD 'conversations';
CREATE SCHEMA conversations AUTHORIZATION conversations;

CREATE USER agent WITH PASSWORD 'agent';
CREATE SCHEMA agent AUTHORIZATION agent;