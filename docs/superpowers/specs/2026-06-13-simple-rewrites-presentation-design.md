# Simple Rewrites Presentation Design

## Context

The presentation is for the Meetup event "Du spaghetti au scalable : refonte applicative avec DDD et Go". The visible deck title is English: `From Spaghetti to Scalable: Simple rewrites with DDD and Go`. The deck will be built with reveal.js and optimized for a 45-minute talk plus 15 minutes of Q&A.

The primary audience is experienced application developers and tech leads who have lived with legacy systems. The deck should help them turn rewrite pressure into a simple, repeatable modernization strategy.

The one sentence a developer should remember is:

> Rewrites should be simple.

## Language Rules

- Visible slide content is in English.
- Oral delivery will be in French.
- Reveal speaker notes or asides may be in French.
- Public content must not use WiiSmile-specific domains, internal examples, or private implementation details.

## Narrative Approach

Use an evidence-first modular structure. Start by defining simplicity and modernization, then walk through movable use-case stacks. After the evidence, synthesize the simple rewrite strategy with DDD and Go.

The main horizontal flow is:

1. Title: `From Spaghetti to Scalable: Simple rewrites with DDD and Go`
2. `Simplicity` vertical stack
3. `Modernization` slide defining `Port` and `Rewrite`
4. Modular use-case stacks
5. Synthesis: simple rewrite loop, DDD, Go
6. Closing: `Rewrites should be simple`

## Simplicity Stack

The `Simplicity` stack is a vertical reveal.js section.

Slides:

1. Top slide containing only `Simplicity`
2. Dijkstra quote: `Simplicity is prerequisite for reliability.`
3. Hoare quote: the two ways of constructing software design, contrasting obvious absence of deficiencies with hidden deficiencies
4. Rich Hickey quote: simple things have one role / are not intertwined, plus `Simple != easy`
5. Working definition from `README.md`: simplicity is the disciplined removal of accidental entanglement so a system can be understood, reasoned about, changed, and trusted. It is not minimalism, familiarity, or clever brevity. It is hard-won clarity.

After this stack, the deck uses the word `simple` according to this definition without repeatedly re-explaining the distinction between simple and easy.

## Modernization Slide

After `Simplicity`, add one `Modernization` slide.

Definition:

> Modernization is changing an existing system so it can keep serving the business under today's constraints.

Split modernization into two sibling strategies:

- `Port`: same external behavior, new technical substrate.
- `Rewrite`: preserve the business capability, but change design boundaries and implementation to reduce accidental complexity.

This slide leads into concrete use cases rather than an abstract definition stack.

## Use-Case Stack Template

Each use case is a vertical stack so it can be moved independently.

Standard structure:

1. Top slide: `Port` or `Rewrite` plus company/project name, with visual cue.
2. Summary slide: context, modernization goal, why the change was needed.
3. Process slide or slides: how they proceeded.
4. Safety mechanism slide: tests, shadow traffic, compatibility/versioning plan, deterministic simulation testing, characterization tests, or absence of a safety mechanism.
5. Outcome table.
6. Bottom failure slide only for true cautionary cases.

Positive examples stop at the outcome table. Do not add a bottom failure slide for positive examples.

Outcome table structure:

- Positive examples: `What worked` / `What to watch`
- Failure examples: `What was attempted` / `What failed`
- Each row should be source-backed or explicitly framed as the presentation's interpretation.
- Prefer concrete evidence: metric, mechanism, quote, observed failure mode, or documented trade-off.
- Avoid unsupported generalized claims.

Visual cues:

- Validated patterns: green styling with a small 🚀 rocket accent.
- Cautionary examples: red styling with a visible 💥 explosion accent.
- Emojis are visual accents, not the main slide title.
- Styling should remain minimal in this phase; detailed visual polish can happen later.

## Initial Use Cases

Initial use-case order:

1. 🚀 `Port: TypeScript`
2. 🚀 `Rewrite: Reddit`
3. 🚀 `Rewrite: Turso`
4. 💥 `Rewrite: Netscape`
5. 💥 `Rewrite: Borland`
6. 💥 `Rewrite: Same team, same mess`

Additional use cases may be added later using the same stack template.

This ordering is intentional: positive examples establish what controlled modernization looks like before failure examples expose the anti-patterns. The synthesis then extracts the repeatable loop from both sets of evidence.

### TypeScript Native Port

Purpose: show the `Port` category.

Core points:

- Microsoft announced a native port of the TypeScript compiler and tooling to Go.
- The goal is behavior parity with a new technical substrate.
- The expected benefits are about 10x faster builds, faster editor startup, and lower memory usage.
- TypeScript 6 remains the JavaScript implementation while TypeScript 7 becomes the native implementation once mature.

Use this case early to show that not every rewrite-like effort is a product rewrite. Some are ports with a clear technical target and compatibility plan.

### Reddit Comment Backend

Purpose: show a successful incremental business-system rewrite.

Core points:

