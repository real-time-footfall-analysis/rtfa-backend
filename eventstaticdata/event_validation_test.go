package eventstaticdata

import "testing"

func TestInvalidEventThrowsError(t *testing.T) {

	validEvent := Event{
		Name:          "Google DevFest 2018",
		OrganiserID:   "DocSoc",
		Location:      "Huxley, Imperial College London",
		StartDate:     "2018-12-05",
		EndDate:       "2018-12-10",
		MaxAttendance: "500",
		CoverPhotoURL: "http://imgur.com/image.jpg"}
	if err := validateEvent(&validEvent); err != nil {
		t.Errorf("Not expecting an error when validating a valid event: %+v",
			validEvent)
	}

	invalidEvent := Event{
		Name:          "Google DevFest 2018",
		OrganiserID:   "DocSoc",
		Location:      "",
		StartDate:     "2018-12-10",
		EndDate:       "2018-12-05",
		MaxAttendance: "many",
		CoverPhotoURL: "image"}
	if err := validateEvent(&invalidEvent); err == nil {
		t.Errorf("Expecting an error when validating an invalid event: %+v",
			invalidEvent)
	}

}
func TestValidateOrganiserId(t *testing.T) {

	if err := validateOrganiserID(""); err == nil {
		t.Error("Expected an error when validating an empty event Organiser ID")
	}

	validOrganiserID := "DocSoc"
	if err := validateOrganiserID(validOrganiserID); err != nil {
		t.Errorf("Not expecting an error when validating a non-empty "+
			"event Organiser ID: %s", validOrganiserID)
	}

}

func TestValidateName(t *testing.T) {

	if err := validateName(""); err == nil {
		t.Error("Expected an error when validating an empty event Name")
	}

	validEventName := "Google DevFest 2018"
	if err := validateOrganiserID(validEventName); err != nil {
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

	if err := validateDates("", "2018-12-05"); err == nil {
		t.Error("Expected an error when validating an empty event Start Date")
	}
	if err := validateDates("2018-12-05", ""); err == nil {
		t.Error("Expected an error when validating an empty event End Date")
	}

	invalidStartDate := "05 Dec 18 10:00 UTC"
	if err := validateDates(invalidStartDate, "2018-12-10"); err == nil {
		t.Errorf("Expected an error when validating an incorrectly formatted"+
			" event Start Date: %s", invalidStartDate)
	}

	invalidEndDate := "10 Dec 18 10:00 UTC"
	if err := validateDates("2018-12-05", invalidEndDate); err == nil {
		t.Errorf("Expected an error when validating an incorrectly formatted"+
			" event End Date: %s", invalidEndDate)
	}

	if err := validateDates("2018-12-10", "2018-12-05"); err == nil {
		t.Error("Expected an error when validating start date and end dates" +
			" that are in the wrong order")
	}

	if err := validateDates("2018-12-05", "2018-12-10"); err != nil {
		t.Error("Not expecting an error when validating valid event dates")
	}
	if err := validateDates("2018-12-05", "2018-12-05"); err != nil {
		t.Error("Not expecting an error when validating valid identical event dates")
	}

}

func TestValidateMaxAttendance(t *testing.T) {

	if err := validateMaxAttendance(""); err == nil {
		t.Error("Expecting an error when validating an empty event Max Attendance")
	}

	invalidMaxAttendance := "many"
	if err := validateMaxAttendance("many"); err == nil {
		t.Errorf("Expecting an error when validating an event Max Attendance"+
			" that isn't a number: %s", invalidMaxAttendance)
	}

	validMaxAttendance := "15000"
	if err := validateMaxAttendance(validMaxAttendance); err != nil {
		t.Error("Not expecting an error when validating a valid event"+
			" Max Attendance: %s", validMaxAttendance)
	}

}

func TestValidateCoverPhotoUrl(t *testing.T) {

	if err := validateCoverPhotoUrl(""); err == nil {
		t.Error("Expecting an error when validating an empty event" +
			" Cover Photo URL")
	}
	invalidUrl := "image"
	if err := validateCoverPhotoUrl(invalidUrl); err == nil {
		t.Errorf("Expecting an error when validating an event Cover "+
			"Photo URL which is not an absolute URL: %s", invalidUrl)
	}
	validUrl := "http://imgur.com/image.jpg"
	if err := validateCoverPhotoUrl("http://imgur.com/image.jpg"); err != nil {
		t.Errorf("Not expecting an error when validating a valid event "+
			"Cover Photo URL: %s", validUrl)
	}

}
