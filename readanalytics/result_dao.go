package readanalytics

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const (
	AWS_REGION   = "eu-central-1"
	DYNAMO_TABLE = "analytics_results"
	DYNAMO_PK    = "EventID-TaskID"
)

func fetchAnalyticsResult(eventID, taskID int) (map[string]interface{}, error) {

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(AWS_REGION)},
	)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	svc := dynamodb.New(sess)

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(DYNAMO_TABLE),
		Key: map[string]*dynamodb.AttributeValue{
			DYNAMO_PK: {
				S: aws.String(fmt.Sprintf("%d-%d", eventID, taskID)),
			},
		},
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	m := map[string]interface{}{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &m)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return m, nil
}
