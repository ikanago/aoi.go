package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	firestore "cloud.google.com/go/firestore"
)

const (
	memoCollectionName        string = "Memo"
	tweetFilterCollectionName string = "Filters"
	defaultChannelID          string = "626795866556203021"
)

// FilterDocument represents  a filter of tweet.
type FilterDocument struct {
	ID        string `firestore:"id"`
	Keywords  []string `firestore:"keywords"`
	ChannelID string `firestore:"channelID"`
}

var firestoreClient *firestore.Client

var tweetFilters = make(map[string]FilterDocument)

func loadFirestore() (err error) {
	ctx := context.Background()
	projectID := os.Getenv("PROJECT_ID")
	firestoreClient, err = firestore.NewClient(ctx, projectID)

	filters, err := fetchFilters()
	if err != nil {
		return
	}
	for _, filter := range filters {
		tweetFilters[filter.ID] = filter
	}
	return
}

func createFilter(id string, filters []string) (message string, err error) {
	if _, exists := tweetFilters[id]; exists {
		return "", errors.New("そのアカウントのフィルターは作成済みです\n`@Aoi tweet add ID KEYWORDS` を使ってください")
	}

	collection := firestoreClient.Collection(tweetFilterCollectionName)
	doc := collection.Doc(id)
	ctx := context.Background()
	filter := FilterDocument{
		ID:        id,
		Keywords:  filters,
		ChannelID: defaultChannelID,
	}
	_, err = doc.Create(ctx, filter)
	if err != nil {
		return
	}

	tweetFilters[id] = filter
	message = fmt.Sprintf("@%s のフィルターを作成しました\n現在のキーワード: %s", id, strings.Join(filters, ", "))
	return
}

func fetchFilters() (filters []FilterDocument, err error) {
	collection := firestoreClient.Collection(tweetFilterCollectionName)
	ctx := context.Background()
	docs, err := collection.Documents(ctx).GetAll()
	if err != nil {
		return
	}

	for _, doc := range docs {
		var filter FilterDocument
		if err = doc.DataTo(&filter); err != nil {
			return
		}
		filters = append(filters, filter)
	}
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
	_, err = doc.Create(ctx, MemoDocument{
		ChannelID: channelID,
		Text:      text,
		Timestamp: currentTime,
	})
	if err != nil {
		log.Println(err)
		return
	}

	message = fmt.Sprintf("%s を記録しました", text)
	return
}

func fetchMemo(channelID string) (memos []MemoDocument, err error) {
	collection := firestoreClient.Collection(memoCollectionName)
	query := collection.Where("channelID", "==", channelID).OrderBy("timestamp", firestore.Asc)
	ctx := context.Background()
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return
	}

	for _, doc := range docs {
		var memo MemoDocument
		if err = doc.DataTo(&memo); err != nil {
			return
		}
		memos = append(memos, memo)
	}
	return
}
