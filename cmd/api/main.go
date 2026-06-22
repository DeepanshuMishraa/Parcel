package main

import (
	"log"

	"github.com/DeepanshuMishraa/mini-job-queue/config"
	"github.com/DeepanshuMishraa/mini-job-queue/db"
	"github.com/DeepanshuMishraa/mini-job-queue/handlers"
	"github.com/DeepanshuMishraa/mini-job-queue/services"
	"github.com/DeepanshuMishraa/mini-job-queue/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	router := gin.Default()
	router.SetTrustedProxies(nil)

	if err != nil {
		log.Fatal("Failed to load env vars")
	}

	dbx, err := db.ConnectDB(cfg.DATABASE_URL)

	if err != nil {
		log.Fatal("Failed to connect to the database with error: ", err)
	}

	redis, err := utils.Connect(cfg.REDIS_URL)

	if err != nil {
		log.Fatal("Failed to connect to redis with error: ", err)
	}

	jobService := &services.JobService{
		DB:    dbx,
		Redis: redis,
	}

	go RunWorker(dbx, redis, jobService)

	router.POST("/api/user/register", handlers.RegisterRequestHandler(dbx))
	router.POST("/api/user/login", handlers.LoginRequestHandler(dbx, cfg))
	router.POST("/api/jobs/create", handlers.CreateJobHandler(jobService))
	router.GET("/api/job/:id", handlers.GetJobByIdHandler(dbx))
	router.GET("/api/jobs/:id", handlers.GetAllJobHandler(dbx))

	router.GET("/api/health", gin.HandlerFunc(func(c *gin.Context) {
		c.JSON(201, gin.H{
			"message": "OK",
		})
	}))

	log.Println("[API] SERVER RUNNING ON PORT: ", cfg.PORT)
	router.Run(":" + cfg.PORT)
}
