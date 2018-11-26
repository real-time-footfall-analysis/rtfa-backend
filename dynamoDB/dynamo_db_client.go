package dynamoDB

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
)

type DynamoDBInterface interface {
	InitConn(tableName string) error
	GetTableScan() []interface{}
	SendItem(req interface{})
}

type DynamoDBClient struct {
	connection *dynamodb.DynamoDB
	tableName  string
}

// initConn opens the connection to the dynamo DB database
func (db *DynamoDBClient) InitConn(tableName string) error {
	// Save the table name
	db.tableName = tableName

	// Create a session in a given AWS region
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1")},
	)
	if err != nil {
		log.Println("Got error creating session:")
		log.Println(err.Error())
		return err
	}

	// Create DynamoDB client
	db.connection = dynamodb.New(sess)
	return nil
}

// Get a scan of the entire table
func (db *DynamoDBClient) GetTableScan() []interface{} {
	// Take a scan of the table
	params := &dynamodb.ScanInput{
		TableName: aws.String(db.tableName),
	}
	result, err := db.connection.Scan(params)

	// Check for errors
	if err != nil {
		log.Println("Got error doing scan:", err.Error())
		return nil
	}

	// Create list to store result in
	var allRows = make([]interface{}, len(result.Items))

	// Unmarshall to list of maps
	for index, row := range result.Items {
		err = dynamodbattribute.UnmarshalMap(row, &allRows[index])
		if err != nil {
			log.Println("Got error unmarshalling:", err.Error())
			return nil
		}
	}
	return allRows
}

// Pre: the event object is valid
func (db *DynamoDBClient) SendItem(req interface{}) {
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
		TableName: aws.String(db.tableName),
	}

	// Send the item
	_, err = db.connection.PutItem(input)
	if err != nil {
		log.Println("Got an error putting item in DynamoDB")
		log.Println(err.Error())
	}
}
