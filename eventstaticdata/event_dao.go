package eventstaticdata

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-pg/pg"
)

const (
	DATABASE_HOST string = "rtfa-crowd-static-data.chh9za0s2bso.eu-central-1.rds.amazonaws.com"
	DATABASE_PORT int    = 5432
	DATABASE_NAME string = "rtfa_crowd_static_data"
)

type Event struct {
	tableName     struct{}  `sql:"event"`
	ID            int32     `json:"eventID,omitempty"`
	OrganiserID   int32     `json:"organiserID"`
	Name          string    `json:"name"`
	Location      string    `json:"location"`
	StartDate     time.Time `json:"startDate"`
	EndDate       time.Time `json:"endDate"`
	IndoorOutdoor string    `json:"indoorOutdoor"`
	MaxAttendance int64     `json:"maxAttendance"`
	CoverPhotoURL string    `json:"coverPhotoUrl"`
}

type AllEventsRequest struct {
	OrganiserID int32 `json:"organiserId,string"`
}

type Map struct {
	tableName struct{} `sql:"map"`
	ID        int32    `json:"mapID,omitempty"`
	Type      string   `json:"type"`
	Zoom      int32    `json:"zoom"`
	EventID   int32    `json:"eventID"`
	Lat       float64  `json:"lat"`
	Lng       float64  `json:"lng"`
}

type Region struct {
	tableName struct{} `sql:"region"`
	ID        int32    `json:"regionID,omitempty"`
	Name      string   `json:"name"`
	Type      string   `json:"type"`
	UUID      string   `json:"uuid,omitempty"`
	Major     int32    `json:"major,omitempty"`
	Minor     int32    `json:"minor,omitempty"`
	Lat       float64  `json:"lat,omitempty"`
	Lng       float64  `json:"lng,omitempty"`
	Radius    int32    `json:"radius,omitempty"`
	EventID   int32    `json:"eventID"`
}

var dbUsername string
var dbPassword string

// fetchEnvVars fetches the environment variables used for database connection
func fetchEnvVars() {

	dbUsername = os.Getenv("RTFA_STATICDATA_DB_USER")
	if dbUsername == "" && flag.Lookup("test.v") == nil {
		log.Fatal("RTFA_STATICDATA_DB_USER not set.")
	}
	dbPassword = os.Getenv("RTFA_STATICDATA_DB_PASSWORD")
	if dbPassword == "" && flag.Lookup("test.v") == nil {
		log.Fatal("RTFA_STATICDATA_DB_PASSWORD not set.")
	}

}

func connectDB() *pg.DB {

	return pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:%d", DATABASE_HOST, DATABASE_PORT),
		User:     dbUsername,
		Password: dbPassword,
		Database: DATABASE_NAME,
	})

}

// All the add{Object} functions assume that the object is valid
// and modify the object, inserting the new ID

func addEvent(event *Event) error {

	db := connectDB()
	defer db.Close()

	err := db.Insert(event)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil

}

func getAllEventsByOrganiserID(organiserID int32) ([]Event, error) {

	db := connectDB()
	defer db.Close()

	var events []Event

	err := db.Model(&events).Where("organiser_id = ?", organiserID).Select()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return events, nil

}

func getEventByID(id int) (*Event, error) {

	db := connectDB()
	defer db.Close()

	event := &Event{ID: int32(id)}
	err := db.Select(event)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return event, nil

}

func addMap(eventMap *Map) error {

	db := connectDB()
	defer db.Close()

	err := db.Insert(eventMap)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil

}

func addRegions(regions *[]Region) error {

	db := connectDB()
	defer db.Close()

	_, err := db.Model(regions).Insert()
	if err != nil {
		log.Println(err)
		return err
	}

	return nil

}

func updateRegion(region *Region) error {

	db := connectDB()
	defer db.Close()

	if err := db.Update(region); err != nil {
		log.Println(err)
		return err
	}

	return nil

}

func getRegionsByEventID(eventID int) (*[]Region, error) {

	db := connectDB()
	defer db.Close()

	var regions = make([]Region, 0)
	err := db.Model(&regions).Where("event_id = ?", eventID).Select()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &regions, nil

}

func getRegionByID(eventID, regionID int) (*Region, error) {

	db := connectDB()
	defer db.Close()

	var region Region

	err := db.Model(&region).Where("id = ?", regionID).Where("event_id = ?", eventID).Select()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &region, nil

}
