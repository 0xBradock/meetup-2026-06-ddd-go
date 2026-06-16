# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository Layout

This repository is a working copy of **reveal.js** (the HTML presentation framework). The project root contains two empty placeholder files (`AGENT.md`, `README.md`) and two subdirectories:

- `reveal.js/` — the full reveal.js source (core library + React wrapper)
- `references/` — reference documents (architecture notes, examples)

All meaningful work happens inside `reveal.js/`.

## Two Sub-Projects

### Core library (`reveal.js/`)
Built with Vite + TypeScript. Entry point: `js/index.ts`. The compiled output goes to `dist/`.

- `js/reveal.js` — main Reveal factory function and public API
- `js/controllers/` — feature controllers (keyboard, fragments, auto-animate, scroll view, etc.)
- `js/utils/` — shared utilities (color, device detection, constants, DOM helpers)
- `css/reveal.scss` + `css/theme/*.scss` — core styles and built-in themes
- `plugin/` — official plugins (highlight, markdown, math, notes, search, zoom), each with its own `vite.config.ts`

### React wrapper (`reveal.js/react/`)
A separate npm package (`@revealjs/react`). Entry: `react/src/index.ts`. Built with Vite; tests use Vitest.

- `react/src/components/` — `Deck`, `Slide`, `Stack`, `Fragment`, `Code`, `Markdown` components (kebab-case filenames)
- `react/src/utils/` — `slide-attributes.ts`, `markdown.ts`
- `react/src/reveal-context.ts` — React context carrying the Reveal instance
- `react/src/types.ts` — shared TypeScript types
- Component tests are colocated as `*.test.tsx`; test setup in `react/src/__tests__/setup.ts`
- Demo app: `react/demo/`

## Commands

All commands are run from `reveal.js/` unless noted.

```bash
# Dev server (port 8000 by default, override with --port)
npm start

# Build core library only
npm run build:core

# Full build (core + all plugins + themes)
npm run build

# Core tests (QUnit via Puppeteer — requires a browser)
# Runs ALL test/*.html files; no CLI flag to run a single file.
# To isolate one test, open it via the dev server: npm start, then visit http://localhost:8000/test/<file>.html
npm test

# React wrapper — from reveal.js/ root
npm run react:test    # vitest run (all)
npm run react:build
npm run react:demo    # dev server for the demo app

# Or cd into react/ and run directly
npm test                                          # vitest run (all)
npm run test:watch                                # vitest watch mode
npx vitest run src/components/deck.test.tsx       # single file
npx vitest run --reporter=verbose -t "test name"  # by test name pattern
```

## Code Style

`.prettierrc` enforces: **tabs** for indentation, **100-character** line width, trailing commas. Run `npx prettier --write <file>` to format.

## Architecture Notes

**Reveal factory pattern**: `reveal.js` exports a single factory function. Calling it returns an instance object (`Reveal`). There is no class; the instance API is built up by closures over the factory scope.

**Controller pattern**: Each UI feature (keyboard, fragments, backgrounds, etc.) is a standalone controller object instantiated by the Reveal factory and coordinated through it. Controllers do not communicate directly with each other.

**React wrapper lifecycle**: `Deck` creates exactly one `Reveal` instance on mount and destroys it on unmount. It must stay safe under React StrictMode (no double-init). `Reveal.sync()` is expensive — it is called only when the rendered slide structure changes (slides added/removed/reordered), never for ordinary content updates inside an existing slide.

**React config updates**: Config is shallow-compared. A new config object with identical shallow values must not trigger `configure()`. After `configure()` Reveal performs its own sync, so the wrapper must not call `sync()` immediately afterward.

**Sync responsibility split**:
- `Deck` — Reveal lifecycle, config, event wiring, structure-level sync
- `Slide` — per-slide `data-*` attribute mapping, calls `syncSlide()` only on attribute change after mount
- `Markdown` — markdown parsing, separator/notes/comment-attribute support; keeps DOM post-processing local
- `Code` — explicit code block rendering and Reveal highlight integration

**Plugin builds**: Each plugin under `plugin/` has its own `vite.config.ts` and is built as a separate library entry. The root `npm run build` invokes all of them in sequence.

**Test runner**: Core tests use QUnit loaded in HTML fixtures under `test/`. `scripts/test.js` spins up a Vite dev server, then runs each `test/*.html` file through Puppeteer via `node-qunit-puppeteer`.
