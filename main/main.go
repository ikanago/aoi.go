package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"github.com/bwmarrin/discordgo"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

func main() {
	projectID := os.Getenv("PROJECT_ID")
	secretID := os.Getenv("SECRET_ID")
	credential, err := accessSecretVersion(projectID, secretID)
	if err != nil {
		log.Fatal(err)
		return
	}

	discordClient, err := discordgo.New("Bot " + credential.DiscordToken)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer discordClient.Close()
	discordClient.AddHandler(onMessageCreate)

	if err = discordClient.Open(); err != nil {
		log.Fatal(err)
		return
	}

	if err = LoadFirestore(projectID); err != nil {
		log.Fatal(err)
		return
	}

	if err = InitClient(credential); err != nil {
		log.Fatal(err)
		return
	}
	stream, demux, err := InitStream(discordClient)
	defer stream.Stop()
	if err != nil {
		log.Fatal(err)
		return
	}
	go demux.HandleChan(stream.Messages)

	log.Println("Bot is running")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

// Credential is a struct to hold a token fetched from Google Secret Manager.
type Credential struct {
	ConsumerKey       string `json:"CONSUMER_KEY"`
	ConsumerSecret    string `json:"CONSUMER_SECRET"`
	AccessToken       string `json:"ACCESS_TOKEN"`
	AccessTokenSecret string `json:"ACCESS_TOKEN_SECRET"`
	DiscordToken      string `json:"DISCORD_TOKEN"`
}

func accessSecretVersion(projectID string, secretID string) (*Credential, error) {
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	secretURI := "projects/" + projectID + "/secrets/" + secretID + "/versions/latest"
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: secretURI,
	}

	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return nil, err
	}

	var credential Credential
	if err := json.Unmarshal(result.Payload.Data, &credential); err != nil {
		return nil, err
	}
	return &credential, nil
}

// onMessageCreate is called when there is a new message in a guild this bot is belogns to.
// If this bot is mentioned, parse command and do corresponding actions.
func onMessageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == session.State.User.ID || len(message.Mentions) == 0 || message.Mentions[0].Username != session.State.User.Username {
		return
	}

	command, err := ParseCommand(message.Content)
	if err != nil {
		log.Println(err)
		session.ChannelMessageSend(message.ChannelID, err.Error())
		return
	}

	if err := command.Handle(session, message.Message); err != nil {
		log.Println(err)
		session.ChannelMessageSend(message.ChannelID, err.Error())
	}
}
