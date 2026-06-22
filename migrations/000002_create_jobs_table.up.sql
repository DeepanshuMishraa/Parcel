CREATE TYPE Status AS ENUM('queued',
'running',
'finished',
'failed');
CREATE TABLE IF NOT EXISTS jobs(
  job_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  job_name TEXT NOT NULL,
  status Status NOT NULL DEFAULT 'queued',
  user_id UUID NOT NULL REFERENCES users(id),
  payload JSONB NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  started_at TIMESTAMPTZ,
  completed_at TIMESTAMPTZ
);
