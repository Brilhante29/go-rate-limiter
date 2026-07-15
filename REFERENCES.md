# References

No source code was copied from the references below. The implementation, fixtures, Lua script, tests, and benchmark are project-specific.

| Reference | License | Used for | Copied code? |
|---|---|---|---|
| [Go 1.26 release history](https://go.dev/doc/devel/release) | BSD-3-Clause | Supported runtime and security patch selection | no |
| [go-chi/chi](https://github.com/go-chi/chi) | MIT | `net/http`-compatible routing and middleware | no |
| [redis/go-redis](https://github.com/redis/go-redis) | BSD-2-Clause | Official Redis client and script execution | no |
| [Redis rate limiter guide](https://redis.io/docs/latest/develop/use-cases/rate-limiter/) | Documentation terms | Token-bucket and atomicity guidance | no |
| [Redis Lua scripting](https://redis.io/docs/latest/develop/programmability/eval-intro/) | Documentation terms | Atomic server-side execution guarantee | no |
| [Grafana k6 thresholds](https://grafana.com/docs/k6/latest/using-k6/thresholds/) | AGPL-3.0 for k6 | Independent load profile and pass/fail thresholds | no |
| [envoyproxy/ratelimit](https://github.com/envoyproxy/ratelimit) | Apache-2.0 | Distributed rate-limit service organization reference | no |
| [uber-go/ratelimit](https://github.com/uber-go/ratelimit) | MIT | Go limiter API and benchmark organization reference | no |

## Reuse Boundary

- Reused: official libraries under their licenses, organization ideas, contract shape, and benchmark conventions.
- Project-owned: domain port, Redis key model, Lua token-bucket script, HTTP response, Docker demo, correctness invariant, and result JSON.
