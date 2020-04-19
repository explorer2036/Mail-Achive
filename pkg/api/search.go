package api

import (
	"Mail-Achive/pkg/model"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// filter defines the fields for search
type filter struct {
	Query string `json:"query"`
	Skip  int    `json:"skip"`
	Take  int    `json:"take"`
}

// EmailResponse for the frontend
type EmailResponse struct {
	Name      string `json:"name"`
	From      string `json:"from"`
	CreatedAt string `json:"created_at"`
	Title     string `json:"title"`
	Content   string `json:"content"`
}

// format the reponse emails
func format(emails []*model.Email) []EmailResponse {
	result := []EmailResponse{}
	for _, email := range emails {
		result = append(result, EmailResponse{
			Name:      email.Name,
			From:      fromPrefix + " " + email.From,
			Title:     titlePrefix + " " + email.Title,
			CreatedAt: timePrefix + " " + email.CreatedAt.Format(time.RFC1123),
			Content:   email.Content,
		})
	}
	return result
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
		send(c, http.StatusInternalServerError, fmt.Sprintf("Elastic search: %v", err))
		return
	}
	c.JSON(http.StatusOK, format(emails))
}
