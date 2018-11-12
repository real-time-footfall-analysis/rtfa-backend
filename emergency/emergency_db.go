package emergency

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
	"os"
)

const (
	DYNOMODB_TABLE = "emergency_events"
)

type emergencyDbAdapter interface {
	initConn() error
	getTableScan() (*dynamodb.ScanOutput, error)
	sendItem(req emergency_request) error
}

type dynamoDbAdaptor struct {
	db         *dynamodb.DynamoDB
	streamName string
}

// initConn opens the connection to the dynamo DB database
func (db *dynamoDbAdaptor) initConn() error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1")},
	)
	if err != nil {
		log.Println("Got error creating session:")
		log.Println(err.Error())
		os.Exit(1)
	}
	// Create DynamoDB client
	db.db = dynamodb.New(sess)

	return nil
}

// Get a scan of the entire table
func (db *dynamoDbAdaptor) getTableScan() (*dynamodb.ScanOutput, error) {
	params := &dynamodb.ScanInput{
		TableName: aws.String(DYNOMODB_TABLE),
	}
	result, err := db.db.Scan(params)
	if err != nil {
		log.Print("Got error doing scan:", err.Error())
		return nil, err
	}
	return result, nil
}

// Pre: the event object is valid
func (db *dynamoDbAdaptor) sendItem(req emergency_request) (err error) {
	// Encode the data
	encoded, err := dynamodbattribute.MarshalMap(req)
	if err != nil {
		fmt.Println("Got error trying to marshal request:")
		fmt.Println(err.Error())
		return
	}

	// Wrap the item up in a request
	input := &dynamodb.PutItemInput{
		Item:      encoded,
		TableName: aws.String(DYNOMODB_TABLE),
	}

	// Send the item
	_, err = db.db.PutItem(input)
	if err != nil {
		fmt.Println("Got error calling PutItem:")
		fmt.Println(err.Error())
		return
	}

	return nil
}
