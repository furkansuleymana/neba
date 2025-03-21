package network

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	"github.com/koron/go-ssdp"
)

const (
	ssdpServiceType    = "urn:axis-com:service:BasicService:1"
	ssdpMaxWaitTimeSec = 5
)

// Root represents the root element of the XML response from a device
type Root struct {
	XMLName     xml.Name `xml:"root"`
	SpecVersion SpecVersion
	Device      Device `xml:"device"`
	URLBase     string `xml:"URLBase"`
}

// SpecVersion represents the version of the XML specification
type SpecVersion struct {
	Major int `xml:"major"`
	Minor int `xml:"minor"`
}

// Device represents the device information in the XML response
type Device struct {
	DeviceType       string      `xml:"deviceType"`
	FriendlyName     string      `xml:"friendlyName"`
	Manufacturer     string      `xml:"manufacturer"`
	ManufacturerURL  string      `xml:"manufacturerURL"`
	ModelDescription string      `xml:"modelDescription"`
	ModelName        string      `xml:"modelName"`
	ModelNumber      string      `xml:"modelNumber"`
	ModelURL         string      `xml:"modelURL"`
	SerialNumber     string      `xml:"serialNumber"`
	UDN              string      `xml:"UDN"`
	ServiceList      ServiceList `xml:"serviceList"`
	PresentationURL  string      `xml:"presentationURL"`
}

// ServiceList represents a list of services provided by the device
type ServiceList struct {
	Services []Service `xml:"service"`
}

// Service represents a single service provided by the device
type Service struct {
	ServiceType string `xml:"serviceType"`
	ServiceId   string `xml:"serviceId"`
	ControlURL  string `xml:"controlURL"`
	EventSubURL string `xml:"eventSubURL"`
	SCPDURL     string `xml:"SCPDURL"`
}

// DiscoverSSDP performs a Simple Service Discovery Protocol (SSDP) search
// to discover devices on the network. It retrieves and parses the XML
// data from the discovered devices to extract relevant information.
//
// Returns:
//   - A slice of maps, where each map contains device information such as
//     "FriendlyName", "ModelName", "SerialNumber", and "PresentationURL".
//   - An error if the SSDP search fails or no devices are found.
//
// Notes:
//   - If fetching or unmarshaling XML data for a device fails, the error
//     is logged, and the device is skipped.
//   - The function requires the `ssdp` package for performing the SSDP search
//     and assumes the existence of a `fetchXMLData` function to retrieve XML
//     data from a given URL.
func DiscoverSSDP() ([]map[string]string, error) {
	// Perform SSDP search for devices matching the specified service type
	ssdpResponses, err := ssdp.Search(ssdpServiceType, ssdpMaxWaitTimeSec, "")
	if err != nil {
		return nil, fmt.Errorf("failed to SSDP search: %w", err)
	}
	if len(ssdpResponses) == 0 {
		return nil, fmt.Errorf("no SSDP devices found")
	}

	// Initialize a slice to store discovered device information
	devices := make([]map[string]string, 0, len(ssdpResponses))
	for _, device := range ssdpResponses {
		// Fetch XML data from the device's location URL
		xmlData, err := fetchXMLData(device.Location)
		if err != nil {
			fmt.Printf("failed to fetch XML data from %s: %v\n", device.Location, err)
			continue // Skip this device if fetching XML data fails
		}

		// Unmarshal the XML data into the Root struct
		var root Root
		err = xml.Unmarshal(xmlData, &root)
		if err != nil {
			fmt.Printf("failed to unmarshal XML data: %v\n", err)
			continue // Skip this device if unmarshaling fails
		}

		// Create a map with relevant device information
		deviceInfo := map[string]string{
			"FriendlyName":    root.Device.FriendlyName,
			"ModelName":       root.Device.ModelName,
			"SerialNumber":    root.Device.SerialNumber,
			"PresentationURL": root.Device.PresentationURL,
		}

		// Add the device information to the list of devices
		devices = append(devices, deviceInfo)
	}

	return devices, nil
}

func fetchXMLData(urlStr string) ([]byte, error) {
	resp, err := http.Get(urlStr) // TODO: Add context for cancellation and timeouts
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch URL: status code %d", resp.StatusCode)
	}

	xmlData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return xmlData, nil
}
