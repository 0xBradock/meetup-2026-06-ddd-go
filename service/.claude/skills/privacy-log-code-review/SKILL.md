---
name: privacy-log-code-review
description: analyse go application logs and source code to identify potential personal data protection issues, then add actionable todo comments directly in the code. use this skill when asked to inspect go logs, go repositories, logging statements, data flows, persistence, exports, integrations, or application behavior for gdpr/privacy compliance risks and produce in-code remediation tasks.
---

# Privacy Log and Code Review

## Goal

Inspect a Go repository and its logs to identify code changes needed to better respect personal data protection rules.

The expected result is not a legal opinion. The expected result is a practical engineering review that adds clear `TODO` comments in the relevant source files, near the code that should be reviewed or changed.

## Scope

Use this skill for Go projects.

Review:

- Go source files: `*.go`
- configuration files that influence logging, storage, exports, integrations, access control, retention, or observability
- log files or captured log outputs provided by the user
- tests only when they reveal logging or data-handling behavior that may also exist in production code

Ignore, unless directly relevant:

- generated files
- vendored dependencies
- binary files
- build artifacts
- unrelated documentation
- minified or compiled assets

## Core privacy rules to check

Use the following checklist when reviewing logs and code.

### Personal data in logs

Look for logs that may expose personal data, including:

- names, first names, usernames
- email addresses
- phone numbers
- postal addresses
- IP addresses
- user IDs, beneficiary IDs, customer IDs, employee IDs, account IDs
- payment, banking, fiscal, or financial data
- authentication tokens, API keys, session IDs, cookies, JWTs
- precise location data
- uploaded document metadata
- free-text fields that may contain personal data
- full request or response payloads
- database records dumped into logs
- errors that include raw payloads or sensitive context

Add TODOs when logs should be reduced, redacted, hashed, pseudonymized, structured differently, or removed.

### Sensitive or special-category data

Be especially strict with data that may reveal:

- health information
- sexual orientation or sex life
- racial or ethnic origin
- political, religious, or philosophical opinions
- trade union membership
- genetic data
- biometric identification data

Add TODOs when code may collect, process, persist, export, or log these data without a clearly visible reason.

### Minimization

Check whether the code collects, stores, returns, exports, or logs more data than necessary.

Add TODOs when:

- entire structs are logged instead of selected safe fields
- complete API responses are logged
- complete request bodies are logged
- database objects are serialized into logs
- unnecessary fields are persisted
- exports include columns that are not obviously needed
- debug logs are too verbose for production

### Purpose limitation

Check whether data collected for one purpose is reused for another purpose.

Add TODOs when:

- user, beneficiary, customer, or employee data is reused in marketing, analytics, support, or reporting flows without an obvious boundary
- data from one domain is copied into another domain without an explicit reason
- code sends data to external systems without clear purpose in the surrounding code
- background jobs process personal data in a way that is not obvious from their name, input, or comments

### Retention and deletion

Check whether the code creates data that may need a retention rule.

Add TODOs when:

- logs are retained without visible retention configuration
- audit or event tables grow indefinitely
- exported files are created without deletion or expiration
- temporary files may contain personal data and are not removed
- cache entries containing personal data have no TTL
- backups, archives, or object storage writes do not show an expiration strategy

### Access control and confidentiality

Check whether access to personal data is restricted.

Add TODOs when:

- handlers return personal data without authorization checks nearby
- admin/back-office endpoints expose broad user data
- repository methods return complete records when only a few fields are needed
- internal tools or debug endpoints expose personal data
- logs or metrics are sent to places where access control is unclear

### Security

Check for weak handling of personal data and secrets.

Add TODOs when:

- tokens, passwords, cookies, API keys, or secrets can be logged
- errors expose implementation details or raw payloads
- data is transmitted to external systems without visible security controls
- files containing personal data are written with unsafe permissions
- encryption, hashing, or pseudonymization appears missing where risk is high

### Third parties and external transfers

Check integrations with external tools or services.

Add TODOs when:

- personal data is sent to SaaS tools, analytics tools, support tools, CRMs, observability tools, or webhooks
- the destination country, provider, or purpose is unclear
- payloads sent to external services include more data than needed
- code sends personal data to logs, metrics, traces, error trackers, or monitoring platforms

### Data subject rights

Check whether the code can support user rights such as access, rectification, deletion, limitation, and portability.

Add TODOs when:

- personal data is duplicated across stores without clear synchronization
- deletion only removes part of the data
- exports are incomplete or inconsistent
- immutable logs or audit events may contain personal data without mitigation
- identifiers make it hard to find all data linked to a person

### Data breach detection

Check logs and code for patterns that may indicate a personal data breach risk.

Examples:

- logs showing data sent to a wrong recipient
- unauthorized access errors around personal data
- exposed debug endpoints
- unexpected bulk exports
- panics dumping request bodies or structs
- access denied messages followed by successful access
- unusually broad queries or scans of personal data

Add TODOs where code should improve detection, alerting, redaction, or incident traceability.

## Workflow

### 1. Discover the repository

Start by identifying:

- project layout
- main packages
- logging libraries
- HTTP framework or transport layer
- persistence layer
- background workers
- external integrations
- configuration files
- available log files

