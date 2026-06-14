---
name: ddd-implementation
description: Comprehensive DDD modeling reference covering Entities, Value Objects, Aggregates, Domain/Application Services, Domain Events, and Hexagonal Architecture patterns.
---

# Skill: Domain-Driven Design Feature Modeling and Implementation Guidance

## Purpose

Use this skill when implementing a feature that contains meaningful business rules, domain language, workflows, invariants, policies, or integrations that should be modeled using Domain-Driven Design (DDD) and, when appropriate, Hexagonal Architecture.

This skill helps the agent decide:

- how to model the domain
- when to use an Entity, Value Object, Aggregate, Domain Service, Application Service, or Infrastructure Service
- when to publish or consume Domain Events
- how to reason about bounded contexts and context mapping
- how to keep business logic separate from technical concerns
- how to choose designs that remain testable with unit, integration, and functional tests

## Core Objective

When implementing a feature, prefer a model that:

- captures business concepts explicitly
- protects invariants in the right place
- keeps domain logic inside the domain model
- isolates infrastructure and framework concerns
- avoids anemic domain models when the domain is rich
- avoids overengineering when the feature is simple CRUD with no meaningful business rules
- preserves ubiquitous language from the business domain

## Fundamental DDD Definitions

### Domain

The domain is the business problem space the software is solving.

The agent must identify the business concepts, rules, terms, constraints, and workflows before deciding on implementation structure.

### Ubiquitous Language

A shared language between business experts and developers.

The agent should use names from the business language consistently in code, tests, events, APIs, and documentation.

### Bounded Context

A bounded context is a clear boundary within which a domain model is valid and consistent.

The same word may have different meanings across different bounded contexts.

Examples:

- `Customer` in Billing may mean an invoiced legal entity
- `Customer` in Support may mean a person receiving help

The agent should avoid merging meanings from different contexts into one model.

### Entity

An Entity is an object defined primarily by identity and continuity over time, not just by its attributes.

Use an Entity when:

- identity matters
- the object changes over time
- lifecycle matters
- business rules refer to the same thing across state changes

Examples:

- Order
- Customer
- Subscription
- Invoice

### Value Object

A Value Object is defined by its attributes and has no conceptual identity.

Use a Value Object when:

- identity does not matter
- equality is based on value
- it represents a descriptive concept
- it should be immutable whenever possible

Examples:

- Money
- Address
- DateRange
- EmailAddress
- Percentage

### Aggregate

An Aggregate is a consistency boundary that groups one root Entity and related objects that must change together under enforced invariants.

Use an Aggregate when:

- multiple objects must remain consistent together
- invariants must be enforced transactionally
- external code should not modify internal parts directly
- one root should control access to the cluster

The Aggregate Root is the only entry point from outside the aggregate.

Examples:

- Order as root, with OrderLines inside
- Cart as root, with CartItems inside

### Repository

A Repository provides collection-like access to Aggregates or sometimes Entities.

Use a Repository for:

- loading an Aggregate needed for business operations
- saving an Aggregate after behavior has changed

Do not use repositories for:

- arbitrary query dumping
- business logic
- technical orchestration

### Domain Service

A Domain Service contains domain logic that does not naturally belong to a single Entity or Value Object.

Use a Domain Service when:

- the logic is domain-specific
- it spans multiple domain objects
- putting it on one entity would be unnatural
- it expresses business rules, policies, or calculations

Examples:

- pricing policy
- commission calculation
- allocation logic
- eligibility policy

### Application Service

An Application Service coordinates a use case.

Use an Application Service to:

- receive a command or request
- load aggregates through repositories
- invoke domain behavior
- call domain services if needed
- persist results
- publish events
- coordinate transactions
- return a response DTO or result

Application Services should not contain core business rules that belong in the domain.

### Infrastructure Service

An Infrastructure Service implements technical capabilities or external integrations.

Use an Infrastructure Service for:

