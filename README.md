# Mini Job Queue

Lightweight HTTP job queue built with Go. Users submit jobs with arbitrary JSON payloads, and a background worker picks them up, processes them (simulated for now), and tracks status through the lifecycle.

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
2. User creates a job with a `job_name`, `user_id`, and any JSON `payload`.
3. Job is saved to PostgreSQL and pushed onto a Redis list.
4. A worker goroutine blocks on BRPOP, picks up the job ID, marks it `running`, sleeps for a random 0–60s, then marks it `finished`.

## Schema

**users** — `id` (UUID), `name`, `email` (unique), `password` (bcrypt), timestamps

**jobs** — `job_id` (UUID), `job_name`, `status` (queued / running / finished / failed), `user_id` (FK), `payload` (JSONB), timestamps

## Endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/user/register` | Create user (name, email, password) |
| POST | `/api/user/login` | Login, returns JWT |
| POST | `/api/jobs/create` | Create job (auth required) |
| GET | `/api/job/:id` | Get single job by job_id |
| GET | `/api/jobs/:id` | Get all jobs for a user_id |
| GET | `/api/health` | Health check |

## Setup

```bash
# environment
cp .env.example .env
# fill in DATABASE_URL, REDIS_URL, JWT_SECRET, PORT

# run migrations (use your migration tool of choice)
# then start the server
go run cmd/api/main.go
```


