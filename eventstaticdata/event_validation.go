package eventstaticdata

import (
	"errors"
	"net/url"
	"strconv"
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
	if err := validateMaxAttendance(event.MaxAttendance); err != nil {
		return err
	}
	if err := validateCoverPhotoUrl(event.CoverPhotoURL); err != nil {
		return err
	}

	return nil

}

func validateOrganiserID(organiserID string) error {
	if organiserID == "" {
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

func validateDates(startDate, endDate string) error {

	// To be consistent with the time package
	// (see https://stackoverflow.com/a/50557983/3837124)
	const RFC3339FullDate = "2006-01-02"

	if startDate == "" {
		return errors.New("Missing Event Start Date")
	}
	if endDate == "" {
		return errors.New("Missing Event End Date")
	}
	startTime, err := time.Parse(RFC3339FullDate, startDate)
	if err != nil {
		return errors.New("The Event Start Date must be of the form \"YYYY-MM-DD\"")
	}
	endTime, err := time.Parse(RFC3339FullDate, endDate)
	if err != nil {
		return errors.New("The Event End Date must be of the form \"YYYY-MM-DD\"")
	}
	if !startTime.Before(endTime) && !startTime.Equal(endTime) {
		return errors.New("The Event Start Date must be before the End Date")
	}

	return nil
}

func validateMaxAttendance(maxAttendance string) error {

	if maxAttendance == "" {
		return errors.New("Missing Event Max Attendance")
	}

	if _, err := strconv.Atoi(maxAttendance); err != nil {
		return errors.New("The Event Max Attendance value must be a number " +
			" and be between 0 and 9223372036854775807")
	}

	return nil
}

func validateCoverPhotoUrl(coverPhotoURL string) error {

	if coverPhotoURL == "" {
		return errors.New("Missing Event Cover Photo URL")
	}
	if _, err := url.ParseRequestURI(coverPhotoURL); err != nil {
		return errors.New("The Event Cover Photo URL must be a valid URL")
	}

	return nil
}
