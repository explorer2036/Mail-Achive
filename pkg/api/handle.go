package api

import (
	"Mail-Achive/pkg/model"
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/google/martian/log"
)

const (
	fromPrefix  = "Von:"
	timePrefix  = "Gesendet:"
	titlePrefix = "Betreff:"
)

// germanyMonths table
var germanyMonths = map[string]string{
	"Januar":    "Jan",
	"Februar":   "Feb",
	"März":      "Mar",
	"April":     "Apr",
	"Mai":       "May",
	"Juni":      "Jun",
	"Juli":      "Jul",
	"August":    "Aug",
	"September": "Sep",
	"Oktober":   "Oct",
	"November":  "Nov",
	"Dezember":  "Dec",
}

func timeStr2Time(ts string) (*time.Time, error) {
	parts := strings.Split(ts, " ")
	if len(parts) != 4 {
		return nil, fmt.Errorf("invalid time format: %s(%d)", ts, len(parts))
	}
	parts[0] = strings.TrimSuffix(parts[0], ".")
	parts[1], _ = germanyMonths[strings.TrimSpace(parts[1])]
	joined := strings.Join(parts, " ")
	parsed, err := time.Parse("02 Jan 2006 15:04", joined)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

// unzip and read the files in mails.zip
func (s *Server) unzip(path string) error {
	reader, err := zip.OpenReader(path)
	if err != nil {
		return err
	}
	defer reader.Close()

	emails := []*model.Email{}
	// loop the email files in zip files
	for _, file := range reader.File {
		if file.FileInfo().IsDir() {
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		size := int64(file.UncompressedSize64)
		data := make([]byte, size)
		buffer := bytes.NewBuffer(data)
		if _, err := io.CopyN(buffer, fileReader, size); err != nil {
			return err
		}

		// parse the email content to structure
		email, err := s.parse(buffer)
		if err != nil {
			return err
		}
		if email != nil {
			emails = append(emails, email)
		}
	}

	// handle the emails
	if err := s.handle(emails); err != nil {
		return err
	}

	return nil
}

// parse the email content to normal fields
func (s *Server) parse(buffer *bytes.Buffer) (*model.Email, error) {
	e := &model.Email{}

	// loop the buffer reader, and read line one by one
	for {
		line, err := buffer.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("parse buffer: %v", err)
		}
		runeLine := []rune(line)

		line = strings.TrimSpace(string(runeLine))
		// Von: Claudia Fehrenberg [mailto:claudia.fehrenberg@gmx.de]
		if strings.HasPrefix(line, fromPrefix) {
			left := line[len(fromPrefix):]
			e.From = strings.TrimSpace(left)
		}

		// Gesendet: Freitag, 31. Januar 2014 13:59
		if strings.HasPrefix(line, timePrefix) {
			// get the time of email
			parts := strings.Split(line[len(timePrefix):], ",")
			format, err := timeStr2Time(strings.TrimSpace(parts[1]))
			if err != nil {
				log.Errorf("invalid time format: %v", line)
				return nil, nil
			}
			e.CreatedAt = *format
		}

		// Betreff: KT: Atopie-Hund mit Dauerfieber
		if strings.HasPrefix(line, titlePrefix) {
			// get the title of email
			e.Title = strings.TrimSpace(line[len(titlePrefix):])

			// read the left email body as content
			// the email file content, use the []rune to solve the invalid UTF-8 characters
			content := []rune(strings.TrimSpace(buffer.String()))
			e.Content = string(content)
		}
	}
	return e, nil
}

// handle the emails
func (s *Server) handle(emails []*model.Email) error {
	if len(emails) == 0 {
		return nil
	}

	// upset the email content into firebase
	if err := s.firebaseHandler.Set(emails); err != nil {
		return err
	}

	// upset the email content into elastic
	return s.elasticHandler.Bulk(emails)
}
