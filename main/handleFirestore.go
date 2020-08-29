package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	firestore "cloud.google.com/go/firestore"
)

const memoCollectionName = "Memo"

var firestoreClient *firestore.Client

func loadFirestore() (err error) {
	ctx := context.Background()
	projectID := os.Getenv("PROJECT_ID")
	firestoreClient, err = firestore.NewClient(ctx, projectID)
	return
}

// MemoDocument represents "Memo" document in the firestore.
type MemoDocument struct {
	ChannelID string    `firestore:"channelID"`
	Text      string    `firestore:"text"`
	Timestamp time.Time `firestore:"timestamp"`
}

func addMemo(channelID string, text string) (message string, err error) {
	currentTime := time.Now()
	documentID := currentTime.String()
	
	collection := firestoreClient.Collection(memoCollectionName)
	doc := collection.Doc(documentID)
	ctx := context.Background()
	_, err = doc.Create(ctx, MemoDocument {
		ChannelID: channelID,
		Text: text,
		Timestamp: currentTime,
	})
	if err != nil {
		log.Println(err)
		return
	}

	message = fmt.Sprintf("%s を記録しました", text)
	return
}
