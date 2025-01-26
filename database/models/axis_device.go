package models

// AxisDevice represents a device with its attributes.
type AxisDevice struct {
	SerialNumber string `json:"serial_number"`
	Model        string `json:"model"`
	IPAddress    string `json:"ip_address"`
	OSVersion    string `json:"os_version"`
	Username     string `json:"username"`
	Password     string `json:"password"` // This is a bad idea.
}

// TODO: Add a method to validate the device.
