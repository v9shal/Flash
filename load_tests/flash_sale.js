import http from 'k6/http';
import { check } from 'k6';
import { Counter } from 'k6/metrics';
import crypto from 'k6/crypto';
import encoding from 'k6/encoding';

const successfulBookings = new Counter('successful_bookings');
const conflictBookings = new Counter('conflict_bookings');

export const options = {
  scenarios: {
    flash_sale_spike: {
      executor: 'per-vu-iterations',
      vus: 100,              // 100 distinct users
      iterations: 1,         // 1 attempt each
      maxDuration: '10s',
    },
  },
};

const JWT_SECRET = "supersecretjwtkey12345"; // Match your JWT_SECRET from .env

// Helper to generate a valid JWT for User ID = __VU
function generateToken(userId) {
  const header = encoding.b64encode(JSON.stringify({ alg: "HS256", typ: "JWT" }), "rawurl");
  const payload = encoding.b64encode(JSON.stringify({
    user_id: userId,
    role: "user",
    exp: Math.floor(Date.now() / 1000) + 3600
  }), "rawurl");

  // Changed "rawurl" -> "base64rawurl"
  const signature = crypto.hmac("sha256", JWT_SECRET, `${header}.${payload}`, "base64rawurl");
  return `${header}.${payload}.${signature}`;
}

export default function () {
  // Each VU represents a unique user (User 1 to 100)
  const userId = __VU;
  const token = generateToken(userId);

  // Pick a seat randomly between 1 and 10
  const seatId = Math.floor(Math.random() * 10) + 1;
  const seatNumber = seatId <= 5 ? `A${seatId}` : `B${seatId - 5}`;
  
 // Generates a valid UUID for User 1 to 100: e.g. "00000000-0000-4000-8000-000000000007"
const userPadded = String(userId).padStart(12, '0');
const idempotencyKey = `00000000-0000-4000-8000-${userPadded}`;

  const payload = JSON.stringify({
    event_id: 1,
    seat_id: seatId,
    seat_number: seatNumber,
    expected_version: 1,
    price: "1500.00",
    idempotency_key: idempotencyKey,
  });

  const params = {
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`,
    },
  };

  const res = http.post('http://localhost:8080/bookings', payload, params);

  if (res.status === 201) {
    successfulBookings.add(1);
  } else if (res.status === 409) {
    conflictBookings.add(1);
  }

  check(res, {
    'status is 201 or 409': (r) => r.status === 201 || r.status === 409,
  });
}