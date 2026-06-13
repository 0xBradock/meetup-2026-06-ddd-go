# Simple Rewrites Deck Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace the current reveal.js deck with the approved evidence-first presentation: `From Spaghetti to Scalable: Simple rewrites with DDD and Go`.

**Architecture:** Keep the presentation as a self-contained reveal.js `index.html` file, following the existing project pattern. Use horizontal sections for major narrative beats and vertical nested sections for modular stacks. Add only small CSS helpers inside the page; do not introduce new build tooling or split files unless the deck becomes too difficult to maintain.

**Tech Stack:** reveal.js 6, HTML, inline CSS, Mermaid, RevealHighlight, RevealNotes, Vite/npm scripts from `reveal.js/package.json`.

---

## File Structure

- Modify: `reveal.js/index.html`
  - Replace existing presentation content with the approved deck structure.
  - Keep existing reveal.js CSS/script imports unless a change is required.
  - Keep Mermaid support because several slides use diagrams.
  - Add small utility classes for case-study cues, outcome tables, strategy loop, and code comparison.
- Reference only: `docs/superpowers/specs/2026-06-13-simple-rewrites-presentation-design.md`
  - Source of truth for structure and content rules.
- Reference only: `README.md`
  - Source of the working simplicity definition.
- Reference only: `references/*.md` and `references/reddit.html`
  - Source-backed claims for case studies.

Do not edit generated reveal.js `dist/` files. Do not change package dependencies.

## Command Conventions

- Commit this plan file from the root repository: `/home/bradock/Code/pres/2026-ddd-go`.
- Commit deck changes from the nested reveal.js repository: `/home/bradock/Code/pres/2026-ddd-go/reveal.js`.
- Run npm commands from the reveal.js project directory: `/home/bradock/Code/pres/2026-ddd-go/reveal.js`.
- `npm run build:core` validates the reveal.js build assets, not the semantic correctness of `index.html`. Use it as a regression check that the project still builds. Use the final local presentation smoke test to validate deck markup, navigation, Mermaid, and content.

## Chunk 1: Deck Skeleton and Opening

### Task 1: Replace Metadata and Add Minimal Styling

**Files:**
- Modify: `reveal.js/index.html:1-20`

- [ ] **Step 1: Update document title**

Change the `<title>` to:

```html
<title>From Spaghetti to Scalable</title>
```

- [ ] **Step 2: Replace inline CSS with deck helpers**

Keep the existing Mermaid and blockquote support, then add helper classes:

```css
.reveal .mermaid svg { max-height: 65vh; }
.reveal blockquote { font-style: italic; border-left: 4px solid #42affa; padding-left: 1rem; }
.reveal .tag { display: inline-block; background: #42affa22; border: 1px solid #42affa55; border-radius: 4px; padding: 2px 10px; font-size: 0.6em; margin: 2px; }
.reveal .cue { font-size: 0.7em; text-transform: uppercase; letter-spacing: 0.08em; }
.reveal .success { color: #69db7c; }
.reveal .failure { color: #ff6b6b; }
.reveal section.success-case { border-top: 0.35rem solid #69db7c; }
.reveal section.failure-case { border-top: 0.35rem solid #ff6b6b; }
.reveal .outcome { width: 100%; font-size: 0.55em; }
.reveal .outcome th { color: #42affa; }
.reveal .loop { display: grid; grid-template-columns: repeat(7, 1fr); gap: 0.4rem; align-items: stretch; }
.reveal .loop-step { border: 1px solid #42affa66; border-radius: 0.4rem; padding: 0.6rem; font-size: 0.5em; background: #42affa11; }
.reveal .code-compare { display: grid; grid-template-columns: 1fr 1fr; gap: 1rem; }
.reveal .code-compare pre { font-size: 0.42em; }
```

- [ ] **Step 3: Run a build regression check**

Run from `reveal.js/`: `npm run build:core`

Expected: Build completes without TypeScript/Vite errors. This does not validate deck markup or slide semantics; it only confirms the project still builds after editing `index.html`.

- [ ] **Step 4: Commit**

```bash
git add index.html
git commit -m "Update deck metadata and styles"
```

### Task 2: Implement Title and Simplicity Vertical Stack

**Files:**
- Modify: `reveal.js/index.html:24-37`

- [ ] **Step 1: Replace the first slide with title content**

Use:

```html
<section>
	<h1>From Spaghetti to Scalable</h1>
	<p>Simple rewrites with DDD and Go</p>
	<br />
	<blockquote>
		"Make the change easy, then make the easy change."
		<br /><small>— Kent Beck</small>
	</blockquote>
</section>
```

