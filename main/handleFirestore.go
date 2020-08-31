package main

import (
	"context"
	"errors"
	"fmt"
	"log"
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
	ID         string   `firestore:"id"`
	ScreenName string   `firestore:"screenName"`
	Keywords   []string `firestore:"keywords"`
	ChannelID  string   `firestore:"channelID"`
}

var firestoreClient *firestore.Client

var tweetFilters = make(map[string]FilterDocument)

// GetFilter returns filter corresponding to given screen name.
func GetFilter(screenName string) *FilterDocument {
	if filter, exists := tweetFilters[screenName]; exists {
		return &filter
	}
	return nil
}

// Map twitter screen name (@NAME) to id string (111111111111).
var screenNameToID = make(map[string]string)

// GetAllIDs yields all registerd filters
func GetAllIDs() (ids []string) {
	for _, id := range screenNameToID {
		ids = append(ids, id)
	}
	return
}

func loadFirestore(projectID string) (err error) {
	ctx := context.Background()
	firestoreClient, err = firestore.NewClient(ctx, projectID)

	filters, err := fetchFilters()
	if err != nil {
		return
	}

	for _, filter := range filters {
		tweetFilters[filter.ScreenName] = filter
		screenNameToID[filter.ScreenName] = filter.ID
	}
	return
}

func createFilter(screenName string, filters []string) (message string, err error) {
	if _, exists := tweetFilters[screenName]; exists {
		return "", errors.New("そのアカウントのフィルターは作成済みです\n`@Aoi tweet add ID KEYWORDS` を使ってください")
	}

	collection := firestoreClient.Collection(tweetFilterCollectionName)
	doc := collection.Doc(screenName)
	ctx := context.Background()
	id, err := getUserID(screenName)
	if err != nil {
		return
	}

	filter := FilterDocument{
		ID:         id,
		ScreenName: screenName,
		Keywords:   filters,
		ChannelID:  defaultChannelID,
	}
	_, err = doc.Create(ctx, filter)
	if err != nil {
		return
	}

	tweetFilters[screenName] = filter
	message = fmt.Sprintf("@%s のフィルターを作成しました\n現在のキーワード: %s", screenName, strings.Join(filters, ", "))
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
