package firebase

import (
	"Mail-Achive/pkg/config"
	"Mail-Achive/pkg/model"
	"Mail-Achive/pkg/utils"
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

// Handler for firebase
type Handler struct {
	document string
	client   *firestore.Client
	ctx      context.Context
}

// NewHandler returns a handler for elastic
func NewHandler(settings *config.Config) *Handler {
	s := &Handler{
		document: settings.Server.DocumentName,
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
