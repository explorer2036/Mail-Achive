package api

import (
	"Mail-Achive/pkg/model"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// flush prevent the situation elastic is broken or upgraded
// it needs to sync the firebase data to elasticsearch
func (s *Server) flush(c *gin.Context) {
	// check if the elastic index is existed
	if err := s.elasticHandler.CreateIndex(); err != nil {
		send(c, http.StatusInternalServerError, fmt.Sprintf("Elastic create index: %v", err))
		return
	}

	// flush the all firebase data to elastic
	if err := s.firebaseHandler.Flush(func(emails []*model.Email) error {
		return s.elasticHandler.Bulk(emails)
	}); err != nil {
		send(c, http.StatusInternalServerError, fmt.Sprintf("Flush firbase data: %v", err))
		return
	}

	send(c, http.StatusOK, "Flush firebase to elastic success")
}
