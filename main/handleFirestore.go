package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	firestore "cloud.google.com/go/firestore"
)

const (
	memoCollectionName        string = "Memo"
	tweetFilterCollectionName string = "Filters"
)

// FilterDocument represents  a filter of tweet.
type FilterDocument struct {
	ID         string   `firestore:"id"`
	ScreenName string   `firestore:"screenName"`
	Keywords   []string `firestore:"keywords"`
	ChannelID  string   `firestore:"channelID"`
}

var firestoreClient *firestore.Client

var tweetFilters = make(map[string]*FilterDocument)

// GetFilter returns filter corresponding to given screen name.
func GetFilter(screenName string) *FilterDocument {
	if filter, exists := tweetFilters[screenName]; exists {
		return filter
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

// LoadFirestore initializes firestore client and fetch data from firestore.
func LoadFirestore(projectID string) (err error) {
	ctx := context.Background()
	firestoreClient, err = firestore.NewClient(ctx, projectID)

	filters, err := FetchFilters()
	if err != nil {
		return
	}

	for _, filter := range filters {
		tweetFilters[filter.ScreenName] = &filter
		screenNameToID[filter.ScreenName] = filter.ID
	}
	return
}

// CreateFilter creates data on firestore.
func CreateFilter(screenName string, keywords []string, channelID string) (err error) {
	if _, exists := tweetFilters[screenName]; exists {
		return errors.New("そのアカウントのフィルターは作成済みです\n`@Aoi tweet add ID KEYWORDS` を使ってください")
	}

	collection := firestoreClient.Collection(tweetFilterCollectionName)
	doc := collection.Doc(screenName)
	ctx := context.Background()
	id, err := GetUserID(screenName)
	if err != nil {
		return
	}

	filter := FilterDocument{
		ID:         id,
		ScreenName: screenName,
		Keywords:   keywords,
		ChannelID:  channelID,
	}
	_, err = doc.Create(ctx, filter)
	if err != nil {
		return
	}

	tweetFilters[screenName] = &filter
	screenNameToID[screenName] = id
	return
}

// AddFilter adds keywords to existing filter.
func AddFilter(screenName string, keywords []string) (updatedKeywords []string, err error) {
	if _, exists := tweetFilters[screenName]; !exists {
		return nil, errors.New("そのアカウントのフィルターは存在しません\n`@Aoi tweet create ID KEYWORDS` でフィルターを作ってください")
	}

	collection := firestoreClient.Collection(tweetFilterCollectionName)
	doc := collection.Doc(screenName)
	ctx := context.Background()
	updatedKeywords = append(tweetFilters[screenName].Keywords, keywords...)
	tweetFilters[screenName].Keywords = updatedKeywords
	_, err = doc.Update(ctx, []firestore.Update{
		{
			Path:  "keywords",
			Value: updatedKeywords,
		},
	})
	return
}

// RemoveFilter removes keywords from a specific filter.
func RemoveFilter(screenName string, keywords []string) (updatedKeywords []string, err error) {
	if _, exists := tweetFilters[screenName]; !exists {
		return nil, errors.New("そのアカウントのフィルターは存在しません\n`@Aoi tweet create ID KEYWORDS` でフィルターを作ってください")
	}

	collection := firestoreClient.Collection(tweetFilterCollectionName)
	doc := collection.Doc(screenName)
	ctx := context.Background()
	updatedKeywords = removeElement(tweetFilters[screenName].Keywords, keywords)
	tweetFilters[screenName].Keywords = updatedKeywords
	_, err = doc.Update(ctx, []firestore.Update{
		{
			Path:  "keywords",
			Value: updatedKeywords,
		},
	})
	return
}

func removeElement(array, toRemove []string) (removed []string) {
	for _, element := range array {
		shouldRemove := false
		for _, x := range toRemove {
			if element == x {
				shouldRemove = true
			}
		}
		if !shouldRemove {
			removed = append(removed, element)
		}
	}
	return
}

// DeleteFilter deletes a filter whose name is `screenName`.
func DeleteFilter(screenName string) (err error) {
	if _, exists := tweetFilters[screenName]; !exists {
		return errors.New("そのアカウントのフィルターは存在しません\n`@Aoi tweet create ID KEYWORDS` でフィルターを作ってください")
	}

	collection := firestoreClient.Collection(tweetFilterCollectionName)
	doc := collection.Doc(screenName)
	ctx := context.Background()
	_, err = doc.Delete(ctx)
	delete(tweetFilters, screenName)
	delete(screenNameToID, screenName)
	return
}

// ChangeFilterChannel changes a channel to which a tweet from `screenName` is sent.
func ChangeFilterChannel(screenName, channelID string) (err error) {
	if _, exists := tweetFilters[screenName]; !exists {
		return errors.New("そのアカウントのフィルターは存在しません\n`@Aoi tweet create ID KEYWORDS` でフィルターを作ってください")
	}

	collection := firestoreClient.Collection(tweetFilterCollectionName)
	doc := collection.Doc(screenName)
	ctx := context.Background()
	tweetFilters[screenName].ChannelID = channelID
	_, err = doc.Update(ctx, []firestore.Update{
		{
			Path:  "channelID",
			Value: channelID,
		},
	})
	return
}

// FetchFilters gets filters from firestore.
func FetchFilters() (filters []FilterDocument, err error) {
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

// CreateMemo creates a new memo on firestore.
func CreateMemo(channelID string, text string) (message string, err error) {
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

// FetchMemo fetches existing memos belong to a given channel.
func FetchMemo(channelID string) (memos []MemoDocument, err error) {
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
