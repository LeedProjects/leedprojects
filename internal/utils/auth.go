package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/ronyv89/leedprojects/internal/models"
)

var jwtKey = []byte("NPV3gnmiZMEEZyujo5ZDNx7maCTaGQFTAG9egYMlAhzEdFKlHqFnebCRcFGCWj0")

// Claims defines the contents of the JWT decoded token
type Claims struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.StandardClaims
}

// UserToken defines the JWT encoded token for a user
func UserToken(user models.User) string {
	expirationTime := time.Now().Add(5 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Username: user.Username,
		Email:    user.Email,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, _ := token.SignedString(jwtKey)
	return tokenString
}
