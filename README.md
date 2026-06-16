# Software Rewrite — Presentation References

A curated collection of articles and case studies supporting a presentation on software rewrite strategies. The central thesis: **incremental, domain-driven migration beats big-bang rewrites every time**.

> *"Make the change easy, then make the easy change."* — Kent Beck

---

## Defining Simplicity

### How do we tell truths that might hurt? — Edsger W. Dijkstra (1975)
**Source:** https://www.cs.virginia.edu/~evans/cs655/readings/ewd498.html

Dijkstra's aphorism connects simplicity directly to dependability. In this framing,
simplicity is not a matter of taste or minimalism; it is a precondition for systems that can
be trusted.

> *"Simplicity is prerequisite for reliability."*

### On the Nature of Computing Science — Edsger W. Dijkstra (1984)
**Source:** https://en.wikiquote.org/wiki/Edsger_W._Dijkstra

A useful thesis quote for the presentation: simplicity is hard to create, requires education to
recognise, and is often less marketable than visible complexity.

> *"Simplicity is a great virtue but it requires hard work to achieve it and education to appreciate it. And to make matters worse: complexity sells better."*

### The Emperor's Old Clothes — C. A. R. Hoare (1981)
**Source:** https://doi.org/10.1145/358549.358561

Hoare defines simplicity as the kind of design where defects become obvious rather than hidden.
The quote is especially useful for distinguishing simple design from naive design: the simple
version is harder to reach.

> *"There are two ways of constructing a software design: One way is to make it so simple that there are obviously no deficiencies and the other way is to make it so complicated that there are no obvious deficiencies. The first method is far more difficult."*

Hoare also makes the reliability connection explicit:

> *"The price of reliability is the pursuit of the utmost simplicity."*

### Software Tools — Brian Kernighan and P. J. Plauger (1976)
**Source:** https://en.wikiquote.org/wiki/Brian_Kernighan

Kernighan and Plauger place complexity control at the centre of programming itself. This supports
a definition of simplicity as active management of accidental complexity, not merely short code.

> *"Controlling complexity is the essence of computer programming."*

### The Elements of Programming Style — Brian Kernighan and P. J. Plauger (1978)
**Source:** https://en.wikiquote.org/wiki/Brian_Kernighan

This quote is a practical warning against cleverness. Simple code stays within the future
maintainer's ability to debug and reason about it.

> *"Everyone knows that debugging is twice as hard as writing a program in the first place. So if you're as clever as you can be when you write it, how will you ever debug it?"*

### Simple Made Easy — Rich Hickey (2011)
**Source:** https://www.infoq.com/presentations/Simple-Made-Easy/

Hickey provides a precise software definition: simple means not intertwined. This is useful for
explaining why something can be familiar and easy while still being complex.

> *"Simple things are those which have one role. They fulfill one task, they have one objective, they cover one concept."*

> *"Simple != easy."*

### Dirac: A Scientific Biography — Paul Dirac, quoted by Helge Kragh (1990)
**Source:** https://en.wikiquote.org/wiki/Paul_Dirac

Dirac gives a broader scientific framing: simplicity is explanatory power, the ability to make
difficult things understandable.

> *"The aim of science is to make difficult things understandable in a simpler way."*

### The Evolution of the Physicist's Picture of Nature — Paul Dirac (1963)
**Source:** https://en.wikiquote.org/wiki/Paul_Dirac

Dirac's link between beauty, insight, and progress can support a closing note that simple systems
often feel elegant because their structure reveals the underlying idea.

> *"It seems that if one is working from the point of view of getting beauty in one's equations, and if one has really a sound insight, one is on a sure line of progress."*

Working definition for the presentation: **simplicity is the disciplined removal of accidental
entanglement so that a system can be understood, reasoned about, changed, and trusted. It is not
minimalism, familiarity, or clever brevity. It is hard-won clarity.**

---

## Core Strategy

### Strangler Fig Pattern — Martin Fowler (2024)
**File:** `references/strangler-fig-application-martinfowler.md`
**Source:** https://martinfowler.com/bliki/StranglerFigApplication.html

The foundational pattern for incremental legacy migration. New components are built alongside the old system, gradually taking over functionality until the legacy system can be retired — just as a strangler fig vine eventually replaces its host tree. Key insight: transitional architecture (code that only exists to bridge old and new) looks wasteful but dramatically reduces risk.

---

## Why Rewrites Fail

### Things You Should Never Do, Part I — Joel Spolsky (2000)
**File:** `references/joel-on-software-2000-things-you-should-never-do.md`
**Source:** https://www.joelonsoftware.com/2000/04/06/things-you-should-never-do-part-i/

The definitive case against big-bang rewrites, using Netscape's catastrophic three-year rewrite as the central example. Spolsky argues that old code has hidden value: every "ugly" line often represents a hard-won bug fix for a real edge case. Discarding that code means gifting competitors years of development time.

> *"It's harder to read code than to write it."* — the root cause of the rewrite impulse

Also covers: Borland's failed dBase-for-Windows rewrite (eaten by Microsoft Access), and the Quattro Pro from-scratch disaster.

