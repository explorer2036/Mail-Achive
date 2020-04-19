package api

import (
	"Mail-Achive/pkg/model"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var (
	// jwt secret for generating the jwt token
	jwtSecret = []byte("pngOEtu35ewyKJWZ")
	// token expire time (unit: minute)
	tokenExpire = 120
	// http header: Authorization
	headerAuthorization = "Authorization"
)

// Claims - Create a struct that will be encoded to a JWT
type Claims struct {
	Key string `json:"key"`
	jwt.StandardClaims
}

func extractToken(c *gin.Context) string {
	bearToken := c.Request.Header.Get(headerAuthorization)
	parts := strings.Split(bearToken, " ")
	if len(parts) == 2 {
		return parts[1]
	}
	return ""
}

// login with username and password
func (s *Server) login(c *gin.Context) {
	var u model.User
	if err := c.ShouldBindJSON(&u); err != nil {
		send(c, http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}

	// validate the username and password
	user, ok := s.settings.Users[u.Username]
	if !ok {
		send(c, http.StatusUnauthorized, "Username is not existed")
		return
	}
	if user.Username != u.Username || user.Password != u.Password {
		send(c, http.StatusUnauthorized, "Username or password is wrong")
		return
	}

	// declare the expiration time of the token
	expiration := time.Now().Add(time.Duration(tokenExpire) * time.Minute)
	// create the jwt claims, which include the username and expiration time
	claims := &Claims{
		Key: u.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiration.Unix(),
		},
	}

	// decleare the token with the algorithem used for signing
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtSecret)
	if err != nil {
		send(c, http.StatusInternalServerError, fmt.Sprintf("Generate token: %v", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":   token,
		"manager": user.Manager,
	})
}

// refresh the token
func (s *Server) refresh(c *gin.Context) {
	tokenStr := extractToken(c)

	claims := &Claims{}
	// parse the jwt string into claims
	if _, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	}); err != nil {
		send(c, http.StatusInternalServerError, fmt.Sprintf("Parse token: %v", err))
		return
	}

	// declare the expiration time of the token
	expiration := time.Now().Add(time.Duration(tokenExpire) * time.Minute)
	// create the jwt claims, which include the username and expiration time
	claims.StandardClaims = jwt.StandardClaims{ExpiresAt: expiration.Unix()}

	// decleare the token with the algorithem used for signing
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtSecret)
	if err != nil {
		send(c, http.StatusInternalServerError, fmt.Sprintf("Generate token: %v", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
