package network

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/koron/go-ssdp"
)

const (
	ssdpServiceType    = "urn:axis-com:service:BasicService:1"
	ssdpMaxWaitTimeSec = 5
)

var ssdpHTTPClient = &http.Client{
	Timeout: 5 * time.Second,
}

type Root struct {
	XMLName xml.Name `xml:"root"`
	Device  Device   `xml:"device"`
	URLBase string   `xml:"URLBase"`
}

type Device struct {
	SerialNumber    string `xml:"serialNumber"`
	PresentationURL string `xml:"presentationURL"`
}

// DiscoverSSDP performs a Simple Service Discovery Protocol (SSDP) search to discover devices on the network.
// It returns a slice of maps, where each map contains information about a discovered device, such as its
// serial number and IP address.
//
// Returns:
// - []map[string]string: A slice of maps containing device information.
// - error: An error if the SSDP search fails or no devices are found.
//
// Example device information map:
//
//	{
//	    "SerialNumber": "123456789",
//	    "IPAddress": "192.168.1.100",
//	}
//
// Errors:
// - If the SSDP search fails, an error is returned.
// - If no SSDP devices are found, an error is returned.
func DiscoverSSDP() ([]map[string]string, error) {
	ssdpResponses, err := ssdp.Search(ssdpServiceType, ssdpMaxWaitTimeSec, "")
	if err != nil {
		return nil, fmt.Errorf("failed to SSDP search: %w", err)
	}
	if len(ssdpResponses) == 0 {
		return nil, fmt.Errorf("no SSDP devices found")
	}

	devices := make([]map[string]string, 0, len(ssdpResponses))
	for _, device := range ssdpResponses {
		xmlData, err := fetchXMLData(device.Location)
		if err != nil {
			fmt.Printf("Failed to fetch XML data from %s: %v\n", device.Location, err)
			continue
		}

		root, err := unmarshalXML(xmlData)
		if err != nil {
			fmt.Printf("Failed to unmarshal XML data: %v\n", err)
			continue
		}

		ipAddress := extractIPAddress(root.URLBase)

		deviceInfo := map[string]string{
			"SerialNumber": root.Device.SerialNumber,
			"IPAddress":    ipAddress,
		}

		devices = append(devices, deviceInfo)
	}

	return devices, nil
}

func fetchXMLData(urlStr string) ([]byte, error) {
	httpResponse, err := ssdpHTTPClient.Get(urlStr)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch URL: status code %d", httpResponse.StatusCode)
	}

	return io.ReadAll(httpResponse.Body)
}

func unmarshalXML(data []byte) (*Root, error) {
	var root Root
	err := xml.Unmarshal(data, &root)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling XML: %w", err)
	}
	return &root, nil
}

func extractIPAddress(urlStr string) string {
	parsedURL, err := url.Parse(urlStr)
	if err == nil {
		hostParts := strings.Split(parsedURL.Host, ":")
		return hostParts[0]
	}

	schemeOffset := strings.Index(urlStr, "http://") + len("http://")
	hostEndPos := strings.IndexAny(urlStr[schemeOffset:], ":/")
	if hostEndPos == -1 {
		return urlStr[schemeOffset:]
	}
	return urlStr[schemeOffset : schemeOffset+hostEndPos]
}
