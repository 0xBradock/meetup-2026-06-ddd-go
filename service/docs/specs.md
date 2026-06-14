# Technical Specification — Service Template

## 1. Context and objectives

### 1.1 Context
The service is the foundation of the WiiSmile ecosystem. It defines
who the main entities are (people, companies and contracts), in which context they operate, and
how they relate to each other. All products use this service as the single source
of truth for these fundamental structures.

This project is a rewrite of multiple incoherent legacy systems. The migration
will follow a Strangler Fig approach to gradually replace existing functionality.

### 1.2 Objectives
- Business: enforce a shared domain language and consistent structural rules.
- Technical: deliver a backend that centralizes service domain rules for people,
  companies, and contracts.
- Success indicators: a single, authoritative system with consistent behavior.

### 1.3 Atlassian scope reference
The following epics define the authoritative scope for Atlassian agent/MCP calls:
- ST-101: SERVICE IDENTITAIRE - Fondations d'architecture du service identitaire
- ST-103: SERVICE IDENTITAIRE - Migration incrementale du domaine client depuis le systeme legacy vers le nouveau service DDD
- ST-117: SERVICE IDENTITAIRE - Cartographie des APIs

## 2. Scope

Included:
- Company, group, and establishment management.
- People management in the context of a company.
- Multiple and independent contract lifecycle management in the context of a company.
- Roles and rights.
- service-level contract lifecycle management.
- Relationships between companies, people, and contracts.

Excluded:
- Billing computation.
- Product-specific business logic beyond this service domain.
- Financial transaction processing.
- Sending any kind of communication to a company or person.

## 3. Functional vision (macro level)

### 3.1 Aggregates

Person:
- `PersonMembership`: Person with a role in a company **dans une entreprise** (clé `(PersonId, CompanyId)`)

### Contrat

Each contract type is an aggregate (composition of `ContractCore`):
- `ContratWSService`
- `ContratTR`
- `ContratMeltem`
- `ContratMesAchatsPro`
- `ContratEcommerce`
- `ContratChequeCadeauxDistributor`
- `ContratPartnershipConvention`

Company:
- `Company`: WiiSmile client that contains one or more Persons and one or more Contracts.
- `Group`: Group of companies

### 3.2 Key use cases

- Create and manage companies, groups, and establishments.
- Create and manage person memberships per company.
- Activate/deactivate roles and manage rights per membership.
- Create and manage typed contracts and their lifecycle.
- Attach/detach contracts to role-specific data.

## 4. DDD vision

### 4.1 Domain and sub-domains

- Primary domain: service domain.
- Sub-domains: Person, Company, Contract.

### 4.2 Bounded context

The service is a single bounded context with one deployment unit. It exposes
domain operations via HTTP and message subscriptions.

### 4.3 Ubiquitous language (English)

See the glossary in Section 14. All code and docs must use these terms.

## 5. Target architecture

### 5.1 Overview

- Single deployable unit in a domain-scoped service architecture.
- Clean Architecture + DDD layering.
- Inbound drivers: HTTP and message subscriptions.

### 5.2 Structure

The repository uses a flat package layout rather than a deep DDD directory tree:

- `health/`: domain aggregates, value objects, use-case service, and repository/cache port interfaces.
- `db/`: persistence adapters (`pg/`, `mssql/`, `inmem/`, `valkey/`) and migration assets.
- `config/`: runtime configuration loading and validation.
- `httpserver/`: HTTP handlers, routes, and DTO mapping — transport only.
- `cmd/`: binary entrypoints (`server/`, `task_alive/`) that wire all layers together.

Composition root note: `cmd/main.go` wires process-level dependencies and injects them into `run(...)` (including environment access via `os.Getenv` as `getenv func(string) string`). This keeps runtime configuration loading testable and avoids hard-coding process-global calls inside `run`.

### 5.3 Key technical choices

- Language: Go.
- Database: PostgreSQL (primary).
- Cache: Memcache (read-through or invalidation strategy defined in Section 8).
- Messaging: abstracted interfaces; initial implementation via SNS/SQS.
- Auth: JWT at service boundary.

## 6. Technical modeling

### 6.1 Domain model (from `v3.mmd`)

Aggregates
- Company
- Group
- Establishment
- PartnerUnderConvention
- PersonMembership
- Contract (per type): ServiceWS, TR, Meltem, MesAchatsPro,
  Ecommerce, ChequeCadeauxDistributor, PartnershipConvention

