---
name: refactor
description: 'Use when refactoring, simplifying, cleaning up, restructuring, or rewriting code for clarity. Prioritizes the simplest correct final result over minimal diff. Trigger words: refactor, simplify, rewrite, clean up, restructure.'
argument-hint: 'Describe what should be refactored and any behavior that must remain unchanged.'
user-invocable: true
---

# Refactor

Refactor code toward the simplest correct final result, not the smallest diff.

## When to Use

- Refactoring code for clarity or maintainability
- Simplifying complex logic or structure
- Cleaning up duplicated, dead, or awkward code
- Restructuring modules, functions, or classes without changing behavior
- Rewriting an implementation when the rewritten version is clearly better than an incremental patch

## Mindset

- Ignore scope minimization when it blocks a cleaner result.
- Rewrite, restructure, rename, merge, split, or delete freely if that produces a better final design.
- Prefer deleting code over preserving unnecessary code.
- Prefer merging similar things over keeping parallel implementations.
- Prefer flatter structures over deeply nested ones.
- Make every remaining line earn its place.

## Procedure

1. Identify the observable behavior that must remain unchanged.
2. Find the simplest structure that preserves that behavior.
3. Remove duplication, dead paths, and incidental complexity.
4. Rename symbols when it improves clarity.
5. Merge or split units when it produces a clearer design.
6. Verify the result with the project's relevant build, test, and lint steps.

## Constraints

- Do not change observable behavior unless explicitly asked.
- Do not introduce new dependencies or new features.
- Do not skip verification: build, test, and lint after changes.