- [ ] **Step 2: Add the `Simplicity` vertical stack after the title**

Use one top slide plus one quote per slide and the working definition:

```html
<section>
	<section>
		<h2>Simplicity</h2>
	</section>
	<section>
		<blockquote>
			"Simplicity is prerequisite for reliability."
			<br /><small>— Edsger W. Dijkstra</small>
		</blockquote>
	</section>
	<section>
		<blockquote>
			"There are two ways of constructing a software design: One way is to make it so simple that there are obviously no deficiencies and the other way is to make it so complicated that there are no obvious deficiencies."
			<br /><small>— C. A. R. Hoare</small>
		</blockquote>
	</section>
	<section>
		<blockquote>
			"Simple things are those which have one role."
			<br /><small>— Rich Hickey</small>
		</blockquote>
		<p><strong>Simple != easy.</strong></p>
	</section>
	<section>
		<h3>Working Definition</h3>
		<p>Simplicity is the disciplined removal of accidental entanglement so that a system can be understood, reasoned about, changed, and trusted.</p>
		<p>It is not minimalism, familiarity, or clever brevity. It is hard-won clarity.</p>
	</section>
</section>
```

- [ ] **Step 3: Verify the deck still builds**

Run from `reveal.js/`: `npm run build:core`

Expected: Build completes without errors.

- [ ] **Step 4: Commit**

```bash
git add index.html
git commit -m "Add title and simplicity stack"
```

### Task 3: Add Modernization Definition Slide

**Files:**
- Modify: `reveal.js/index.html` after the simplicity stack

- [ ] **Step 1: Add a horizontal `Modernization` slide**

Use:

```html
<section>
	<h2>Modernization</h2>
	<p>Changing an existing system so it can keep serving the business under today's constraints.</p>
	<div class="code-compare">
		<div>
			<h3>Port</h3>
			<p>Same external behavior, new technical substrate.</p>
		</div>
		<div>
			<h3>Rewrite</h3>
			<p>Preserve the business capability, but change design boundaries and implementation to reduce accidental complexity.</p>
		</div>
	</div>
</section>
```

- [ ] **Step 2: Verify build**

Run from `reveal.js/`: `npm run build:core`

Expected: Build completes without errors.

- [ ] **Step 3: Commit**

```bash
git add index.html
git commit -m "Define modernization strategies"
```

## Chunk 2: Evidence Use-Case Stacks

### Task 4: Add TypeScript Port Stack

**Files:**
- Modify: `reveal.js/index.html` after Modernization slide

- [ ] **Step 1: Add top and summary slides**

Create a vertical stack with `section.success-case` on the top slide:

```html
<section>
	<section class="success-case">
		<p class="cue success">🚀 Port</p>
		<h2>TypeScript Native</h2>
		<p>Same language and tooling experience, new native implementation substrate.</p>
	</section>
	<section>
		<h3>Why Modernize?</h3>
		<ul>
			<li>Large TypeScript codebases need faster editor startup and type checking.</li>
			<li>The goal is better developer experience at scale.</li>
			<li>Microsoft reported about 10x faster builds, faster editor startup, and lower memory usage.</li>
			<li>The implementation target is Go.</li>
		</ul>
	</section>
```

- [ ] **Step 2: Add process and safety slides**

Continue the same stack:

```html
	<section>
		<h3>Process</h3>
		<ul>
			<li>Port the compiler and tools to native Go.</li>
			<li>Preview command-line type checking first.</li>
			<li>Move toward project builds and language service parity.</li>
		</ul>
	</section>
	<section>
		<h3>Safety Mechanism</h3>
		<ul>
			<li>Behavior parity is the product promise.</li>
			<li>TypeScript 6 remains the JavaScript line.</li>
			<li>TypeScript 7 becomes native once mature enough.</li>
		</ul>
	</section>
```

- [ ] **Step 3: Add outcome table and close stack**

Use:

```html
	<section>
		<h3>Outcome</h3>
		<table class="outcome">
			<thead><tr><th>What worked</th><th>What to watch</th></tr></thead>
			<tbody>
				<tr><td>Clear technical target: faster tooling.</td><td>Compatibility and API parity must be protected.</td></tr>
				<tr><td>Reported order-of-magnitude speedups on large projects.</td><td>Existing projects may remain on TypeScript 6 until TypeScript 7 fits their constraints.</td></tr>
				<tr><td>Go supports the operational goal with native binaries and strong tooling.</td><td>The port succeeds only if behavior remains familiar.</td></tr>
			</tbody>
		</table>
	</section>
</section>
```