- Reddit migrated a high-throughput comment write service from Python to Go.
- They used tap compare for validation and sister datastores for write safety.
- They achieved zero user disruption and a large p99 write-latency improvement.
- Watch-outs included hidden ORM optimizations, DB pressure, race conditions, and expensive edge-case discovery through production logs.

### Turso Limbo

Purpose: show when a full rewrite can be justified and what makes it safer.

Core points:

- Turso pursued a SQLite rewrite in Rust after a fork reached limits.
- Some key SQLite test infrastructure, such as TH3, is proprietary, and C made some changes risky.
- Deterministic Simulation Testing was built in from day one.
- Antithesis fault injection found rare failure modes.

Use this case to clarify that full rewrites are not categorically impossible, but they require exceptional safety mechanisms and justification.

### Netscape

Purpose: show the canonical big-bang rewrite failure.

Core points:

- Netscape rewrote from scratch and lost years of competitive ground.
- The failure pattern is long tunnel delivery, hidden behavior rediscovery, and insufficient feedback.
- Bottom 💥 slide: `The big bang removed the feedback loop for too long.`

### Borland

Purpose: show scope, timing, and feature-parity failure.

Core points:

- Treat Borland as one cautionary evidence stack with two brief examples from Joel Spolsky: dBase-for-Windows and Quattro Pro.
- The shared lesson is that rewrite scope and feature parity are easy to underestimate, and late delivery gives competitors time to win.
- The failure pattern is underestimating feature parity and business timing.
- Bottom 💥 slide: `A rewrite that misses the market is not modernization.`

### Same Team, Same Mess

Purpose: show that rewrites do not fix team habits by themselves.

Core points:

- Based on DaedTech's rewrite critique.
- If the same team operates under the same pressures without changing habits, the new codebase tends to become the same mess.
- Bottom 💥 slide: `A new codebase does not repair the system that produced the old one.`

## Synthesis Section Order

After the use-case stacks, use separate horizontal slides for the synthesis. Keep the section compact so it feels like the conclusion drawn from the evidence.

Slide order:

1. `The simple rewrite loop`
2. `DDD helps define capability and seams`
3. `The simplest target architecture`
4. `Go as a simple target candidate`
5. `One slice: BorrowBook`

## Synthesis: Simple Rewrite Loop

After the use cases, synthesize the evidence into a repeatable loop:

1. Define the business capability.
2. Pick the simplest target architecture.
3. Add characterization tests around current behavior.
4. Find domain seams.
5. Migrate one thin slice.
6. Measure.
7. Repeat.

This must be presented as a loop, not a one-time checklist.

## DDD Role

DDD is not presented as a full tutorial. It is used as a simplification tool for two steps in the loop:

- `Define the business capability`: use ubiquitous language and domain conversations to identify what must be preserved or improved.
- `Find domain seams`: use bounded contexts, EventStorming, and consistency boundaries to choose a migration slice.

Avoid heavy DDD content such as CQRS and event sourcing in the main path. Mention these only as optional escalation patterns when needed.

The `BorrowBook` code example supports this section. It demonstrates how to define one business capability and carve a seam around it before changing implementation details.

## Target Architecture

The simplest target architecture starts as a domain-oriented modular monolith.

Guidance:

- Use DDD to define coarse business modules with clear boundaries.
- Deploy together until there is a proven reason to split.
- If more separation is needed, evolve toward a service-based architecture with roughly 4-12 coarse services and a simple shared template.
- Do not start with microservices by default.

## Go Role

Go appears under `Pick the simplest target`.

Core claim:

> DDD helps us choose the right slice. Go helps us keep the new slice operationally simple.

Go is a strong default candidate for rewritten backend parts when the current stack is hard to maintain, outdated, too operationally heavy, or too slow for typical backend needs.

The argument for Go is operational simplicity and maintainability, not universal raw performance.

Emphasize:

- Small language surface
- Orthogonal concepts
- Strong standard library
- Standard tooling
- Static binaries
- Fast builds
- Straightforward deployment
- Clear concurrency model
- Good performance for typical backend services
- Go 1 compatibility promise since 2012

Be explicit that extreme-performance domains, such as some trading or communication systems, may require different choices.

Go serves the migration strategy; it is not the strategy itself.

## Code Example

Include one small code example only.

Domain:

- Fictional neutral domain: library loans.

Slice:

- `BorrowBook`

Purpose:

- Demonstrate DDD as simplification.
- Show a tangled operation becoming a small use case around domain concepts.
- Do not teach Go syntax.

Suggested before/after:

- Before: one handler checks member status, book copy availability, dates, persistence, and HTTP response together.
- After: a `BorrowBook` use case coordinates domain rules such as `Member.CanBorrow()` and `Copy.MarkBorrowed()`.

Placement:

- Put this example after the DDD and target-architecture slides, before the closing.
- Use it to demonstrate `Define the business capability`, `Find domain seams`, and `Migrate one thin slice`.

## Closing

End with the one-sentence takeaway:

> Rewrites should be simple.

The closing should connect simple back to the opening definition: simple rewrites reduce accidental entanglement while preserving business capability and feedback.
