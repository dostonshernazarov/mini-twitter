import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  vus: 50,
  duration: '1m',
};

const BASE_URL = 'http://localhost:7777/v1';

export default function () {
  const payload = JSON.stringify({
    user_id: 'test-user-id',
    content: 'This is a load test tweet',
  });

  const params = {
    headers: { 'Content-Type': 'application/json' },
  };

  const res = http.post(`${BASE_URL}/tweets`, payload, params);
  check(res, { 'status was 201': (r) => r.status === 201 });
  sleep(1);
}