- [ ] **Step 4: Verify build and commit**

Run from `reveal.js/`: `npm run build:core`

Expected: Build completes without errors.

Commit:

```bash
git add index.html
git commit -m "Add TypeScript port case study"
```

### Task 5: Add Reddit Rewrite Stack

**Files:**
- Modify: `reveal.js/index.html` after TypeScript stack

- [ ] **Step 1: Add top and summary slides**

Create a `Rewrite: Reddit` vertical success stack:

```html
<section>
	<section class="success-case">
		<p class="cue success">🚀 Rewrite</p>
		<h2>Reddit Comment Backend</h2>
		<p>Incremental migration of a high-throughput write path from Python to Go.</p>
	</section>
	<section>
		<h3>Why Modernize?</h3>
		<ul>
			<li>Comments are a high-throughput, user-visible write path.</li>
			<li>The target was lower latency without user disruption.</li>
			<li>The new service was written in Go.</li>
		</ul>
	</section>
```

- [ ] **Step 2: Add tap-compare process diagram**

Use this slide:

```html
	<section>
		<h3>Process: Tap Compare for Writes</h3>
		<div class="mermaid">
sequenceDiagram
    actor Client
    participant Go as Go Service (new)
    participant Py as Python Service (legacy)
    participant PD as Production Stores<br/>(Postgres + Memcached + Redis)
    participant SD as Sister Stores<br/>(isolated — Go only)

    Client->>Go: Write request (small % of traffic)
    Go->>Py: Delegate real write
    Py->>PD: Write to production
    Py-->>Go: Response
    Go->>SD: Write independently (no production impact)
    Go->>Go: Compare production vs sister stores
    Note over Go: Log differences — never return new response to user
    Go-->>Client: Return legacy response
		</div>
	</section>
```

- [ ] **Step 3: Add safety mechanism slide**

Use:

```html
	<section>
		<h3>Safety Mechanism</h3>
		<ul>
			<li>Tap compare for write-path behavior validation.</li>
			<li>Sister datastores for write-path validation.</li>
			<li>The legacy response remained the user-facing truth during validation.</li>
		</ul>
	</section>
```

- [ ] **Step 4: Add outcome table and close stack**

Use:

```html
	<section>
		<h3>Outcome</h3>
		<table class="outcome">
			<thead><tr><th>What worked</th><th>What to watch</th></tr></thead>
			<tbody>
				<tr><td>Zero user disruption during migration.</td><td>Shadow validation does not remove the need for strong local tests.</td></tr>
				<tr><td>p99 write latency was reduced by about 50%.</td><td>Raw Go queries exposed DB pressure hidden by the Python ORM.</td></tr>
				<tr><td>Sister datastores protected production writes.</td><td>Race conditions created false tap-compare mismatches.</td></tr>
				<tr><td>Tap compare gave production confidence before cutover.</td><td>Finding edge cases through production logs is expensive.</td></tr>
			</tbody>
		</table>
	</section>
</section>
```

- [ ] **Step 5: Verify build and commit**

Run from `reveal.js/`: `npm run build:core`

Expected: Build completes without errors.

Commit:

```bash
git add index.html
git commit -m "Add Reddit rewrite case study"
```

### Task 6: Add Turso Rewrite Stack

**Files:**
- Modify: `reveal.js/index.html` after Reddit stack

- [ ] **Step 1: Add top and summary slides**

Create this opening stack:

```html
<section>
	<section class="success-case">
		<p class="cue success">🚀 Rewrite</p>
		<h2>Turso Limbo</h2>
		<p>A SQLite-compatible rewrite in Rust, designed with reliability testing from day one.</p>
	</section>
	<section>
		<h3>Why Modernize?</h3>
		<ul>
			<li>libSQL started as a SQLite fork and reached changeability limits.</li>
			<li>Some key SQLite test infrastructure, such as TH3, is proprietary.</li>
			<li>C made invasive evolution riskier for the target changes.</li>
		</ul>
	</section>
```

- [ ] **Step 2: Add process slide**

Use:

```html
	<section>
		<h3>Process</h3>
		<ul>
			<li>Rewrite SQLite-compatible behavior in Rust.</li>
			<li>Design for async I/O from the beginning.</li>
			<li>Keep compatibility visible by verifying generated bytecode against SQLite.</li>
			<li>Regularly fuzz inputs to find edge cases.</li>
		</ul>
	</section>
```

- [ ] **Step 3: Add safety mechanism slide**

Use:

