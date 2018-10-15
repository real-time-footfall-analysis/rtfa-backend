package eventstaticdata

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	// Postgres driver
	_ "github.com/lib/pq"
)

const (
	DATABASE_HOST string = "rtfa-crowd-static-data.chh9za0s2bso.eu-central-1.rds.amazonaws.com"
	DATABASE_NAME string = "rtfa_crowd_static_data"
)

type Event struct {
	ID            string `json:"id"`
	OrganiserID   string `json:"organiserId"`
	Name          string `json:"name"`
	Location      string `json:"location"`
	StartDate     string `json:"startDate"`
	EndDate       string `json:"endDate"`
	MaxAttendance string `json:"maxAttendance"`
	CoverPhotoURL string `json:"coverPhotoUrl"`
}

type Map struct {
	EventID     string `json:"eventId"`
	ImageURL    string `json:"imageUrl"`
	ImageWidth  string `json:"imageWidth"`
	ImageHeight string `json:"imageHeight"`
	MapWidth    string `json:"mapWidth"`
}

var db *sql.DB

// initConn opens the connection to postgres and sets the global db variable
func initConn() {

	username := os.Getenv("RTFA_STATICDATA_DB_USER")
	if username == "" {
		log.Fatal("RTFA_STATICDATA_DB_USER not set.")
	}
	password := os.Getenv("RTFA_STATICDATA_DB_PASSWORD")
	if password == "" {
		log.Fatal("RTFA_STATICDATA_DB_PASSWORD not set.")
	}

	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s",
		DATABASE_HOST,
		username,
		password,
		DATABASE_NAME)

	openedDb, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	db = openedDb

}

// Pre: the event object is valid
func addEvent(event *Event) (*Event, error) {

	// TODO: and insert into database
	// probably with a raw SQL call
	return nil, nil

}
