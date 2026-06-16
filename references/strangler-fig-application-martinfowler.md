# Strangler Fig Application

**Author:** Martin Fowler
**Date:** 22 August 2024
**Source:** https://martinfowler.com/bliki/StranglerFigApplication.html
**Tags:** application architecture, legacy modernization

---

## Overview

Martin Fowler describes the Strangler Fig pattern as an approach to modernizing legacy software systems, inspired by the botanical process of strangler fig vines gradually replacing their host trees.

## The Biological Metaphor

During a 2001 vacation in Queensland's rain forests, Fowler observed strangler figs — vines that germinate in tree crevices, gradually drawing nutrients from the host while growing toward sunlight and soil. Eventually, the host tree dies while the fig becomes self-sustaining. This natural process became an analogy for gradual software modernization.

## Why Wholesale Replacement Fails

Organizations facing outdated systems often attempt complete rewrites. However, this approach typically fails because:

- Modern users cannot wait for replacements to complete
- Hidden behavioral details prove difficult to replicate
- Much existing behavior shouldn't be rebuilt at all

## The Gradual Modernization Approach

Rather than replacement, the Strangler Fig method involves building new components alongside legacy systems. New features develop independently while gradually migrating functionality from old to new code.

## Four Essential Activities

According to Cartwright, Horn, and Lewis, incremental modernization requires:

1. Clarifying desired outcomes and maintaining alignment
2. Identifying seams to decompose systems into manageable pieces
3. Delivering replacements with reduced risk through smaller components
4. Implementing organizational changes to support sustainable development

## Key Considerations

Transitional architecture — temporary code enabling coexistence — appears wasteful but reduces risk and enables earlier value realization. Critically, organizational culture must evolve alongside technical systems; otherwise, new systems become similarly brittle.

## Further Resources

- [Patterns of Legacy Displacement](https://martinfowler.com/articles/patterns-legacy-displacement/) — additional techniques and case studies
