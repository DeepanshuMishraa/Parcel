package services

import (
	"context"
	"database/sql"
	"math/rand"
	"time"

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

func (s *JobService) ProcessJob(
	db *sql.DB,
	jobID string,
) (*models.Job, error) {

	job, err := repository.GetJobById(
		db,
		jobID,
	)

	if err != nil {
		return nil, err
	}

	err = repository.UpdateJobByID(
		db,
		job.JobID,
		models.RUNNING,
	)

	if err != nil {
		return nil, err
	}

	processingTime :=
		time.Duration(rand.Intn(61)) *
			time.Second

	time.Sleep(processingTime)

	err = repository.UpdateJobByID(
		db,
		job.JobID,
		models.FINISHED,
	)

	if err != nil {
		return nil, err
	}

	job.JobStatus = models.FINISHED

	return job, nil
}