```html
	<section>
		<h3>Safety Mechanism</h3>
		<ul>
			<li>Deterministic Simulation Testing built in from day one.</li>
			<li>Failures can be reproduced reliably.</li>
			<li>Antithesis fault injection found rare io_uring partial-write failures.</li>
		</ul>
	</section>
```

- [ ] **Step 4: Add outcome table**

Use:

```html
	<section>
		<h3>Outcome</h3>
		<table class="outcome">
			<thead><tr><th>What worked</th><th>What to watch</th></tr></thead>
			<tbody>
				<tr><td>Reliability strategy was part of the architecture from day one.</td><td>A full rewrite needs exceptional justification.</td></tr>
				<tr><td>DST made rare failures reproducible.</td><td>Compatibility with SQLite behavior remains a large burden.</td></tr>
				<tr><td>Antithesis found failure modes internal tests may have missed.</td><td>The testing investment is not optional; it is the price of the rewrite.</td></tr>
			</tbody>
		</table>
	</section>
</section>
```

- [ ] **Step 5: Verify build and commit**

Run from `reveal.js/`: `npm run build:core`

Expected: Build completes without errors.

Commit:

```bash
git add index.html
git commit -m "Add Turso rewrite case study"
```

### Task 7: Add Failure Case Stacks

**Files:**
- Modify: `reveal.js/index.html` after Turso stack

- [ ] **Step 1: Add Netscape failure stack**

Create a vertical stack with these exact content requirements:

- Top slide: `💥 Rewrite` cue, `Netscape`, and `The canonical big-bang rewrite warning.`
- Summary slide: `Rewritten from scratch`, `three-year competitive gap`, `old behavior rediscovered late`.
- Process slide: `Stop improving the old product`, `build the replacement in a tunnel`, `ship when parity feels close`.
- Safety slide: `No continuous user feedback loop`, `inadequate testing in the referenced story`, `hidden behavior discarded`.
- Outcome table with columns `What was attempted` and `What failed`:
  - `Start fresh` / `Lost years of market feedback`
  - `Replace old code wholesale` / `Discarded hard-won edge cases`
  - `Ship the new version after the tunnel` / `Users saw bugs and missing platform conventions`
- Bottom failure slide: `The big bang removed the feedback loop for too long.`

- [ ] **Step 2: Add Borland failure stack**

Create a vertical stack with these exact content requirements:

- Top slide: `💥 Rewrite` cue, `Borland`, and `Feature parity is a trap.`
- Summary slide: mention dBase-for-Windows and Quattro Pro as two brief examples from Joel Spolsky.
- Process slide: `Rewrite large products`, `underestimate parity`, `compete while not shipping enough value`.
- Safety slide: `No visible incremental replacement path in the cautionary story`, `scope and timing risk dominated`.
- Outcome table with columns `What was attempted` and `What failed`:
  - `Create a Windows-era replacement` / `Microsoft Access captured the market window`
  - `Rewrite Quattro Pro from scratch` / `The rewritten product surprised users with how few features it had`
  - `Treat parity as implementation detail` / `Parity became the project`
- Bottom failure slide: `A rewrite that misses the market is not modernization.`

- [ ] **Step 3: Add DaedTech failure-pattern stack**

Create a vertical stack with these exact content requirements:

- Top slide: `💥 Rewrite` cue, `Same team, same mess`, and `A new codebase is not a new operating model.`
- Summary slide: based on DaedTech's critique that the group created the mess and must clean it.
- Process slide: `Same deadlines`, `same shortcuts`, `same lack of tests`, `new repository`.
- Safety slide: `Change team habits`, `add characterization tests`, `make debt work explicit`, `keep refactoring continuous`.
- Outcome table with columns `What was attempted` and `What failed`:
  - `Escape the old codebase` / `The causes of complexity remained`
  - `Promise to do it right this time` / `Same pressures recreated the mess`
  - `Treat rewrite as reset button` / `No sustained test and refactoring discipline`
- Bottom failure slide: `A new codebase does not repair the system that produced the old one.`

- [ ] **Step 4: Verify build and commit**

Run from `reveal.js/`: `npm run build:core`

Expected: Build completes without errors.

Commit:

```bash
git add index.html
git commit -m "Add rewrite failure case studies"
```

## Chunk 3: Synthesis and Verification

### Task 8: Add Simple Rewrite Loop and DDD Slides

**Files:**
- Modify: `reveal.js/index.html` after failure stacks

- [ ] **Step 1: Add loop slide**

Use `.loop` and `.loop-step` classes to render:

1. Define the business capability
2. Pick the simplest target architecture
3. Add characterization tests around current behavior
4. Find domain seams
5. Migrate one thin slice
6. Measure
7. Repeat

