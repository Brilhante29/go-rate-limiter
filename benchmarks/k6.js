import http from 'k6/http';
import { check, sleep } from 'k6';
import { Counter } from 'k6/metrics';

const targets = (__ENV.TARGETS || 'http://node-a:8080,http://node-b:8080').split(',');
const accepted = new Counter('rate_limit_accepted');
const rejected = new Counter('rate_limit_rejected');

http.setResponseCallback(http.expectedStatuses(200, 429));

export const options = {
  vus: 32,
  duration: '5s',
  thresholds: {
    checks: ['rate>0.99'],
    http_req_duration: ['p(95)<250'],
    http_req_failed: ['rate<0.01'],
  },
};

export function setup() {
  for (let attempt = 0; attempt < 50; attempt += 1) {
    const responses = targets.map((target) => http.get(`${target}/healthz`));
    if (responses.every((response) => response.status === 200)) {
      return;
    }
    sleep(0.1);
  }
  throw new Error('rate-limiter nodes did not become ready');
}

export default function () {
  const target = targets[__ITER % targets.length];
  const response = http.post(
    `${target}/v1/limits/k6-shared/check`,
    JSON.stringify({ cost: 1 }),
    { headers: { 'Content-Type': 'application/json' } },
  );

  check(response, {
    'decision is valid': (result) => result.status === 200 || result.status === 429,
    'node is identified': (result) => Boolean(result.headers['X-Rate-Limiter-Node']),
  });
  if (response.status === 200) {
    accepted.add(1);
  } else if (response.status === 429) {
    rejected.add(1);
  }
}