- email sending
- payment gateways
- storage
- message buses
- clocks
- UUID generation
- external APIs
- persistence implementations

Infrastructure Services should not define business rules.

### Domain Event

A Domain Event represents something meaningful that happened in the domain.

Use a Domain Event when:

- the business cares that something happened
- other parts of the same context or other contexts may react
- the event is part of the business narrative
- the event helps decouple side effects from core decisions

Examples:

- OrderPlaced
- PaymentCaptured
- SubscriptionCancelled
- InvoiceIssued

A Domain Event should describe a past fact, not an instruction.

## Decision Guide: How to Model a Concept

When the implementation introduces a new business concept, follow this decision process.

### Step 1: Is it a domain concern at all?

Treat it as a domain concern if it contains:

- business meaning
- business rules
- invariants
- lifecycle rules
- policy decisions
- meaningful state transitions

If it is only technical plumbing, keep it out of the domain layer.

### Step 2: Should it be a Value Object?

Choose a Value Object if most of the following are true:

- it has no business identity
- two instances with the same data are interchangeable
- it describes something rather than being something
- immutability is natural
- it helps express the domain more clearly

Examples:

- `Money`
- `Address`
- `CustomerName`
- `TaxRate`
- `TimeWindow`

Smells that something should be a Value Object:

- many primitive fields travel together
- validation rules repeat everywhere
- the concept has business meaning but no lifecycle
- strings or numbers are being passed around without meaning

Prefer a Value Object instead of primitive obsession.

### Step 3: Should it be an Entity?

Choose an Entity if most of the following are true:

- identity matters
- the object persists across state changes
- business rules refer to this specific instance over time
- history or lifecycle matters
- equality by value would be wrong

Examples:

- `Order`
- `Product`
- `Account`
- `Reservation`

Smells that something should be an Entity:

- it changes state over time
- users refer to that exact one
- audit/history matters
- it has status transitions

### Step 4: Should it be an Aggregate?

Choose an Aggregate when you need a consistency boundary.

Use an Aggregate if:

- multiple related objects must remain consistent together
- there are invariants that must always hold after a transaction
- the root should control all modifications
- concurrent modifications must be protected at the boundary

Examples:

- `Order` with `OrderLines`
- `Invoice` with `InvoiceItems`
- `Cart` with `CartItems`

Do not make something an Aggregate just because objects are related.

Ask:

- must these things be changed in one transaction?
- is there an invariant spanning them?
- should external code access only the root?

Keep Aggregates as small as possible while still protecting invariants.

### Step 5: Should it be a Domain Service?

Choose a Domain Service when:

- the behavior is business logic
- the behavior does not fit naturally on one Entity or Value Object
- the behavior coordinates rules across domain objects
- the behavior is not merely orchestration

Examples:

- `PricingPolicy`
- `FraudAssessmentPolicy`
- `ShipmentRoutingService`
- `DiscountCalculationService`

Do not use a Domain Service to avoid putting valid behavior on Entities or Value Objects.

If logic clearly belongs to an Aggregate, put it there first.

### Step 6: Should it be an Application Service?

Choose an Application Service when the need is use-case orchestration:

- start use case
- validate coarse input
- load aggregate(s)
- invoke domain behavior
- save
- publish events
- call external ports
- return result

A good rule:

- business decision => domain
- use case coordination => application
- technical integration => infrastructure

### Step 7: Should it be an Infrastructure Service?

Choose an Infrastructure Service when the capability is technical rather than business-specific.

Examples:

- SMTP email sender
- Stripe payment adapter
- database repository implementation
- Kafka event publisher
- file storage adapter

Use ports and interfaces so the application and domain remain independent of implementation details.

## Decision Guide: Where Should Logic Live?

When the implementation requires a function or operation, use the following rules.

### Put logic inside an Entity or Aggregate when:

- it changes the entity’s own state
- it protects invariants
- it validates valid state transitions
- it belongs to the entity’s lifecycle
- it uses the entity’s own data to make decisions

