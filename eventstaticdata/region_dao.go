package eventstaticdata

/*
name
type
major
minor
lat
long
radius
*/

type Region struct {
	ID      int32   `json:"id,omitempty"`
	Name    string  `json:"name"`
	Type    string  `json:"type"`
	Major   int32   `json:"major,omitempty"`
	Minor   int32   `json:"minor,omitempty"`
	Lat     float64 `json:"lat,omitempty"`
	Long    float64 `json:"lat,omitempty"`
	Radius  int32   `json:"radius,omitempty"`
	EventID int32   `json:"eventid"`
}
