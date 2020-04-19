package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// auth for the token authorization
// Authorization: Bearer <token>
func (s *Server) auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// try to extract the token from bear token
		tokenStr := extractToken(c)
		if tokenStr == "" {
			send(c, http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		// initialize a new instance of 'Claims'
		claims := &Claims{}
		// parse the jwt string into claims
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		// check if the token is valid or if the signature is matched
		if err != nil {
			send(c, http.StatusUnauthorized, fmt.Sprintf("Parse token: %v", err))
			c.Abort()
			return
		}
		if token.Valid == false {
			send(c, http.StatusUnauthorized, "Token is invalid")
			c.Abort()
			return
		}

		now := time.Now()
		// check if the token is expired
		if claims.ExpiresAt-now.Unix() < 30 {
			send(c, http.StatusUnauthorized, "Token is expired")
		}
	}
}
