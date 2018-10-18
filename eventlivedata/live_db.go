package eventlivedata

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"os"
	"strconv"
)

type live_db_adapter interface {
	initConn() error
	getLiveHeatMap(event int) (map[string]int, error)
}

type dynamodbAdaptor struct {
	db         *dynamodb.DynamoDB
	streamName string
}

// initConn opens the connection to the location event kinesis queue
func (db *dynamodbAdaptor) initConn() error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1")},
	)
	if err != nil {
		fmt.Println("Got error creating session:")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	// Create DynamoDB client
	db.db = dynamodb.New(sess)

	return nil
}

// Pre: the event object is valid
func (db *dynamodbAdaptor) getLiveHeatMap(event int) (map[string]int, error) {
	params := &dynamodb.ScanInput{
		TableName: aws.String("current_position"),
	}
	result, err := db.db.Scan(params)
	if err != nil {
		fmt.Println("Got error doing scan:")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	region_count := make(map[string]int, 0)
	for _, row := range result.Items {
		eventId, _ := strconv.Atoi(*row["eventId"].N)
		if eventId == event {
			regionId := *row["regionId"].N
			count, ok := region_count[regionId]
			if !ok {
				region_count[regionId] = 1
			} else {
				region_count[regionId] = count + 1
			}
		}
	}
	fmt.Print(region_count)
	return region_count, nil

}