- [ ] **Step 2: Add DDD slide**

Show DDD as a tool for two loop steps:

- Define capability: ubiquitous language and domain conversations.
- Find domain seam: bounded contexts, EventStorming, consistency boundaries.

- [ ] **Step 3: Verify build and commit**

Run from `reveal.js/`: `npm run build:core`

Expected: Build completes without errors.

Commit:

```bash
git add index.html
git commit -m "Add simple rewrite loop and DDD synthesis"
```

### Task 9: Add Target Architecture and Go Slides

**Files:**
- Modify: `reveal.js/index.html` after DDD slide

- [ ] **Step 1: Add target architecture slide**

Include:

- Start with a domain-oriented modular monolith.
- Deploy together until splitting is proven necessary.
- Evolve toward 4-12 coarse services with a simple template when separation pays for itself.
- Do not start with microservices by default.

- [ ] **Step 2: Add Go slide**

Include:

- DDD helps choose the right slice.
- Go helps keep the new slice operationally simple.
- Small language, standard library, standard tooling, static binaries, fast builds, straightforward deployment, clear concurrency, typical backend performance, Go 1 compatibility.
- Avoid claiming Go is the right choice for extreme-performance domains.

- [ ] **Step 3: Verify build and commit**

Run from `reveal.js/`: `npm run build:core`

Expected: Build completes without errors.

Commit:

```bash
git add index.html
git commit -m "Add target architecture and Go synthesis"
```

### Task 10: Add BorrowBook Code Example and Closing

**Files:**
- Modify: `reveal.js/index.html` after Go slide and before closing scripts

- [ ] **Step 1: Add `BorrowBook` before/after slide**

Use two small Go snippets. Keep them compact.

Before snippet:

```go
func BorrowBook(w http.ResponseWriter, r *http.Request) {
    member := loadMember(r.FormValue("member_id"))
    copy := loadCopy(r.FormValue("copy_id"))

    if member.Blocked || member.ActiveLoans >= 5 || copy.Borrowed {
        http.Error(w, "cannot borrow", http.StatusConflict)
        return
    }

    copy.Borrowed = true
    copy.BorrowedBy = member.ID
    copy.DueAt = time.Now().AddDate(0, 0, 21)
    saveCopy(copy)
}
```

After snippet:

```go
func (h BorrowBookHandler) Handle(ctx context.Context, cmd BorrowBook) error {
    member := h.members.Get(ctx, cmd.MemberID)
    copy := h.copies.Get(ctx, cmd.CopyID)

    if err := member.CanBorrow(); err != nil {
        return err
    }

    loan := copy.MarkBorrowedBy(member.ID, h.clock.Today())
    return h.loans.Save(ctx, loan)
}
```

- [ ] **Step 2: Add closing slide**

Use:

```html
<section>
	<h2>Rewrites should be simple.</h2>
	<p>Preserve the business capability. Reduce accidental entanglement. Keep feedback alive.</p>
</section>
```

- [ ] **Step 3: Verify build and commit**

Run from `reveal.js/`: `npm run build:core`

Expected: Build completes without errors.

Commit:

```bash
git add index.html
git commit -m "Add code example and closing"
```

### Task 11: Final Deck Verification

**Files:**
- Verify: `reveal.js/index.html`

- [ ] **Step 1: Run final build**

Run from `reveal.js/`: `npm run build:core`

Expected: Build completes without errors.

- [ ] **Step 2: Run a local presentation smoke test**

Run from `reveal.js/`: `npm start -- --host 127.0.0.1 --port 8000`

Expected: Vite serves the deck at `http://127.0.0.1:8000/`.

Manual checks:

- Title slide shows `From Spaghetti to Scalable`.
- Vertical navigation works for `Simplicity` and all use-case stacks.
- Required use-case order is TypeScript, Reddit, Turso, Netscape, Borland, Same team/same mess.
- Positive examples stop at `What worked` / `What to watch` outcome tables.
- Failure examples include `What was attempted` / `What failed` tables and bottom failure slides.
- Mermaid diagrams render after Reveal initialization.
- Highlighted Go code appears legible.
- Visible slide text is English.
- No WiiSmile-specific content appears in visible slides.
- Case-study claims are source-backed or clearly framed as interpretation.

- [ ] **Step 3: Stop dev server**

Stop the `npm start` process after manual verification.

- [ ] **Step 4: Check final git status**

Run: `git status --short`

Expected: no unexpected modified files besides any intentionally uncommitted user files.

- [ ] **Step 5: Commit any final fixes**

If verification required fixes:

```bash
git add index.html
git commit -m "Polish presentation deck"
```