### Netscape Goes Bonkers — Joel Spolsky (2000)
**File:** `references/joel-on-software-2000-netscape-goes-bonkers.md`
**Source:** https://www.joelonsoftware.com/2000/11/20/netscape-goes-bonkers/

A follow-up post-mortem on Netscape 6.0's release. When the browser finally shipped after three years, it was riddled with bugs, missed standard platform conventions (right-click, keyboard shortcuts, accessibility), and lacked features. Lou Montulli, one of the five original Navigator engineers, confirmed the rewrite decision was a primary reason he resigned.

> *"They shouldn't have rewritten from scratch. They should have done this all in steps. Big chunky steps, fine, but steps."* — Lou Montulli

### TSB Bank Migration Failure (2018)
**Source:** https://icedq.com/resources/case-studies/tsb-bank-data-migration-failure

April 2018: TSB attempted to migrate 5 million customers to a new platform (Proteo4UK) in a single weekend cutover. 1.9 million customers could not access their accounts after go-live. Root cause: the team simultaneously upgraded from Proteo3 to Proteo4 *and* migrated the data in the same event, meaning QA ran against an incomplete system and real missing-data scenarios were never exercised.

> *Note: this source is a secondary case study — verify figures before citing.*

Lesson: never combine a platform upgrade and a data migration into the same cutover. Test against the exact system that will run in production.

### The Myth of the Software Rewrite — Erik Dietrich / DaedTech (2015)
**File:** `references/daedtech-2015-myth-of-the-software-rewrite.md`
**Source:** https://daedtech.com/the-myth-of-the-software-rewrite/

Addresses the psychological and organisational dynamics that drive rewrite demands. The root problem isn't the codebase — it's the group's habits. A rewrite with the same team and the same pressures produces the same mess. The alternative: characterisation tests, incremental refactoring, and paying down technical debt deliberately.

> *"The software is a mess because the group made it a mess, and it'll only get and stay clean if the group cleans it."*

> *"Automated tests are your friend — characterize the system's current behavior with lots of automated tests and then work on refactoring."*

---

## How to Succeed

### Introducing Limbo: A Complete Rewrite of SQLite in Rust — Turso (2024)
**File:** `references/turso-2024-limbo-sqlite-rewrite-dst.md`
**Source:** https://turso.tech/blog/introducing-limbo-a-complete-rewrite-of-sqlite-in-rust

When a rewrite *is* justified (proprietary test suite blocking large changes, unsafe C codebase), how do you do it right? Turso's answer: **Deterministic Simulation Testing (DST) built in from day one**. DST allows testing years of execution with different event orderings, and reproduces any found failure 100% reliably. They also partnered with Antithesis (a deterministic hypervisor) to catch failure modes no internal framework would find — such as io_uring partial-write bugs.

Key lesson: testing confidence is the foundation that makes a rewrite safe.

### Deconstructing the Monolith: Shopify — Shopify Engineering (2019)
**Source:** https://shopify.engineering/shopify-monolith

Shopify decomposed their Rails monolith using DDD subdomains as the boundary strategy — but deliberately chose a **modular monolith** (Rails Engines) over microservices, arguing that distributing a system adds complexity that must be paid for in coordination costs. They used a "vision-guided" approach (DDD first, then engineering) rather than a "spin-off" approach (extract whatever is painful today). Result: 37 components with explicit public entrypoints.

> *"Splitting a single monolithic application into a distributed system of services increases the overall complexity considerably."*

Lesson: DDD gives you the boundaries without forcing a distributed system. The decomposition can happen inside a deployable monolith.

### Online Migrations at Scale — Stripe Engineering (2017)
**Source:** https://stripe.com/blog/online-migrations

Stripe migrated hundreds of millions of Subscriptions objects with zero downtime using a strict **four-phase dual-write pattern**: (1) dual write to old and new tables, (2) change read paths to new table, (3) change write paths exclusively to new table, (4) remove old data. Each phase is independently verifiable and reversible — there is no single cutover moment.

Lesson: any data migration can be made safe by ensuring consistency is maintained at every phase boundary and that each step can be reversed before proceeding.

### Beating the Odds: Khan Academy's Monolith→Services Rewrite — Khan Academy Engineering (2021)
**Source:** https://blog.khanacademy.org/beating-the-odds-khan-academys-successful-monolith→services-rewrite/

Khan Academy migrated ~1 million lines of Python 2 to Go ("Project Goliath"), running from early 2020 to August 2021 (95% of traffic on the new backend at that milestone). The key to finishing: a strict **"direct port only" policy** enforced by a single decision-maker — no new features mixed into the migration. Scope creep is what kills language rewrites.

> *Note: this source was not adversarially verified in the research pass — verify claims before citing.*

### Airbnb SOA Migration — InfoQ / Airbnb Engineering
**Source:** https://www.infoq.com/presentations/airbnb-soa-migration

