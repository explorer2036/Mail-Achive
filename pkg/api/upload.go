package api

import (
	"Mail-Achive/pkg/log"
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

		start := time.Now()
		// unzip and read the email files
		err := s.unzip(path)
		if err != nil {
			send(c, http.StatusInternalServerError, fmt.Sprintf("Unzip: %v", err))
		} else {
			send(c, http.StatusOK, fmt.Sprintf("Upload %s success", file.Filename))
		}
		log.Infof("Handle upload file %s, %d, %v", file.Filename, file.Size, time.Since(start))
	}
}