Examples:

- `Order.confirm()`
- `Subscription.cancel()`
- `Invoice.markAsPaid()`

### Put logic inside a Value Object when:

- it operates on the value’s own data
- it enforces its own invariants
- it provides domain-specific behavior
- it produces other values

Examples:

- `Money.add()`
- `DateRange.overlaps()`
- `EmailAddress.normalized()`

### Put logic inside a Domain Service when:

- the logic is domain logic
- it spans multiple concepts
- no single object is the natural owner
- the operation expresses policy or calculation

Examples:

- `PricingPolicy.calculateTotal(...)`
- `EligibilityPolicy.isEligible(...)`

### Put logic inside an Application Service when:

- it sequences actions for a use case
- it manages transactional boundaries
- it handles authorization at use-case level
- it coordinates repositories, ports, domain objects, and events

Examples:

- place order workflow
- approve application workflow
- renew subscription workflow

### Put logic inside an Infrastructure Service when:

- it communicates with external systems
- it persists or retrieves data
- it sends messages
- it depends on frameworks or transport concerns

## Aggregate Design Guidelines

### Use an Aggregate to protect invariants

The Aggregate Root must guarantee that business invariants always hold after any state change.

Examples:

- an order cannot confirm with zero lines
- a reservation cannot exceed capacity
- a payment cannot be captured twice

### Keep aggregate boundaries small

Do not include related objects unless they must be consistent in the same transaction.

Avoid large aggregates that:

- load too much data
- create locking issues
- reduce throughput
- hide multiple business concepts in one cluster

### Reference other aggregates by identity

Do not keep direct object graphs across aggregates unless there is a strong and local modeling reason.

Prefer:

- `customerId`
- `productId`
- `invoiceId`

This reduces coupling and preserves boundaries.

### Only expose behavior, not setters

Prefer intention-revealing methods:

- `order.addLine(productId, quantity)`
- `order.confirm()`

Avoid unrestricted setters that allow invalid states.

### Enforce invariants inside the root

If an invariant belongs to the aggregate, do not rely on callers to preserve it.

## Domain Events Guidance

### When to emit a Domain Event

Emit a Domain Event when a completed domain state change is meaningful.

Examples:

- after order is placed
- after invoice is issued
- after account is suspended

The event should represent something that already happened.

### What Domain Events are good for

Use Domain Events to:

- trigger side effects
- notify other parts of the same bounded context
- integrate with other bounded contexts
- preserve explicit domain knowledge
- support eventual consistency

### What not to do

Do not use Domain Events for:

- replacing simple direct method calls inside the same aggregate
- hiding poor modeling
- representing commands
- exposing infrastructure details as domain facts

### Event design guidelines

A Domain Event should:

- be named in past tense
- contain only relevant business data
- avoid technical noise
- be stable and understandable to consumers
- describe facts, not procedures

Good examples:

- `OrderPlaced`
- `CustomerRegistered`
- `PaymentFailed`

Poor examples:

- `DoSendEmail`
- `DatabaseRowInserted`

### Internal vs integration events

The agent should distinguish:

- Domain Events: meaningful inside the domain model
- Integration Events: messages published for external systems or other bounded contexts

An Integration Event may be derived from a Domain Event, but they are not necessarily identical.

## Context Mapping Guidance

When the feature interacts with concepts outside its immediate model, reason explicitly about bounded contexts and their relationship.

### First identify context boundaries

Ask:

- does the same term mean different things in different parts of the business?
- are there separate teams, workflows, models, or policies?
- would merging these models create confusion or leakage?

If yes, prefer separate bounded contexts.

### Common context map relationships

#### Partnership

Use when two contexts evolve closely and coordinate actively.

#### Shared Kernel

Use only when a very small, carefully governed part of the model is truly shared.

Avoid large shared kernels.

#### Customer-Supplier

