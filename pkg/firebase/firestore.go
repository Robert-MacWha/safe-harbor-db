package firebase

import (
	"context"
	"fmt"
	"os"
	"strings"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func NewFirestoreClient() (*firestore.Client, error) {
	creds := os.Getenv("FIREBASE_CREDENTIALS")
	if creds == "" {
		return nil, fmt.Errorf("missing FIREBASE_CREDENTIALS env")
	}

	ctx := context.Background()
	creds = strings.Trim(creds, "'")
	sa := option.WithCredentialsJSON([]byte(creds))
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		return nil, fmt.Errorf("firebase.NewApp: %w", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, fmt.Errorf("app.Firestore: %w", err)
	}

	return client, nil
}