Core components/value objects (examples)
- ContractCore, PersonCore, PersonCompanyKey
- CompanyID, GroupID, EstablishmentID, PersonID, ContractID, PartnerID
- Address, BankInfo, SIREN, SIRET, Role, Right

Key invariants
- A Company may belong to zero or one Group.
- A Group contains at least one Company.
- Establishments belong to exactly one Company and are deleted if the company is deleted.
- A PersonMembership is unique per (PersonID, CompanyID).
- Role-specific data can exist only while the role is active.
- Contract references on a role must belong to the same Company.
- Contract status transitions are constrained by domain methods.

### 6.2 Diagrams

- Class diagram: `v3.mmd`.
- Sequence diagrams: to be added for critical flows.

## 7. API and contracts

### 7.1 Role/action matrix

The following matrix is the authoritative access reference.

| Method / Action (aggregate)                 | Beneficiary | Director | Manager | Accountant | Referrer | Ambassador | Candidate |
|---------------------------------------------|------------:|---------:|--------:|-----------:|---------:|-----------:|----------:|
| **Company**                                 |             |        W |       W |            |          |            |           |
| Company.read                                |           R |        R |       R |          R |        R |          R |         R |
| Company.attach/detachEstablishment          |             |        W |       W |            |          |            |           |
| **Contract**                                |             |        W |       W |            |          |            |           |
| Contract.read                               |           R |        R |       R |          R |        R |          R |         R |
| TRContract.validateConsumptions             |             |          |         |          W |          |            |           |
| **PersonMembership**                        |             |        W |       W |            |          |            |           |
| PersonMembership.grant/revoke/replaceRights |             |        W |       W |            |          |            |           |
| PersonMembership.updatePerson               |           W |        W |       W |          W |        W |          W |         W |
| BeneficiaryData.attach/detach contract      |           W |        W |       W |            |          |            |           |
| DirectorData.attach/detach contract         |             |        W |       W |            |          |            |           |
| AccountantData.attach/detach contract       |             |          |         |          W |          |            |           |
| ReferrerData.attach/detach contract         |             |          |         |            |        W |            |           |
| DependantData.attachBeneficiary             |             |        W |       W |            |          |            |           |

### 7.2 Transport and API shape

- REST over HTTP for external clients.
- Message subscriptions for domain events and cross-service workflows.

REST endpoints (initial)
- `POST /companies`
- `GET /companies/{companyId}`
- `PATCH /companies/{companyId}`
- `DELETE /companies/{companyId}`
- `POST /companies/{companyId}/establishments`
- `PATCH /establishments/{establishmentId}`
- `DELETE /establishments/{establishmentId}`
- `POST /groups`
- `PATCH /groups/{groupId}`
- `POST /groups/{groupId}/companies/{companyId}` (attach)
- `DELETE /groups/{groupId}/companies/{companyId}` (detach)
- `POST /partners`
- `PATCH /partners/{partnerId}`

- `POST /companies/{companyId}/contracts/ws-service`
- `POST /companies/{companyId}/contracts/tr`
- `POST /companies/{companyId}/contracts/meltem`
- `POST /companies/{companyId}/contracts/mes-achats-pro`
- `POST /companies/{companyId}/contracts/ecommerce`
- `POST /companies/{companyId}/contracts/cheque-cadeaux-distributor`
- `POST /partners/{partnerId}/contracts/partnership-convention`
- `PATCH /contracts/{contractId}` (type-specific update)
- `POST /contracts/{contractId}:subscribe`
- `POST /contracts/{contractId}:terminate`
- `POST /contracts/{contractId}:renew`
- `POST /contracts/{contractId}:validate-consumptions` (TR only)

- `POST /companies/{companyId}/memberships`
- `GET /companies/{companyId}/memberships/{personId}`
- `PATCH /companies/{companyId}/memberships/{personId}/person`
- `POST /companies/{companyId}/memberships/{personId}:activate-role`
- `POST /companies/{companyId}/memberships/{personId}:deactivate-role`
- `POST /companies/{companyId}/memberships/{personId}:grant-right`
- `POST /companies/{companyId}/memberships/{personId}:revoke-right`
- `POST /companies/{companyId}/memberships/{personId}/beneficiary:attach-contract`
- `POST /companies/{companyId}/memberships/{personId}/director:attach-contract`
- `POST /companies/{companyId}/memberships/{personId}/accountant:attach-contract`
- `POST /companies/{companyId}/memberships/{personId}/referrer:attach-contract`

### 7.3 External integrations

