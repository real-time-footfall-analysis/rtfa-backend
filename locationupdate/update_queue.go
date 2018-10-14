package locationupdate

type BluetoothUpdate struct {
	UUID        string `json:"uuid"`
	EventID     string `json:"EventId"`
	BeaconMajor int    `json:"BeaconMajor"`
	BeaconMinor int    `json:"BeaconMinor"`
	Entering    bool   `json:"Entering"`
}

// initConn opens the connection to the location event queue
func initConn() {
	//TODO: make connection to queue service
}

// Pre: the event object is valid
func addLocationUpdate(event *BluetoothUpdate) error {

	// TODO: add event to queue
	return nil

}
