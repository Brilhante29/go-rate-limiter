---
name: reuse-improvement-review
description: Review whether work on a portfolio project should improve portfolio-reuse-kit, classify improvements as patch_now/backlog/reject, and update the kit proactively when safe.
---

# Reuse Improvement Review

Use this skill during every portfolio project, especially after scaffold, architecture decision, first working slice, benchmark result, CI failure, and before publication.

Read `decision-brain/reuse-improvement-loop.yaml` in the kit, or `.portfolio/decision-brain/reuse-improvement-loop.yaml` in a generated project.

## Review Questions

Ask whether the current project revealed:

- a repeated decision that belongs in the kit
- a missing template, schema, script, Docker pattern, benchmark helper, or skill
- ambiguous agent instructions
- a reusable local-first adapter or parity-test pattern
- a missing benchmark metric or result shape
- a validation gap

## Classification

- `patch_now`: low-risk, reusable, and improves templates, schemas, validation, skills, docs, harness, metrics, or decision matrices.
- `backlog`: useful but broad, risky, or not fully understood yet.
- `reject`: project-specific, premature, duplicated, or harmful.

## Required Output

Update `sdd/reuse-improvement-review.md`, remove template placeholders, replace the blank finding row, and mark all final gate checks that apply before publication.

If an item is `patch_now`, update `portfolio-reuse-kit`, run validation, commit, and push before treating the dependent project as portfolio-ready.

Do not upstream project-specific business code.