- Messaging: SNS/SQS initially, abstracted behind ports.
- Protocols: HTTP and message subscriptions.
- Constraints: idempotent handlers, at-least-once delivery.

Message subscriptions (initial placeholders)
- `company.created`
- `company.updated`
- `company.deleted`
- `establishment.created`
- `establishment.updated`
- `establishment.deleted`
- `membership.created`
- `membership.updated`
- `membership.role.activated`
- `membership.role.deactivated`
- `contract.created`
- `contract.updated`
- `contract.subscribed`
- `contract.terminated`
- `contract.renewed`

## 8. Data and persistence

### 8.1 Data model

- Relational model aligned to aggregates and value objects.
- Contract tables per contract type or a single contracts table with type
  discrimination (decision TBD).

### 8.2 Persistence and caching

- Primary datastore: PostgreSQL.
- Cache: Memcache for read optimization (policy TBD).
- Transactions: ACID in Postgres; cache never source of truth.
- High availability: required (TBD targets).

## 9. Security

- JWT authentication at the boundary.
- Authorization enforced by role/rights logic in the application layer.
- Sensitive data handled according to least-privilege access.

## 10. Non-functional requirements
- Availability: TBD (service is critical; HA required).
- Scalability: horizontal scaling via stateless services.
- Observability: logs, metrics, and traces for key operations.
- Performance: TBD latency and throughput targets.

## 11. Testing strategy
- Use table-driven test whenever possible (first check for existing tables that can accomodate the new tests or create if no similar tests exist)
- Unit tests for domain use cases and we should always have 100 coverage and mock driven packages
- Integration tests for persistence and adapters.
- Contract tests for HTTP boundaries.
- End-to-end tests for critical workflows.

## 12. Deployment and operations

### 12.1 Deployment model (ECS)
- Deployed as a single ECS service with multiple tasks.
- Health checks and readiness probes required.
- Auto-scaling based on CPU and latency.
- Environment configuration via ECS task definitions.

### 12.2 Edge and routing
- API Gateway in front of the service with routing rules.
- HTTP endpoints exposed via gateway configuration.

### 12.3 Operational concerns
- Centralized logging, metrics, and tracing.
- Backups and database monitoring.
- Alerting on SLA/SLO breaches (TBD targets).

## 13. Migration strategy (Strangler Fig)

### 13.1 Approach
- Incrementally route traffic from legacy systems to this service.
- Use API Gateway rules to direct specific routes to the new service.

### 13.2 Data migration
- Preferred: AWS DMS for replication/migration.
- Fallback: dual-write with regular reconciliation checks.
- Define a cutover plan per domain area (Company, Person, Contract).

### 13.3 Coexistence and decommission
- Run legacy and new systems in parallel during migration.
- Decommission legacy components once parity and data integrity are proven.

### 13.4 Risks and mitigations
- Data divergence: mitigated by DMS or dual-write reconciliation.
- Contract incompatibility: mitigated by explicit API contracts and testing.
- Latency shifts: monitored via gateway and service metrics.

## 14. Glossary (English)

### Ubiquitous Language

|    Fr | En    |
| Personne   | Person    |
|Entreprise     | Company     |
|Groupe     | Group     |
|Etablissement     | Establishment     |
|PartenaireSousConvention | PartnerUnderConvention     |
|Entreprise     | Company     |
|Entreprise     | Company     |
|Role     | Role     |
|Droit     | Right     |

Note: every time we introduce a new transaction, it must be added to the ubiquitous language (glossary) as well.

Person: an individual in the system, scoped to a company via PersonMembership.
Company: a WiiSmile client entity that contains one or more persons and one or more contracts.
Group: a collection of one or more companies under a shared grouping.
Establishment: a legal establishment that belongs to exactly one company.
PartnerUnderConvention: a partner entity operating under a partnership convention.
PersonMembership: the association between a person and a company, unique per (PersonID, CompanyID), carrying roles and rights.
Contract: a typed agreement aggregate (e.g., ServiceWS, TR, Meltem) with its own lifecycle.
Role: the position a person holds within a PersonMembership (e.g., Beneficiary, Director, Manager). Roles can be activated or deactivated and gate role-specific data.
Right: a specific permission/action that can be granted or revoked within a PersonMembership. Rights are independent permissions associated with roles.

## 15. Annexes 

- [Current list of APIs](https://wiismilefr.sharepoint.com/:x:/r/sites/service.it.po2/Documents%20partages/Dette%20Technique/listing-api.xlsx?d=w4f83a95e28e845bfa80549736b46635e&csf=1&web=1&e=H1P6iz)
