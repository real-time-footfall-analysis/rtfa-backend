package eventlivedata

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

var db live_db_adapter = &dynamodb_adapter{}

func Init(r *mux.Router) {

	db.initConn()

	r.HandleFunc("/live/heatmap", heatmapHandler).Methods("POST")
}

type mapRequest struct {
	EventId string `json:"event_id"`
}

func heatmapHandler(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)

	var req mapRequest

	err := decoder.Decode(&req)

	if err != nil {
		log.Println("Cannot decode request:", err)
		http.Error(
			writer,
			fmt.Sprintf("Failed to decode request: %s", err),
			http.StatusBadRequest)
		return
	}

	heatMap, _ := db.getLiveHeatMap()

	json.NewEncoder(writer).Encode(heatMap)

	return
}

func main() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1")},
	)
	if err != nil {
		fmt.Println("Got error creating session:")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	// Create DynamoDB client
	svc := dynamodb.New(sess)
	/*proj := expression.NamesList(expression.Name("pKey"), expression.Name("regionId"), expression.Name("eventId"))
	expr, err := expression.NewBuilder().WithProjection(proj).Build()
	*/
	if err != nil {
		fmt.Println("Got error creating expression:")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	params := &dynamodb.ScanInput{
		TableName: aws.String("current_position"),
	}
	result, err := svc.Scan(params)
	if err != nil {
		fmt.Println("Got error doing scan:")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	region_count := make(map[string]int, 0)
	for _, row := range result.Items {
		regionId := *row["regionId"].N
		count, ok := region_count[regionId]
		if !ok {
			region_count[regionId] = 1
			continue
		}
		region_count[regionId] = count + 1
	}
	fmt.Print(region_count)
}