Use when one context depends on another and can influence its evolution.

#### Conformist

Use when one context must adopt another’s model with little negotiation power.

#### Anti-Corruption Layer (ACL)

Use when one context must integrate with another model without polluting its own domain.

The ACL is especially important when:

- integrating with legacy systems
- consuming poorly designed APIs
- integrating external vendor models
- avoiding leakage of foreign concepts into the local model

#### Open Host Service

Use when a context exposes a clear published protocol for others.

#### Published Language

Use when contexts exchange messages or contracts through a well-defined shared language.

## How to Decide if a New Bounded Context Is Needed

Create or preserve a separate bounded context when:

- the same words have different meanings
- rules differ significantly
- workflows differ significantly
- the model would become inconsistent or contradictory if unified
- team ownership differs
- integration with an external or legacy model would distort the internal domain

Do not create a new bounded context for every module or package.

A bounded context is a semantic boundary, not just a folder structure.

## Heuristics for Rich Domain vs Simple CRUD

### Prefer rich domain modeling when:

- there are many business rules
- invariants matter
- workflows have multiple state transitions
- calculations and policies are central
- the business language is rich
- errors arise from invalid business behavior, not only invalid input

### Prefer a simpler service-oriented or CRUD style when:

- the feature mostly stores and retrieves data
- there are almost no business rules
- no meaningful invariants exist beyond basic validation
- state transitions are trivial
- the domain is not a source of business complexity

Do not force DDD tactical patterns where they add no value.

## Hexagonal Architecture Alignment

When implementing the solution, separate responsibilities into the following layers.

### Domain

Contains:

- Entities
- Value Objects
- Aggregates
- Domain Services
- Domain Events
- repository interfaces if needed by the domain or application boundary

Must not depend on frameworks, transports, databases, or external APIs.

### Application

Contains:

- use case handlers or services
- orchestration logic
- transaction boundaries
- ports for external capabilities
- DTOs or commands if needed

### Infrastructure

Contains:

- repository implementations
- database mappers
- message brokers
- external API clients
- web framework adapters
- queue consumers and producers
- email, SMS, or payment adapters

### Interface / Delivery

Contains:

- controllers
- CLI commands
- message listeners
- HTTP request and response mapping
- transport-specific validation

## Implementation Procedure for the Agent

When asked to implement a feature, follow this sequence.

### 1. Extract the business language

Identify:

- nouns
- verbs
- invariants
- policies
- lifecycle changes
- external systems
- events
- roles and actors

### 2. Identify the bounded context

Decide where this feature belongs and whether it crosses contexts.

### 3. Classify concepts

For each concept, determine whether it is:

- Entity
- Value Object
- Aggregate Root
- Domain Service
- Application Service
- Infrastructure Service
- Domain Event
- Repository
- external port or adapter

### 4. Define invariants explicitly

State what must always be true.

Examples:

- an order must have at least one line before confirmation
- a booking must not exceed available seats
- a coupon cannot be applied after expiration

### 5. Design behavior around invariants

Put behavior where invariants can be enforced reliably.

### 6. Identify side effects and integrations

Decide which are:

- direct within the use case
- event-driven
- external integrations behind ports

### 7. Produce implementation structure

Return a design that shows:

- bounded context
- aggregate boundaries
- entities and value objects
- services and responsibilities
- events
- repositories
- ports and adapters

### 8. Validate against architectural quality

Check that:

- domain logic is not leaking into controllers or infrastructure
- application services are orchestrating, not deciding core rules
- aggregates are not oversized
- value objects are used instead of loose primitives where appropriate
- events describe facts
- foreign models are isolated with ACLs where needed

## Decision Checklist

Before finalizing the implementation, verify the following.

### Modeling checklist

- Does each concept have a clear business meaning?
- Are Value Objects used for descriptive, identity-free concepts?
- Are Entities used where identity and lifecycle matter?
- Are Aggregates used only where consistency boundaries are needed?
- Are aggregate roots the only external entry points?

