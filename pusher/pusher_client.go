package pusher

import (
	"flag"
	"github.com/pusher/pusher-http-go"
	"log"
	"os"
)

type PusherInterface interface {
	InitConn()
	SendItem(channelName string, eventName string, data []byte)
}

type PusherClient struct {
	client pusher.Client
}

func (pc *PusherClient) InitConn() {
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

func (pc *PusherClient) SendItem(channelName string, eventName string, data []byte) {
	_, err := pc.client.Trigger(channelName, eventName, data)
	if err != nil {
		log.Println("Got an error sending item to Pusher")
		log.Println(err.Error())
	}
}
