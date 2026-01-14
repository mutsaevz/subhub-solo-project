import http from 'k6/http';
import { check } from 'k6';

export const options = {
  vus: 20,
  duration: '10s',
  thresholds: {
    http_req_failed: ['rate<0.01'],
http_req_duration: ['p(95)<1500']
  },
};

export default function () {
  const res = http.get('http://localhost:8080/api/payments?limit=100');

  check(res, {
    'GET /payments?limit=100 status 200': (r) => r.status === 200,
  });
}