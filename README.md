# Mini Job Queue

Lightweight HTTP job queue built with Go. Users submit jobs with arbitrary JSON payloads, and a background worker picks them up, processes them, and tracks status through the lifecycle.

It's called "mini" because it's not a full production job queue yet — next iteration will add RabbitMQ, retries, dead-letter queues, and proper job execution.

## Arch Diagram

![Arch Diagram](./arch.png)

## Tech Stack

- **Go** — API server + worker
- **Gin** — HTTP router
- **PostgreSQL** — job/user persistence
- **Redis** — queue (LPUSH / BRPOP)
- **JWT** — auth tokens

## How It Works

1. User registers and logs in (gets a JWT).
2. User creates a job with a `job_name` and JSON `payload`. `user_id` is extracted from the JWT session.
3. Job is saved to PostgreSQL and pushed onto a Redis list.
4. A worker goroutine blocks on BRPOP, picks up the job ID, marks it `running`, processes it, then marks it `finished` or `failed`.

## Job Types & Payloads

### `send_email`

Sends an email via Resend(configured via env vars).

```json
{
  "job_name": "send_email",
  "payload": {
    "to": "user@example.com",
    "subject": "Hello",
    "body": "This is the email body."
  }
}
```

### `send_pasta`

Sends pasta emoji via Telegram.

```json
{
  "job_name": "send_pasta",
  "payload": {
    "who": "@username"
  }
}
```

## Schema

**users** — `id` (UUID), `name`, `email` (unique), `password` (bcrypt), timestamps

**jobs** — `job_id` (UUID), `job_name`, `status` (queued / running / finished / failed), `user_id` (FK), `payload` (JSONB), timestamps

## Endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | `/api/user/register` | No | Create user (name, email, password) |
| POST | `/api/user/login` | No | Login, returns JWT |
| POST | `/api/jobs/create` | Yes | Create job (`job_name` + `payload`) |
| GET | `/api/job/:id` | Yes | Get single job by job_id |
| GET | `/api/jobs` | Yes | Get all jobs for authenticated user |
| GET | `/api/health` | No | Health check |

## cURL Examples

```bash
# Register
curl -X POST localhost:8080/api/user/register \
  -H "Content-Type: application/json" \
  -d '{"name": "test", "email": "test@test.com", "password": "password123"}'

# Login (save the token)
curl -X POST localhost:8080/api/user/login \
  -H "Content-Type: application/json" \
  -d '{"email": "test@test.com", "password": "password123"}'

# Create a send_email job
curl -X POST localhost:8080/api/jobs/create \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <TOKEN>" \
  -d '{"job_name": "send_email", "payload": {"to": "user@example.com", "subject": "Hello", "body": "Email body here"}}'

# Create a send_pasta job
curl -X POST localhost:8080/api/jobs/create \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <TOKEN>" \
  -d '{"job_name": "send_pasta", "payload": {"who": "@username"}}'

# Get all jobs for the authenticated user
curl -X GET localhost:8080/api/jobs \
  -H "Authorization: Bearer <TOKEN>"

# Get a specific job
curl -X GET localhost:8080/api/job/<JOB_ID> \
  -H "Authorization: Bearer <TOKEN>"
```

## Setup

```bash
# environment
cp .env.example .env
# fill in DATABASE_URL, REDIS_URL, JWT_SECRET, PORT

# run migrations (use your migration tool of choice)
# then start the server
go run cmd/api/main.go
```


