package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// filter defines the fields for search
type filter struct {
	Query string `json:"query"`
	Skip  int    `json:"skip"`
	Take  int    `json:"take"`
}

// search the email from elasticsearch with query string
func (s *Server) search(c *gin.Context) {
	var u filter
	if err := c.ShouldBindJSON(&u); err != nil {
		send(c, http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	if u.Query == "" {
		send(c, http.StatusBadRequest, "Query is empty")
		return
	}

	// search the elastic with the fileter fields
	emails, err := s.elasticHandler.Search(c.Request.Context(), u.Query, u.Skip, u.Take)
	if err != nil {
		send(c, http.StatusInternalServerError, fmt.Sprintf("elastic search: %v", err))
		return
	}
	c.JSON(http.StatusOK, emails)
}
