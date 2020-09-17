package controllers

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ronyv89/leedprojects/internal/connections"
	"github.com/ronyv89/leedprojects/internal/decorators"
	"github.com/ronyv89/leedprojects/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SignupParams ...
type SignupParams struct {
	Username string `binding:"required"`
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

// HashPassword gives a hashed password from the plain text password provided
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash checks if the plain text password provided and the salted password are the same
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func hashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
	}
	stringHash := string(hash)
	return stringHash
}

// AuthLogin logs in user
func AuthLogin(c *gin.Context) {
	var authParams, user models.User
	c.BindJSON(&authParams)
	result := connections.DB().Where(models.User{Username: authParams.Username}).Or(models.User{Email: authParams.Username}).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(401, gin.H{
				"error": "User not found",
			})
		} else {
			c.JSON(500, gin.H{
				"error": "Internal server error",
			})
		}
	} else if user.ID > 0 {
		match := CheckPasswordHash(authParams.Password, user.Password)
		if match {
			c.JSON(200, gin.H{
				"user": decorators.UserResponse(user),
			})
		} else {
			c.JSON(401, gin.H{
				"error": "Invalid Password",
			})
		}
	} else {
		c.JSON(500, gin.H{
			"error": "Internal server error",
		})
	}

}

// AuthSignup signs up a new user
func AuthSignup(c *gin.Context) {
	var user models.User
	var existingUsers []models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := connections.DB()
	result := db.Where(models.User{Username: user.Username}).Or(models.User{Email: user.Email}).Find(&existingUsers)
	if result.Error != nil {
		c.JSON(500, gin.H{
			"error": "Internal server error",
		})
	} else if result.RowsAffected == 0 {
		hash, _ := HashPassword(user.Password)
		user.Password = hash
		result = db.Create(&user)
		c.JSON(200, gin.H{
			"user": decorators.UserResponse(user),
		})
	} else {
		c.JSON(400, gin.H{
			"error": "username/email already in use",
		})
	}
}
