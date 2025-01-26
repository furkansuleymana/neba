package api

import (
	"fmt"
	"io"
	"net/http"

	"github.com/furkansuleymana/neba/tools"
)

func Param(ip string, username string, password string) error {
	url := fmt.Sprintf("http://%s/axis-cgi/param.cgi?action=list", ip)

	client := &http.Client{}
	req, err := tools.Authenticate(username, password, url, client)
	if err != nil {
		return fmt.Errorf("authentication failed: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %v", err)
	}
	fmt.Println(string(body))

	return nil
}