### Behavior checklist

- Is domain logic inside domain objects or domain services?
- Are application services coordinating rather than holding core rules?
- Are infrastructure concerns isolated behind ports and adapters?
- Are setters avoided in favor of intention-revealing methods?

### Event checklist

- Are important business facts represented as Domain Events?
- Are events named as past facts?
- Are integration concerns separated from internal domain concerns?

### Context checklist

- Is the bounded context clear?
- Are foreign models prevented from leaking into the local model?
- Is an Anti-Corruption Layer needed?

### Simplicity checklist

- Is the design no more complex than necessary?
- Has DDD been applied because the domain needs it, not as ceremony?

## Output Expectations for the Agent

When using this skill to propose or implement a feature, the agent should produce:

1. a short domain analysis
2. identified bounded context or contexts
3. a classification of main concepts:
   - Entity
   - Value Object
   - Aggregate
   - Domain Service
   - Application Service
   - Infrastructure Service
   - Domain Event
4. key invariants
5. suggested object and service responsibilities
6. context mapping notes if other systems or contexts are involved
7. implementation guidance consistent with Hexagonal Architecture
8. testing guidance

## Testing Guidance

The resulting design should support the following.

### Unit tests

Use for:

- Value Object rules
- Entity behavior
- Aggregate invariants
- Domain Service policies

These should run without infrastructure.

### Integration tests

Use for:

- repository implementations
- database mapping
- message bus integration
- external adapter behavior
- anti-corruption layer translation

### Functional tests

Use for:

- end-to-end use case behavior
- controller or transport interaction
- full application workflow across layers

If the design is hard to unit test without infrastructure, the domain and application boundary is probably too coupled.

## Anti-Patterns to Avoid

Avoid these unless there is a strong and explicit reason.

### Anemic domain model

Business logic lives only in services while entities are data containers.

### Large aggregates

Too many objects inside one transactional boundary.

### Primitive obsession

Using raw strings, integers, and arrays where a domain concept deserves a Value Object.

### Service dumping

Putting all logic into generic services because modeling decisions were skipped.

### Leaky infrastructure

Database, framework, or external API concepts leaking into the domain model.

### Shared model without boundaries

Using one model across multiple contexts with different meanings.

### Event misuse

Using events as commands, or as a substitute for proper aggregate modeling.

## Default Decision Biases

When uncertain, prefer these defaults:

- prefer Value Object over primitive
- prefer behavior on domain objects over procedural orchestration
- prefer small Aggregates
- prefer Application Services for orchestration only
- prefer Domain Services only when behavior has no natural home
- prefer explicit Domain Events for meaningful domain facts
- prefer Anti-Corruption Layers when integrating foreign models
- prefer one bounded context language at a time
- prefer simpler design when the domain is truly simple

## Example Classification Patterns

### Example: Ordering

- `Order` => Aggregate Root / Entity
- `OrderLine` => Entity inside Aggregate or internal child concept depending on identity needs
- `Money` => Value Object
- `OrderPlaced` => Domain Event
- `PlaceOrderHandler` => Application Service
- `PricingPolicy` => Domain Service
- `OrderRepository` => Repository
- payment gateway adapter => Infrastructure Service

### Example: Customer registration

- `Customer` => Entity or Aggregate Root if it protects registration invariants
- `EmailAddress` => Value Object
- `CustomerRegistered` => Domain Event
- `RegisterCustomerHandler` => Application Service
- email sender => Infrastructure Service
- identity provider adapter => Infrastructure Service

## Final Instruction to the Agent

Do not begin with frameworks, tables, endpoints, or database schemas.

Begin by understanding the domain language, invariants, and consistency boundaries.

Model the business first.  
Place behavior where it best protects correctness.  
Keep boundaries explicit.  
Use DDD tactically where it adds clarity and safety, and stay pragmatic where the domain is simple.
