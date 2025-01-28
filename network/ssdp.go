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

// DiscoverSSDP discovers devices on the network using the Simple Service Discovery Protocol (SSDP).
// It sends an M-SEARCH request to the SSDP multicast address and listens for responses from devices.
//
// Returns a slice of device locations (URLs) and an error if any occurred during the discovery process.
//
// The function performs the following steps:
// 1. Creates a UDP connection.
// 2. Resolves the SSDP multicast address.
// 3. Sends the M-SEARCH request to the multicast address.
// 4. Sets a deadline for the connection to avoid indefinite blocking.
// 5. Reads responses from the connection until the deadline is reached.
// 6. Extracts device locations from the responses and returns them.
//
// Returns:
// - []string: A slice of device locations (URLs) discovered on the network.
// - error: An error if any occurred during the discovery process.
func DiscoverSSDP() ([]string, error) {
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
