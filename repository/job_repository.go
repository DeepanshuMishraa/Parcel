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

func GetAllJobs(db *sql.DB, user_id string) ([]models.Job, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT job_id, job_name, status, user_id, payload FROM jobs WHERE user_id=$1`

	rows, err := db.QueryContext(ctx, query, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []models.Job
	for rows.Next() {
		var job models.Job
		var payloadByte []byte
		if err := rows.Scan(
			&job.JobID,
			&job.JobName,
			&job.JobStatus,
			&job.UserId,
			&payloadByte,
		); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(payloadByte, &job.Payload); err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}

	return jobs, rows.Err()
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

func UpdateJobByID(db *sql.DB, id string, status models.Status) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `UPDATE jobs SET status=$1 WHERE job_id=$2`

	_, err := db.ExecContext(ctx, query, status, id)
	return err
}
