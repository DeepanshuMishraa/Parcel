package handlers

import (
	"database/sql"
	"net/http"

	"github.com/DeepanshuMishraa/mini-job-queue/models"
	"github.com/DeepanshuMishraa/mini-job-queue/repository"
	"github.com/DeepanshuMishraa/mini-job-queue/types"
	"github.com/DeepanshuMishraa/mini-job-queue/utils"
	"github.com/gin-gonic/gin"
)

func RegisterRequestHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.RegisterRequest

		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if len(req.Password) < 6 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Password Must be greater than or equal to 6 characters",
			})
			return
		}

		hashedPass, err := utils.HashPassword(req.Password)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		user := &models.User{
			Name:     req.Name,
			Email:    req.Email,
			Password: hashedPass,
		}

		createdUser, err := repository.CreateUser(db, *user)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, types.RegisterResponse{
			Id:    createdUser.Id,
			Email: createdUser.Email,
		})
	}
}

func LoginRequestHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.LoginRequest

		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		user, error := repository.GetUserByEmail(db, req.Email)

		if error != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User Doesnt exist",
			})
			return
		}

		ok := utils.CompareHashedPassword(req.Password, user.Password)

		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invaild Password",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Login SuccessFull",
		})
	}
}
