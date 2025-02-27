package handlers

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func RegisterStreamRoute(mux *http.ServeMux) {
	mux.HandleFunc("/api/stream", handleStream)
}

func handleStream(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	streamMJPEG(w, req, "http://192.168.61.207/axis-cgi/mjpg/video.cgi")
}

func streamMJPEG(w http.ResponseWriter, req *http.Request, mjpegURL string) error {
	w.Header().Set("Content-Type", "multipart/x-mixed-replace; boundary=mjpegboundary")
	w.Header().Set("Cache-Control", "no-cache")

	resp, err := (&http.Client{Timeout: 10 * time.Second}).Get(mjpegURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		return err
	}
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}

	for {
		if req.Context().Err() != nil {
			return nil
		}

		// Skip to boundary
		_, err = reader.ReadSlice('-')
		if err != nil {
			return err
		}

		_, err = reader.ReadString('\n')
		if err != nil {
			return err
		}

		// Read headers until empty line
		headerData := ""
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
			if line == "\r\n" || line == "\n" {
				break
			}
			headerData += line
		}

		// Get content length
		contentLength := 0
		for _, line := range strings.Split(headerData, "\n") {
			if strings.HasPrefix(strings.ToLower(line), "content-length:") {
				fmt.Sscanf(strings.TrimSpace(line[15:]), "%d", &contentLength)
				break
			}
		}

		if contentLength <= 0 {
			continue
		}

		// Read frame
		frameData := make([]byte, contentLength)
		if _, err = io.ReadFull(reader, frameData); err != nil {
			return err
		}

		// Write frame
		fmt.Fprintf(w, "--mjpegboundary\r\nContent-Type: image/jpeg\r\nContent-Length: %d\r\n\r\n", len(frameData))
		w.Write(frameData)
		fmt.Fprint(w, "\r\n")

		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}

		// Read to end of part
		_, err = reader.ReadBytes('\n')
		if err != nil {
			return err
		}
	}
}
