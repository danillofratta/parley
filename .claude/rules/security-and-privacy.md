# Security and privacy (LGPD)

- Never log message text, names, phone numbers or emails. Log ids only.
- Secrets come from environment variables (locally) or the cloud secret store.
  Never commit `.env`, keys or Terraform state.
- Compare secrets with constant-time functions.
- Limit request body sizes and set server timeouts.
- Fail closed: a missing security setting stops the service from starting.
- Every query filters by `tenantId`.
- Mask personal data before sending text to an LLM when the feature allows it.
