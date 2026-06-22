package handlers

import (
	"database/sql"
	"net/http"

	"github.com/DeepanshuMishraa/mini-job-queue/models"
	"github.com/DeepanshuMishraa/mini-job-queue/repository"
	"github.com/DeepanshuMishraa/mini-job-queue/services"
	"github.com/DeepanshuMishraa/mini-job-queue/types"
	"github.com/gin-gonic/gin"
)

func CreateJobHandler(jobService *services.JobService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.CreateJobRequest
		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		status := models.QUEUED
		if req.Status != "" {
			status = models.Status(req.Status)
		}

		job := models.Job{
			JobName:   req.JobName,
			JobStatus: status,
			Payload:   req.Payload,
			UserId:    req.UserId,
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
