package handlers

import (
	"database/sql"
	"net/http"

	"github.com/DeepanshuMishraa/Parcel/models"
	"github.com/DeepanshuMishraa/Parcel/repository"
	"github.com/DeepanshuMishraa/Parcel/services"
	"github.com/DeepanshuMishraa/Parcel/types"
	"github.com/gin-gonic/gin"
)

func CreateJobHandler(jobService *services.JobService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.CreateJobRequest
		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userId, _ := c.Get("user_id")

		status := models.QUEUED
		if req.Status != "" {
			status = models.Status(req.Status)
		}

		job := models.Job{
			JobName:   req.JobName,
			JobStatus: status,
			Payload:   req.Payload,
			UserId:    userId.(string),
		}

		createdJob, err := jobService.CreateJobService(job)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, types.CreateJobResponse{
			JobID:   createdJob.JobID,
			JobName: createdJob.JobName,
			Status:  string(createdJob.JobStatus),
		})
	}
}

func GetJobHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Param("id")
		job, err := repository.GetJobById(db, user_id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, job)
	}
}

func GetAllJobHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, _ := c.Get("user_id")
		jobs, err := repository.GetAllJobs(db, userId.(string))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, jobs)
	}
}

func GetJobByIdHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		job, err := repository.GetJobById(db, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, job)
	}
}
