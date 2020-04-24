package firebase

import (
	"Mail-Achive/pkg/config"
	"Mail-Achive/pkg/model"
	"Mail-Achive/pkg/utils"
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

const (
	// collectionOrderKey - the order key for the collection
	collectionOrderKey = "CreatedAt"
)

// Handler for firebase
type Handler struct {
	document string
	limit    int
	client   *firestore.Client
	ctx      context.Context
}

// NewHandler returns a handler for elastic
func NewHandler(settings *config.Config) *Handler {
	s := &Handler{
		document: settings.Server.DocumentName,
		limit:    settings.Server.LimitNumber,
		ctx:      context.Background(),
	}

	// init the firebase client
	creds := option.WithCredentialsFile(settings.Server.FirebaseCreds)
	app, err := firebase.NewApp(s.ctx, nil, creds)
	if err != nil {
		panic(err)
	}
	client, err := app.Firestore(s.ctx)
	if err != nil {
		panic(err)
	}
	s.client = client

	return s
}

// Close - release the connection
func (s *Handler) Close() {
	if s.client != nil {
		s.client.Close()
	}
}

// Set - upset the emails to firebase
func (s *Handler) Set(emails []*model.Email) error {
	// upset the email content into firebase
	for _, email := range emails {
		// md5 the content as the id
		id := utils.MD5Str(email.Content)
		if _, err := s.client.Collection(s.document).Doc(id).Set(s.ctx, email); err != nil {
			return fmt.Errorf("upset the firebase document: %v", err)
		}
	}

	return nil
}

// Flush - fetch the data after the special position
func (s *Handler) Flush(handle func(emails []*model.Email) error) error {
	// connect to the collection
	emails := s.client.Collection(s.document)

	var after *firestore.DocumentSnapshot
	for {
		var iter *firestore.DocumentIterator
		// filter the documents by order and limit
		if after == nil {
			iter = emails.OrderBy(collectionOrderKey, firestore.Asc).
				Limit(s.limit).
				Documents(s.ctx)
		} else {
			iter = emails.OrderBy(collectionOrderKey, firestore.Asc).
				StartAfter(after.Data()[collectionOrderKey]).
				Limit(s.limit).
				Documents(s.ctx)
		}

		// fetch the documents
		docs, err := iter.GetAll()
		if err != nil {
			return err
		}

		// if no more records, break the loop
		if len(docs) == 0 {
			break
		}
		after = docs[len(docs)-1]

		emails := []*model.Email{}
		// format the email structure
		for _, doc := range docs {
			emails = append(emails, &model.Email{
				From:      doc.Data()["From"].(string),
				Title:     doc.Data()["Title"].(string),
				CreatedAt: doc.Data()["CreatedAt"].(time.Time),
				Content:   doc.Data()["Content"].(string),
			})
		}

		// handle the emails
		if err := handle(emails); err != nil {
			return err
		}
	}

	return nil
}
