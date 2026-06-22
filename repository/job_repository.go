package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/DeepanshuMishraa/mini-job-queue/models"
)

func CreateJob(db *sql.DB, job models.Job) (*models.Job, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `INSERT INTO jobs(job_name, payload, user_id) VALUES($1, $2, $3) RETURNING job_id, job_name, status, user_id`

	payloadByte, err := json.Marshal(job.Payload)

	if err != nil {
		return nil, err
	}
	createdJob := &models.Job{}
	err = db.QueryRowContext(ctx, query, job.JobName, payloadByte, job.UserId).Scan(
		&createdJob.JobID,
		&createdJob.JobName,
		&createdJob.JobStatus,
		&createdJob.UserId,
	)

	if err != nil {
		return &models.Job{}, err
	}

	return createdJob, nil
}

func GetJobById(db *sql.DB, id string) (*models.Job, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT job_id, job_name, status, user_id, payload FROM jobs WHERE job_id=$1`

	jobs := &models.Job{}
	var payloadByte []byte
	err := db.QueryRowContext(ctx, query, id).Scan(
		&jobs.JobID,
		&jobs.JobName,
		&jobs.JobStatus,
		&jobs.UserId,
		&payloadByte,
	)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(payloadByte, &jobs.Payload); err != nil {
		return nil, err
	}

	return jobs, nil
}
