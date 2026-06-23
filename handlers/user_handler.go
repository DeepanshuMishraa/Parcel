package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/DeepanshuMishraa/Parcel/config"
	"github.com/DeepanshuMishraa/Parcel/models"
	"github.com/DeepanshuMishraa/Parcel/repository"
	"github.com/DeepanshuMishraa/Parcel/types"
	"github.com/DeepanshuMishraa/Parcel/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

func LoginRequestHandler(db *sql.DB, cfg *config.Config) gin.HandlerFunc {
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

		claims := jwt.MapClaims{
			"user_id": user.Id,
			"email":   user.Email,
			"exp":     time.Now().Add(24 * time.Hour).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		tokenString, err := token.SignedString([]byte(cfg.JWT_SECRET))

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"token":   tokenString,
			"message": "Login SuccessFull",
		})
	}
}
