package network

import (
	"fmt"
	"net"
	"strings"
	"time"
)

const (
	// Configuration constants
	mSearchMessage = `
M-SEARCH * HTTP/1.1
Host:239.255.255.250:1900
ST:urn:axis-com:service:BasicService:1
Man:"ssdp:discover"
MX:3

`
	ssdpMulticastAddress = "239.255.255.250:1900"
	ssdpTimeout          = 5 * time.Second
	bufferSize           = 2048
)

// DiscoverWithSSDP sends an SSDP M-SEARCH request and returns
// a list of discovered device locations.
func DiscoverWithSSDP() ([]string, error) {
	// Create a UDP connection
	conn, err := net.ListenPacket("udp4", ":0")
	if err != nil {
		return nil, fmt.Errorf("create UDP connection: %w", err)
	}
	defer conn.Close()

	// Resolve the SSDP multicast address
	addr, err := net.ResolveUDPAddr("udp4", ssdpMulticastAddress)
	if err != nil {
		return nil, fmt.Errorf("resolve SSDP address: %w", err)
	}

	// Send the M-SEARCH request
	_, err = conn.WriteTo([]byte(mSearchMessage), addr)
	if err != nil {
		return nil, fmt.Errorf("send M-SEARCH request: %w", err)
	}

	// Set a deadline for the connection
	if err := conn.SetDeadline(time.Now().Add(ssdpTimeout)); err != nil {
		return nil, fmt.Errorf("set connection timeout: %w", err)
	}

	var devices []string
	buffer := make([]byte, bufferSize)

	// Read responses from the connection
	for {
		n, _, err := conn.ReadFrom(buffer)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				break
			}
			return nil, fmt.Errorf("read response: %w", err)
		}

		response := string(buffer[:n])

		// Extract the device location from the response
		location := extractDeviceLocation(response)
		if location != "" {
			devices = append(devices, location)
		}
	}

	return devices, nil
}

// extractDeviceLocation extracts the LOCATION header from an SSDP response.
func extractDeviceLocation(response string) string {
	for _, line := range strings.Split(response, "\r\n") {
		if strings.HasPrefix(line, "LOCATION:") {
			return strings.TrimSpace(line)
		}
	}
	return ""
}
