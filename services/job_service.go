package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/DeepanshuMishraa/Parcel/config"
	"github.com/DeepanshuMishraa/Parcel/models"
	"github.com/DeepanshuMishraa/Parcel/repository"
	"github.com/DeepanshuMishraa/Parcel/tools"
	"github.com/DeepanshuMishraa/Parcel/types"
	"github.com/redis/go-redis/v9"
)

type JobService struct {
	DB    *sql.DB
	Redis *redis.Client
}

func (s *JobService) CreateJobService(job models.Job) (*models.Job, error) {
	createdJob, err := repository.CreateJob(s.DB, job, job.UserId)

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
	jobID string,
	cfg *config.Config,
) (*models.Job, error) {

	job, err := repository.GetJobById(s.DB, jobID)
	if err != nil {
		return nil, err
	}

	err = repository.UpdateJobByID(s.DB, job.JobID, models.RUNNING)
	if err != nil {
		return nil, err
	}

	var runErr error
	switch job.JobName {
	case "send_email":
		var payload types.EmailPayload
		payloadBytes, marshalErr := json.Marshal(job.Payload)
		if marshalErr != nil {
			runErr = marshalErr
			break
		}

		if unmarshalErr := json.Unmarshal(payloadBytes, &payload); unmarshalErr != nil {
			runErr = unmarshalErr
			break
		}

		runErr = tools.SendEmail(cfg, payload.To, payload.Subject, payload.Body)

	case "send_pasta":
		var payload types.SendPastaPayload
		payloadBytes, marshalErr := json.Marshal(job.Payload)
		if marshalErr != nil {
			runErr = marshalErr
			break
		}

		if unmarshalErr := json.Unmarshal(payloadBytes, &payload); unmarshalErr != nil {
			runErr = unmarshalErr
			break
		}

		runErr = tools.SendPasta(payload.Who)

	default:
		runErr = fmt.Errorf("unknown job type: %s", job.JobName)
	}

	if runErr != nil {
		_ = repository.UpdateJobByID(s.DB, job.JobID, models.FAILED)
		return nil, runErr
	}

	_ = repository.UpdateJobByID(s.DB, job.JobID, models.FINISHED)
	job.JobStatus = models.FINISHED

	return job, nil
}
