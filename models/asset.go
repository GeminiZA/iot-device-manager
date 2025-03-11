package models

type Telemetry map[string]interface{}

type Asset struct {
	DeviceID  string
	Name      string
	Status    string
	Telemetry Telemetry
}

// Methods
