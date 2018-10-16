package locationupdate

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
)

type queue_adapter interface {
	initConn() error
	addLocationUpdate(event *update) error
}

type Movement struct {
	AttendeeID string `json:"attendee_id"`
	Location   string `json:"location"`
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
	s := session.New(&aws.Config{Region: aws.String(region)})

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
	putOutput, err := kq.kinesis.PutRecord(&kinesis.PutRecordInput{
		Data:         byteEncodedMov,
		StreamName:   aws.String(kq.streamName),
		PartitionKey: aws.String("key1"),
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", putOutput)
	// TODO: add event to queue_adapter
	return nil

}
