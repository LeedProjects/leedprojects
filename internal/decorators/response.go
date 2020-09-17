package decorators

import (
	"time"

	"github.com/ronyv89/leedprojects/internal/models"
	"github.com/ronyv89/leedprojects/internal/utils"
)

// FormattedUser formats user model instance for API responses
type FormattedUser struct {
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserResponse formats a user for JSON response
func UserResponse(user models.User) FormattedUser {
	return FormattedUser{
		Username:  user.Username,
		Email:     user.Email,
		Token:     utils.UserToken(user),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
