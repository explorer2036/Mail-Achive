package api

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// upload the email zip files
func (s *Server) upload(c *gin.Context) {
	// fetch the file uploaded
	file, _ := c.FormFile("file")

	if strings.HasSuffix(file.Filename, ".zip") == false {
		send(c, http.StatusInternalServerError, fmt.Sprintf("%v is not supported", file.Filename))
		return
	}

	// temporary files
	path := fmt.Sprintf("%d_%s", time.Now().UnixNano()/1000000, file.Filename)
	if err := c.SaveUploadedFile(file, path); err != nil {
		send(c, http.StatusInternalServerError, fmt.Sprintf("Save uploaded file: %v", err))
	} else {
		// remove the temporary files
		defer os.RemoveAll(path)

		// unzip and read the email files
		if err := s.unzip(path); err != nil {
			send(c, http.StatusInternalServerError, fmt.Sprintf("Unzip: %v", err))
		} else {
			send(c, http.StatusOK, fmt.Sprintf("Upload %s success", file.Filename))
		}
	}
}
