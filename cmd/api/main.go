package main

import (
	"log"

	"github.com/DeepanshuMishraa/mini-job-queue/config"
	"github.com/DeepanshuMishraa/mini-job-queue/db"
	"github.com/DeepanshuMishraa/mini-job-queue/handlers"
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

	_, err = utils.Connect(cfg.REDIS_URL)

	if err != nil {
		log.Fatal("Failed to connect to redis with error: ", err)
	}

	router.POST("/api/user/register", handlers.RegisterRequestHandler(dbx))
	router.POST("/api/user/login", handlers.LoginRequestHandler(dbx))

	log.Println("[API] SERVER RUNNING ON PORT: ", cfg.PORT)
	router.Run(":" + cfg.PORT)
}