Airbnb extracted their Rails monolith into Java microservices one endpoint at a time, using dual-read/dual-write with configurable feature gates and a 1%→100% traffic ramp per endpoint. The monolith remained the live fallback throughout. Reported outcomes: ~10,000 deploys per week (vs. hour-long monolith deploys), search page ~3× faster, home description page ~10× faster.

> *Note: this source is a secondary summary (InfoQ presentation) — verify figures before citing.*

### Rebuilding Slack on the Desktop — Slack Engineering (2019)
**Source:** https://slack.engineering/rebuilding-slack-on-the-desktop/

The definitive "Ship of Theseus" case study. Slack replaced their Electron desktop app one module at a time while shipping continuously to users. The first modern component — the emoji picker — shipped more than two years before the full rollout. Three explicit benefits over a big-bang approach: immediate value delivery, sustained focus on continuous improvement, and de-risked release.

> *"Had we waited until the entirety of Slack was rewritten before releasing it, our users would have had a worse day-to-day experience."*

> *"Releasing incrementally allowed us to deliver real value to our customers as soon as possible, helped us stay focused on continuous improvement, and de-risked the release."*

### Rewriting the Heart of Our Sync Engine — Dropbox Engineering (2020)
**Source:** https://dropbox.tech/infrastructure/rewriting-the-heart-of-our-sync-engine

One of the rare justified big-bang rewrites. Dropbox spent ~4 years rewriting their Python desktop sync engine from scratch in Rust because the original data model was structurally wrong: files had no stable identifier preserved across moves, making sharing impossible to add incrementally. The key contribution is their **decision checklist** for when a rewrite is actually warranted:

1. Are incremental improvements truly exhausted?
2. Do you have the resources and organisational health?
3. Can the team accept slower feature development for the duration?

> *"Rust has been a force multiplier for our team… one of the best decisions we made."*

Lesson: a big-bang rewrite is justified only when the architecture is structurally wrong — not just ugly or slow.

### Modernizing Reddit's Comment Backend — Reddit Engineering (2025)
**File:** `references/reddit.html`
**Source:** https://www.reddit.com/r/RedditEng/comments/1mbqto6/

A recent real-world case study of migrating Reddit's highest-throughput write service from Python to Go — with **zero user disruption** and a **50% reduction in p99 write latency**. Key technique: **tap compare** (shadow traffic comparison) for reads, and **sister datastores** for writes (isolated parallel stores for the new service to write to during validation, preventing corruption of production data).

Lessons learned:
- Python ORM had hidden DB optimisations; raw Go queries caused unexpected pressure at ramp-up — monitor DB query patterns early.
- Race conditions produced false tap-compare mismatches; solution: version database updates.
- Relying on production tap compare logs for edge-case discovery is costly; invest in comprehensive local tests using real production data *before* starting shadow traffic.

---

## DDD Concepts to Succeed

- Pilars: Ubiquitous language, strategic patterns and tactical patterns
- Domain, sub-domain and bounded contexts
- Event storming
- Architecture (hexagonal proposed in this presentation) is dissociated from DDD
- Value Objects, Entities and Aggregates
- Context Maps
- Factories and Repositories
- Domain x Application
- Services: Domain and application

- [Azure reference tactical patterns](https://learn.microsoft.com/en-us/azure/architecture/microservices/model/tactical-domain-driven-design)


---

## Internal Reference

### WiiSmile — New Architecture Proposition
**File:** `references/wiismile-new-arch-proposition.md`

The internal document motivating this presentation. Proposes migrating WiiSmile's monolithic PHP stack to a modular, DDD-based architecture (PHP + Go backend, React/Next.js frontend) using the Strangler Fig strategy. Includes a SWOT analysis, technology comparison tables, and a step-by-step migration plan:

1. Establish performance metrics (MTTR, page load times)
2. Increase test coverage before touching anything
3. Identify the least critical module to migrate first
4. Run EventStorming sessions with PO + devs + domain experts
5. Migrate one module at a time, freeze feature additions during migration
6. Verify metrics after each module
7. Deploy continuously — never accumulate a big release

---

## Themes at a Glance

| Theme | Key Sources |
|---|---|
| **Never do big bang** | Joel 2000, Joel Nov 2000, DaedTech |
| **Incremental migration pattern** | Strangler Fig (Fowler), Reddit Eng |
| **Tests as the migration safety net** | DaedTech, Turso, Reddit Eng, WiiSmile |
| **Organisational buy-in & team habits** | DaedTech, WiiSmile |
| **Real-world Go migration case study** | Reddit Eng |
| **When a full rewrite is justified** | Turso (with DST from day one), Dropbox (checklist + broken data model) |
| **Ship of Theseus in production** | Slack desktop rebuild |
| **Data migration pattern** | Stripe four-phase dual-write |
| **DDD decomposition without microservices** | Shopify modular monolith |
| **Scope discipline in language ports** | Khan Academy Goliath |
| **Endpoint-by-endpoint SOA extraction** | Airbnb |
| **Big-bang data migration failure** | TSB Bank 2018 |
| **Simplicity as hard-won reliability** | Dijkstra, Hoare, Kernighan, Hickey, Dirac |