Use fast search tools when available, such as `rg`, `grep`, `find`, or language-aware code search.

### 2. Inspect logs first

Read the provided Go application logs.

Look for:

- personal data appearing directly in logs
- payload dumps
- stack traces with raw values
- request or response bodies
- user, beneficiary, customer, employee, or payment fields
- recurring events that point to risky code paths
- file names, function names, endpoints, job names, trace IDs, or error messages that can help locate the code

Extract useful clues:

- package names
- function names
- log messages
- endpoint paths
- job names
- queue names
- event names
- struct names
- field names

Then use those clues to search the code.

### 3. Inspect code

Search for privacy-relevant patterns, including:

```text
log.
logger.
zap.
zerolog.
slog.
logrus.
fmt.Printf
fmt.Println
Errorf
Debug
Debugf
Info
Infof
Warn
Warnf
Error
Errorf
With
WithField
WithFields
WithContext
request body
response body
payload
token
password
secret
cookie
authorization
email
phone
firstname
lastname
name
address
birth
iban
bic
payment
card
ip
user_id
beneficiary
customer
employee
export
csv
xlsx
pdf
webhook
analytics
segment
sentry
datadog
newrelic
openobserve
prometheus
trace
span
```

Also inspect:

- HTTP handlers and middleware
- background jobs and queue consumers
- repository/database code
- serializers and DTOs
- export generation
- file upload/download logic
- audit/event logging
- error handling
- observability and tracing setup
- third-party client code

### 4. Add TODO comments directly in code

When a risk is found, add a `TODO(privacy):` comment as close as possible to the relevant line.

Do not rewrite business logic unless explicitly asked.

Do not make speculative large refactors.

Prefer precise TODOs that an engineer can act on.

Use this format:

```go
// TODO(privacy): <action to take>. Reason: <privacy risk>. Evidence: <log/code clue>.
```

When the TODO is related to logs, prefer:

```go
// TODO(privacy): Redact or remove this logged field before production logging. Reason: it may contain personal data. Evidence: observed in application logs as "<short log clue>".
```

When the TODO is related to retention, prefer:

```go
// TODO(privacy): Define and enforce a retention or deletion rule for this stored/exported data. Reason: personal data must not be kept longer than necessary.
```

When the TODO is related to external services, prefer:

```go
// TODO(privacy): Review the data sent to this external provider and minimize the payload. Reason: personal data sharing with third parties must be necessary, documented, and secured.
```

When the TODO is related to access control, prefer:

```go
// TODO(privacy): Verify that this path enforces authorization before returning personal data. Reason: access to personal data must be limited to authorized users.
```

When the TODO is related to deletion or user rights, prefer:

```go
// TODO(privacy): Ensure this data is included in deletion/access/export workflows. Reason: duplicated personal data can make user-rights handling incomplete.
```

### 5. Avoid noisy TODOs

Do not add a TODO for every occurrence of the same issue.

Group repeated issues when possible:

- add one TODO in the shared helper
- add one TODO in the logger wrapper
- add one TODO in the serializer
- add one TODO near the integration client
- add one TODO in the configuration responsible for retention

Only add multiple TODOs when the code locations require different fixes.

### 6. Use severity internally when deciding priority

Use these severity levels to decide which TODOs matter most.

High severity:

- secrets, tokens, passwords, sessions, cookies, JWTs in logs
- payment, banking, health, biometric, or other sensitive data in logs
- full request/response payload logging in production paths
- unauthorized access to personal data
- exports or files with personal data and no access control
- external transfer of broad personal data without clear need

Medium severity:

- email, phone, address, IP, user IDs, customer IDs, beneficiary IDs in logs
- excessive data returned by handlers
- missing retention or deletion strategy
- duplicated personal data across systems
- unclear external integration payload

Low severity:

- unclear comments or naming around personal data
- missing documentation near a data flow
- logs that use indirect identifiers with limited risk
- tests that encourage unsafe logging patterns

Do not include the severity in every TODO unless useful. Use it when it helps prioritize.

Example:

```go
// TODO(privacy, high): Remove token from this log entry. Reason: authentication secrets must never be logged.
```

## Final response to the user

After modifying the code, provide a short report.

Include:

1. number of TODOs added
2. files changed
3. highest-risk findings
4. any areas that could not be checked
5. recommended next step

Use this format:

```markdown
## Privacy review summary

Added <n> TODO comments across <n> files.

### Files changed

- `<path>`: <short reason>
- `<path>`: <short reason>

### Highest-risk findings

- <finding>
- <finding>

### Not checked / limitations

- <limitation>

### Recommended next step

<one practical next action>
```

If no issues are found, do not modify files. Explain what was checked and why no TODO was added.

## Important constraints

- Do not provide legal advice.
- Do not claim the repository is GDPR-compliant.
- Do not remove or rewrite code unless explicitly asked.
- Do not expose personal data from logs in the final report.
- Quote only short, redacted log clues when needed.
- Prefer adding actionable TODOs over producing a long theoretical report.
- If logs contain obvious personal data, redact it in all explanations.
- If a potential data breach is suspected, tell the user to contact the relevant privacy/security contact immediately.
