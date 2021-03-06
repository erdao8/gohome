package belkin

// DeviceAttributes contains values returned from the FetchAttributes call. Values are pointers, a nil value
// indicated we didn't get any value for this particular attribute
type DeviceAttributes struct {
	// The current state of the switch, 1 -> on, 0 -> off
	Switch *int

	// The sensor value, 1 -> open, 0 -> closed
	Sensor *int

	// The mode of the switch, 0 -> toggle, 1 -> momentary
	SwitchMode *int

	// If the sensor is active, 1 -> yes, 0 -> no
	SensorPresent *int
}
