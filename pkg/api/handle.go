package api

import (
	"Mail-Achive/pkg/model"
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/google/martian/log"
	"golang.org/x/text/encoding/charmap"
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

// Gesendet: Dienstag, 28. Januar 2014 um 13:26 Uhr
// Gesendet: Freitag, 31. Januar 2014 13:59
// Gesendet: Freitag, 1. Januar 2014 13:59
func timeStr2Time(ts string) (*time.Time, error) {
	parts := strings.Split(ts, " ")
	if len(parts) != 4 && len(parts) != 6 {
		return nil, fmt.Errorf("invalid time format: %s(%d)", ts, len(parts))
	}
	if len(parts) == 6 {
		parts[3] = parts[4]
		parts = parts[0:4]
	}
	day, err := strconv.Atoi(strings.TrimSuffix(parts[0], "."))
	if err != nil {
		return nil, err
	}
	parts[0] = fmt.Sprintf("%02d", day)

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
		// check if it's a directory
		if file.FileInfo().IsDir() {
			continue
		}

		// open the file
		fd, err := file.Open()
		if err != nil {
			return err
		}
		defer fd.Close()

		// new a file reader with europe standard
		reader := charmap.ISO8859_1.NewDecoder().Reader(fd)
		// new a buffer reader
		br := bufio.NewReader(reader)

		// parse the email content to structure
		email, err := s.parse(br)
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
func (s *Server) parse(br *bufio.Reader) (*model.Email, error) {
	e := &model.Email{}

	// loop the buffer reader, and read line one by one
	for {
		data, _, err := br.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("parse buffer: %v", err)
		}
		line := strings.TrimSpace(string(data))

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
			left := make([]byte, br.Buffered())
			if _, err := br.Read(left); err != nil {
				log.Errorf("buffer reader the left data: %v", err)
				return nil, err
			}
			e.Content = string(left)
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
