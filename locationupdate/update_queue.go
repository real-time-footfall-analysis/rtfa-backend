package locationupdate

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"log"
	"os"
)

type queue_adapter interface {
	initConn() error
	addLocationUpdate(event *update) error
}

type kenisis_queue struct {
	kinesis    *kinesis.Kinesis
	streamName string
}

// initConn opens the connection to the location event kinesis queue
func (kq *kenisis_queue) initConn() error {
	// Define the stream name and the AWS region it's in
	stream := "movement_event_stream"
	region := "eu-central-1"
	// Create a new AWS session in the reqired region
	s, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}

	// Create a new kinesis adapter (assume stream exists
	kq.kinesis = kinesis.New(s)
	kq.streamName = stream

	return nil
}

// Pre: the event object is valid
func (kq *kenisis_queue) addLocationUpdate(event *update) error {
	// Encode a record into JSON bytes
	byteEncodedMov, _ := json.Marshal(event)

	// Send the record to Kinesis
	_, err := kq.kinesis.PutRecord(&kinesis.PutRecordInput{
		Data:         byteEncodedMov,
		StreamName:   aws.String(kq.streamName),
		PartitionKey: aws.String("key1"),
	})
	if err != nil {
		panic(err)
	}
	return nil

}
