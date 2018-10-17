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
	if err := validateMaxAttendance(event.MaxAttendance); err != nil {
		return err
	}
	if err := validateCoverPhotoURL(event.CoverPhotoURL); err != nil {
		return err
	}

	return nil
}

func validateMap(eventMap *Map) error {

	if err := validateType(eventMap.Type); err != nil {
		return err
	}
	if err := validateZoom(eventMap.Zoom); err != nil {
		return err
	}

	return nil

}

func validateAllEventsRequest(allEventsReq *AllEventsRequest) error {

	if allEventsReq.OrganiserID == 0 {
		return errors.New("Missing Organiser ID")
	}

	return nil

}

func validateOrganiserID(organiserID int32) error {

	if organiserID == 0 {
		return errors.New("Missing Event Organiser ID")
	}

	return nil
}

func validateName(name string) error {
	if name == "" {
		return errors.New("Missing Event Name")
	}
	return nil
}

func validateLocation(location string) error {

	if location == "" {
		return errors.New("Missing Event Location")
	}

	return nil
}

func validateDates(startDate, endDate time.Time) error {

	if startDate.IsZero() {
		return errors.New("Missing Event Start Date")
	}

	if endDate.IsZero() {
		return errors.New("Missing Event End Date")
	}

	if !startDate.Before(endDate) && !startDate.Equal(endDate) {
		return errors.New("The Event Start Date must be before the End Date")
	}

	return nil
}

func validateIndoorOutdoor(indoorOutdoor string) error {

	if indoorOutdoor == "" {
		return errors.New("Missing Event Indoor/Outdoor")
	}

	if indoorOutdoor != "indoor" &&
		indoorOutdoor != "outdoor" {
		return errors.New("The Event Indoor/Outdoor field must be " +
			"\"indoor\" or \"outdoor\"")
	}

	return nil

}

func validateMaxAttendance(maxAttendance int64) error {

	if maxAttendance == 0 {
		return errors.New("An Event Max Attendance > 0 must be specified")
	}

	return nil
}

func validateCoverPhotoURL(coverPhotoURL string) error {

	if coverPhotoURL == "" {
		return errors.New("Missing Event Cover Photo URL")
	}
	if _, err := url.ParseRequestURI(coverPhotoURL); err != nil {
		return errors.New("The Event Cover Photo URL must be a valid URL")
	}

	return nil
}

func validateType(mapType string) error {

	if mapType == "" {
		return errors.New("Missing Map Type")
	}
	if mapType != "image" &&
		mapType != "realMap" {
		return errors.New("The Map Type should be one of " +
			" \"image\" or \"realMap")
	}

	return nil

}

func validateZoom(zoom int32) error {

	if zoom == 0 {
		return errors.New("Missing Map Zoom")
	}

	return nil

}
