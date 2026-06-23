package worker

import (
	"context"
	"log"

	"github.com/DeepanshuMishraa/Parcel/config"
	"github.com/DeepanshuMishraa/Parcel/services"
	"github.com/redis/go-redis/v9"
)

func RunWorker(cfg *config.Config, rd *redis.Client, service *services.JobService) {
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

		_, err = service.ProcessJob(jobId, cfg)

		if err != nil {
			log.Println(err)
			continue
		}

	}
}
