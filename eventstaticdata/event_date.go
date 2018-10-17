package eventstaticdata

import (
	"fmt"
	"time"
)

type EventDate struct {
	time.Time
}

// To be consistent with the time package
// (see https://stackoverflow.com/a/50557983/3837124)
const RFC3339FullDate = "2006-01-02"

func (ed *EventDate) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		ed.Time = time.Time{}
		return nil
	}
	time, err := time.Parse(RFC3339FullDate, string(b))
	if err != nil {
		return err
	}
	ed.Time = time
	return nil
}

func (ed *EventDate) MarshalJSON() ([]byte, error) {
	if ed.IsZero() {
		return make([]byte, 0), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", ed.Format(RFC3339FullDate))), nil
}
