package repository

import (
	"context"
	"database/sql"
	"github.com/DeepanshuMishraa/mini-job-queue/models"
	"time"
)

func CreateJob(db *sql.DB, job models.Job) (*models.Job, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `INSERT INTO jobs(job_name,payload) VALUES($1,$2) RETURNING job_id, job_name, status`

	createdJob := &models.Job{}
	err := db.QueryRowContext(ctx, query, job.JobName, job.Payload).Scan(
		&createdJob.JobID,
		&createdJob.JobName,
		&createdJob.JobStatus,
	)

	if err != nil {
		return &models.Job{}, err
	}

	return createdJob, nil
}
