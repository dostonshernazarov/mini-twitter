import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  vus: 100,
  duration: '30s',
};

const BASE_URL = 'http://localhost:7777/v1';

export default function () {
  const res = http.get(`${BASE_URL}/tweets`);
  check(res, { 'status was 200': (r) => r.status === 200 });
  sleep(1);
}
