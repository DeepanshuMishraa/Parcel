package services

import (
	"context"
	"database/sql"

	"github.com/DeepanshuMishraa/mini-job-queue/models"
	"github.com/DeepanshuMishraa/mini-job-queue/repository"
	"github.com/redis/go-redis/v9"
)

type JobService struct {
	DB    *sql.DB
	Redis *redis.Client
}

func (s *JobService) CreateJobService(job models.Job) (*models.Job, error) {
	createdJob, err := repository.CreateJob(s.DB, job)

	if err != nil {
		return nil, err
	}

	err = s.Redis.LPush(
		context.Background(),
		"jobs",
		createdJob.JobID,
	).Err()

	if err != nil {
		return nil, err
	}

	return createdJob, nil

}
