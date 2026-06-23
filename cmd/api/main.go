package main

import (
	"log"

	"github.com/DeepanshuMishraa/Parcel/config"
	"github.com/DeepanshuMishraa/Parcel/db"
	"github.com/DeepanshuMishraa/Parcel/handlers"
	"github.com/DeepanshuMishraa/Parcel/middleware"
	"github.com/DeepanshuMishraa/Parcel/services"
	"github.com/DeepanshuMishraa/Parcel/utils"
	"github.com/DeepanshuMishraa/Parcel/worker"
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

	go worker.RunWorker(cfg, redis, jobService)

	protectedRouter := router.Group("/api")
	protectedRouter.Use(middleware.AuthMiddleware(cfg))

	router.POST("/api/user/register", handlers.RegisterRequestHandler(dbx))
	router.POST("/api/user/login", handlers.LoginRequestHandler(dbx, cfg))
	protectedRouter.POST("/jobs/create", handlers.CreateJobHandler(jobService))
	protectedRouter.GET("/job/:id", handlers.GetJobByIdHandler(dbx))
	protectedRouter.GET("/jobs", handlers.GetAllJobHandler(dbx))

	router.GET("/api/health", gin.HandlerFunc(func(c *gin.Context) {
		c.JSON(201, gin.H{
			"message": "OK",
		})
	}))

	log.Println("[API] SERVER RUNNING ON PORT: ", cfg.PORT)
	router.Run(":" + cfg.PORT)
}
