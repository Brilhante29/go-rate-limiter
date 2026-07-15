# Reuse Improvement Review

Project: `12 - go-rate-limiter`

## Review Points

- [x] after scaffold
- [x] after architecture decision
- [x] after first working slice
- [x] after benchmark result
- [x] before publication
- [ ] after CI failure, if applicable

## Findings

| Finding | Classification | Kit Area | Action | Status |
|---|---|---|---|---|
| Older scaffolds lacked manifest and complete SDD | `patch_now` | templates/tools/validation | Added conservative backfill and integrated sync switch | complete in kit `175bdb2` |
| Project contract capped IDs at 30 | `patch_now` | contracts/docs | Kept 30 as baseline and removed the growth ceiling | complete in kit `175bdb2` |
| OpenSpec article generator hardcoded a retrieval narrative | `patch_now` | openspec/tools/validation | Made narrative domain-agnostic and added regression gate | complete in kit `792773e` |
| Shared project validator does not detect Go modules | `patch_now` | templates/validation | Added gofmt, test, and vet gates with Docker-aware fallback | complete in kit `d4d30f9` |
| Planner read average latency as p95 and parsed the full manifest with a backtracking-prone regex | `patch_now` | openspec/tools | Separated p95 extraction and bounded component-pack parsing | complete in kit `0060d0b` |
| PowerShell here-strings corrupted Markdown backticks in generated artifacts | `patch_now` | openspec/tools/validation | Escaped literal Markdown and added regression gates | complete in kit `0b65b95` |
| Final validation did not require generated OpenSpec evidence | `patch_now` | templates/validation/docs | Required config and complete artifact graph before release | complete in kit `2087fec` |
| Primary benchmark metric could disagree with the README headline | `patch_now` | templates/validation/docs | Cross-checked manifest metric, committed JSON, and README opening | complete in kit `113d649` |
| Five earlier project manifests contained invalid YAML indentation | `patch_now` | templates/validation | Added top-level shape checks plus full PyYAML parsing when available; repaired all five | complete in kit `113d649` |
| Multi-node correctness benchmark may recur in backend projects | `backlog` | harness | Extract only after a second project proves the same shape | recorded; avoid premature abstraction |
| Redis Lua token bucket could be copied into the kit | `reject` | templates | Keep algorithm project-owned; the kit owns benchmark contracts, not domain behavior | rejected |

## Patch Now Decisions

- Go-aware validation, robust OpenSpec generation, Markdown-safe output, and artifact-graph, manifest, and benchmark-claim release gates were added to the kit and resynced before publication.

## Backlog Decisions

- Reassess a reusable distributed HTTP benchmark after `api-gateway-lite` or `multi-tenant-starter` needs the same invariant and node-observation logic.

## Rejected Improvements

- Do not add a generic Redis adapter, Lua library, or rate-limiter template to the kit.
- Do not add brokers, cloud emulators, or a control plane to make the shared layer appear broader.

## Final Gate

- [x] Reusable improvements were patched or recorded.
- [x] Project-specific implementation was not moved into the kit.
- [x] Validation reflects each repeated mistake discovered during this project, including Go, OpenSpec, YAML structure, and benchmark-claim gates.
