package worker

import (
	"context"
	"database/sql"
	"log"

	"github.com/DeepanshuMishraa/mini-job-queue/services"
	"github.com/redis/go-redis/v9"
)

func RunWorker(db *sql.DB, rd *redis.Client, service *services.JobService) {
	ctx := context.Background()

	for {
		result, err := rd.BRPop(
			ctx,
			0,
			"jobs",
		).Result()

		if err != nil {
			log.Println(err)
			continue
		}

		jobId := result[1]

		log.Printf("Job picked by %s", jobId)

		_, err = service.ProcessJob(db, jobId)

		if err != nil {
			log.Println(err)
			continue
		}

	}
}
