package network

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

const (
	mSearchMessage = `M-SEARCH * HTTP/1.1
Host:239.255.255.250:1900
ST:urn:axis-com:service:BasicService:1
Man:"ssdp:discover"
MX:3

`
	ssdpMulticastAddress = "239.255.255.250:1900"
	ssdpTimeout          = 5 * time.Second
	bufferSize           = 2048
	expectedDevices      = 50
)

func DiscoverSSDP(ctx context.Context) ([]string, error) {
	devices := make([]string, 0, expectedDevices)
	var mu sync.Mutex

	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("list network interfaces: %w", err)
	}

	var wg sync.WaitGroup
	errCh := make(chan error, len(interfaces))

	for _, iface := range interfaces {
		if (iface.Flags&net.FlagUp) == 0 || (iface.Flags&net.FlagLoopback) != 0 {
			continue
		}

		wg.Add(1)
		go func(iface net.Interface) {
			defer wg.Done()
			if err := scanInterface(ctx, iface, &devices, &mu); err != nil {
				errCh <- fmt.Errorf("interface %s: %w", iface.Name, err)
			}
		}(iface)
	}

	wg.Wait()
	close(errCh)

	if err := <-errCh; err != nil {
		return devices, err
	}

	return devices, nil
}

func scanInterface(ctx context.Context, iface net.Interface, devices *[]string, mu *sync.Mutex) error {
	addrs, err := iface.Addrs()
	if err != nil {
		return fmt.Errorf("get addresses: %w", err)
	}

	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if !ok || ipNet.IP.To4() == nil {
			continue
		}

		if err := scanAddress(ctx, ipNet.IP.String(), devices, mu); err != nil {
			return err
		}
	}
	return nil
}

func scanAddress(ctx context.Context, ip string, devices *[]string, mu *sync.Mutex) error {
	conn, err := net.ListenPacket("udp4", ip+":0")
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}
	defer conn.Close()

	udpAddr, err := net.ResolveUDPAddr("udp4", ssdpMulticastAddress)
	if err != nil {
		return fmt.Errorf("resolve multicast address: %w", err)
	}

	if _, err := conn.WriteTo([]byte(mSearchMessage), udpAddr); err != nil {
		return fmt.Errorf("send discovery message: %w", err)
	}

	deadline := time.Now().Add(ssdpTimeout)
	if err := conn.SetDeadline(deadline); err != nil {
		return fmt.Errorf("set deadline: %w", err)
	}

	buffer := make([]byte, bufferSize)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			n, _, err := conn.ReadFrom(buffer)
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					return nil
				}
				return fmt.Errorf("read response: %w", err)
			}

			if location := extractDeviceLocation(string(buffer[:n])); location != "" {
				mu.Lock()
				*devices = append(*devices, location)
				mu.Unlock()
			}
		}
	}
}

func extractDeviceLocation(response string) string {
	const locationPrefix = "location:"
	for _, line := range strings.Split(response, "\r\n") {
		if l := strings.ToLower(line); strings.Contains(l, locationPrefix) {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1])
			}
		}
	}
	return ""
}
