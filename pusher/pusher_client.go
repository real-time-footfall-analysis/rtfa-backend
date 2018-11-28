package pusher

import (
	"flag"
	"github.com/pusher/push-notifications-go"
	"github.com/pusher/pusher-http-go"
	"log"
	"os"
)

type PusherChannelInterface interface {
	InitConn()
	SendItem(channelName string, eventName string, data []byte)
}

type PusherChannelClient struct {
	client pusher.Client
}

func (pc *PusherChannelClient) InitConn() {
	// Get the secret key
	secretKey := os.Getenv("RTFA_PUSHER_SECRET_KEY")
	if secretKey == "" && flag.Lookup("test.v") == nil {
		log.Fatal("RTFA_PUSHER_SECRET_KEY not set.")
	}

	client := pusher.Client{
		AppId:   "648875",
		Key:     "544e69db41ad4dcc08db",
		Secret:  secretKey,
		Cluster: "eu",
	}

	pc.client = client
}

func (pc *PusherChannelClient) SendItem(channelName string, eventName string, data []byte) {
	_, err := pc.client.Trigger(channelName, eventName, data)
	if err != nil {
		log.Println("Got an error sending item to Pusher channel")
		log.Println(err.Error())
	}
}

type PusherBeamsInterface interface {
	InitConn()
	SendNotification(regionIds []string, title string, body string) (publishId string, err error)
}

type PusherBeamsClient struct {
	client pushnotifications.PushNotifications
}

func (pbc *PusherBeamsClient) InitConn() {
	const instanceId = "5cb5ee8c-bcd9-4b07-ab76-95220dc679c1"

	// Get the secret key
	secretKey := os.Getenv("RTFA_PUSHER_BEAMS_SECRET_KEY")
	if secretKey == "" && flag.Lookup("test.v") == nil {
		log.Fatal("RTFA_PUSHER_BEAMS_SECRET_KEY not set.")
	}

	// Create the notification client
	client, err := pushnotifications.New(instanceId, secretKey)
	if err != nil {
		log.Println("Error starting up pusher beams")
		os.Exit(1)
	}
	pbc.client = client
}

func (pbc *PusherBeamsClient) SendNotification(regionIds []string, title string, body string) (publishId string, err error) {
	// Make the request
	publishRequest := map[string]interface{}{
		"apns": map[string]interface{}{
			"aps": map[string]interface{}{
				"alert": map[string]interface{}{
					"title": title,
					"body":  body,
				},
			},
		},
	}

	// Send the notification
	publishId, err = pbc.client.Publish(regionIds, publishRequest)
	if err != nil {
		log.Println("Error sending to pusher beam")
		log.Println(err)
	}
	return publishId, err
}
