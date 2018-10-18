package eventstaticdata

import (
	"testing"
	"time"
)

func TestInvalidEventThrowsError(t *testing.T) {

	validEvent := Event{
		Name:          "Google DevFest 2018",
		OrganiserID:   6,
		Location:      "Huxley, Imperial College London",
		StartDate:     time.Date(2018, 12, 05, 0, 0, 0, 0, time.UTC),
		EndDate:       time.Date(2018, 12, 10, 0, 0, 0, 0, time.UTC),
		IndoorOutdoor: "indoor",
		MaxAttendance: 500,
		CoverPhotoURL: "http://imgur.com/image.jpg"}
	if err := validateEvent(&validEvent); err != nil {
		t.Errorf("Not expecting an error when validating a valid event: %+v",
			validEvent)
	}

	invalidEvent := Event{
		Name:          "Google DevFest 2018",
		OrganiserID:   6,
		Location:      "",
		StartDate:     time.Date(2018, 12, 10, 0, 0, 0, 0, time.UTC),
		EndDate:       time.Date(2018, 12, 05, 0, 0, 0, 0, time.UTC),
		IndoorOutdoor: "neither",
		MaxAttendance: 0,
		CoverPhotoURL: "image"}
	if err := validateEvent(&invalidEvent); err == nil {
		t.Errorf("Expecting an error when validating an invalid event: %+v",
			invalidEvent)
	}

}
func TestValidateOrganiserId(t *testing.T) {

	if err := validateOrganiserID(0); err == nil {
		t.Error("Expected an error when validating an empty event Organiser ID")
	}

	const validOrganiserID int32 = 52
	if err := validateOrganiserID(validOrganiserID); err != nil {
		t.Errorf("Not expecting an error when validating a non-empty "+
			"event Organiser ID: %d", validOrganiserID)
	}

}

func TestValidateName(t *testing.T) {

	if err := validateName(""); err == nil {
		t.Error("Expected an error when validating an empty event Name")
	}

	validEventName := "Google DevFest 2018"
	if err := validateName(validEventName); err != nil {
		t.Errorf("Not expecting an error when validating a non-empty event Name: %s",
			validEventName)
	}

}

func TestValidateLocation(t *testing.T) {

	if err := validateLocation(""); err == nil {
		t.Error("Expected an error when validating an empty event Location")
	}

	validEventLocation := "Huxley, Imperial College London"
	if err := validateLocation(validEventLocation); err != nil {
		t.Errorf("Not expecting an error when validating a non-empty event "+
			"Location: %s", validEventLocation)
	}

}

func TestValidateDates(t *testing.T) {

	validStartDate := time.Date(2018, 12, 05, 0, 0, 0, 0, time.UTC)
	validEndDate := time.Date(2018, 12, 10, 0, 0, 0, 0, time.UTC)

	if err := validateDates(time.Time{}, validEndDate); err == nil {
		t.Error("Expected an error when validating an empty event Start Date")
	}
	if err := validateDates(validStartDate, time.Time{}); err == nil {
		t.Error("Expected an error when validating an empty event End Date")
	}

	if err := validateDates(validEndDate, validStartDate); err == nil {
		t.Error("Expected an error when validating start date and end dates" +
			" that are in the wrong order")
	}

	if err := validateDates(validStartDate, validEndDate); err != nil {
		t.Error("Not expecting an error when validating valid event dates")
	}
	if err := validateDates(validStartDate, validStartDate); err != nil {
		t.Error("Not expecting an error when validating valid identical event dates")
	}

}

func TestValidateIndoorOutdoor(t *testing.T) {

	if err := validateIndoorOutdoor(""); err == nil {
		t.Error("Expected an error when validating an empty event Indoor/Outdoor")
	}

	validIndoorOutdoor := "indoor"
	if err := validateIndoorOutdoor(validIndoorOutdoor); err != nil {
		t.Errorf("Not expecting an error when validating a non-empty Indoor/Outdoor: %s",
			validIndoorOutdoor)
	}

}

func TestValidateCoverPhotoURL(t *testing.T) {

	if err := validateCoverPhotoURL(""); err == nil {
		t.Error("Expecting an error when validating an empty event" +
			" Cover Photo URL")
	}
	invalidUrl := "image"
	if err := validateCoverPhotoURL(invalidUrl); err == nil {
		t.Errorf("Expecting an error when validating an event Cover "+
			"Photo URL which is not an absolute URL: %s", invalidUrl)
	}
	validUrl := "http://imgur.com/image.jpg"
	if err := validateCoverPhotoURL("http://imgur.com/image.jpg"); err != nil {
		t.Errorf("Not expecting an error when validating a valid event "+
			"Cover Photo URL: %s", validUrl)
	}

}
