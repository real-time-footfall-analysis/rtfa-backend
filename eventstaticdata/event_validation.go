package eventstaticdata

import (
	"errors"
	"net/url"
	"time"
)

func validateEvent(event *Event) error {

	if err := validateOrganiserID(event.OrganiserID); err != nil {
		return err
	}
	if err := validateName(event.Name); err != nil {
		return err
	}
	if err := validateLocation(event.Location); err != nil {
		return err
	}
	if err := validateDates(event.StartDate, event.EndDate); err != nil {
		return err
	}
	if err := validateIndoorOutdoor(event.IndoorOutdoor); err != nil {
		return err
	}
	if err := validateRequiredInt(int(event.MaxAttendance), "event max attendance"); err != nil {
		return err
	}
	if err := validateCoverPhotoURL(event.CoverPhotoURL); err != nil {
		return err
	}

	return nil
}

func validateMap(eventMap *Map, eventID int) error {

	if err := validateMapType(eventMap.Type); err != nil {
		return err
	}
	if err := validateRequiredInt(int(eventMap.Zoom), "map zoom field"); err != nil {
		return err
	}
	if err := validateEventID(eventMap.EventID, eventID); err != nil {
		return err
	}

	return nil

}

func validateRegion(region *Region, eventID int) error {

	if err := validateRequiredString(region.Name, "region name"); err != nil {
		return err
	}
	if err := validateRegionType(region.Type); err != nil {
		return err
	}
	if err := validateEventID(region.EventID, eventID); err != nil {
		return err
	}

	return nil

}

func validateOrganiserID(organiserID int32) error {

	if organiserID == 0 {
		return errors.New("Missing event organiser ID")
	}

	return nil
}

func validateName(name string) error {
	if name == "" {
		return errors.New("Missing event name")
	}
	return nil
}

func validateLocation(location string) error {

	if location == "" {
		return errors.New("Missing event location")
	}

	return nil
}

func validateDates(startDate, endDate time.Time) error {

	if startDate.IsZero() {
		return errors.New("Missing event start date")
	}
	if endDate.IsZero() {
		return errors.New("Missing event end date")
	}

	if !startDate.Before(endDate) && !startDate.Equal(endDate) {
		return errors.New("The event start date must be before the event end date")
	}

	return nil
}

func validateIndoorOutdoor(indoorOutdoor string) error {

	if indoorOutdoor == "" {
		return errors.New("Missing event indoorOutdoor field")
	}

	if indoorOutdoor != "indoor" &&
		indoorOutdoor != "outdoor" {
		return errors.New("The event indoorOutdoor field must be " +
			"\"indoor\" or \"outdoor\"")
	}

	return nil

}

func validateCoverPhotoURL(coverPhotoURL string) error {

	if coverPhotoURL == "" {
		return errors.New("Missing event cover photo URL")
	}
	if _, err := url.ParseRequestURI(coverPhotoURL); err != nil {
		return errors.New("The event cover photo URL must be a valid URL")
	}

	return nil
}

func validateMapType(mapType string) error {

	if mapType == "" {
		return errors.New("Missing map type")
	}
	if mapType != "image" &&
		mapType != "realMap" {
		return errors.New("The map type should be one of " +
			" \"image\" or \"realMap")
	}

	return nil

}

func validateRequiredInt(num int, label string) error {

	if num <= 0 {
		return errors.New("Missing a positive " + label)
	}

	return nil

}

func validateRequiredString(str, label string) error {

	if str == "" {
		return errors.New("Missing " + label)
	}

	return nil

}

func validateRegionType(regionType string) error {

	if regionType == "" {
		return errors.New("Missing region type")
	}
	if regionType != "gps" &&
		regionType != "beacon" {
		return errors.New("The region type should be one of " +
			" \"gps\" or \"beacon")
	}

	return nil

}

func validateEventID(objEventID int32, URLEventID int) error {

	if objEventID == 0 {
		return errors.New("Missing event ID in object")
	}
	if int(objEventID) != URLEventID {
		return errors.New("The event ID given in the region object " +
			"is different to the one given in the request URL")
	}

	return nil

}
