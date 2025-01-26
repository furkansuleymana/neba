package api

import (
	"fmt"
	"io"
	"net/http"

	"github.com/furkansuleymana/neba/tools"
)

var DigestAuth *tools.DigestAuth

func Param(ip string, username string, password string) error {
	url := fmt.Sprintf("http://%s/axis-cgi/param.cgi", ip)

	httpClient := http.Client{}
	req, _ := http.NewRequest("GET", url, nil)

	DigestAuth, err := DigestAuth.Authenticate(username, password, url, &httpClient)
	if err != nil {
		req.SetBasicAuth(username, password)
		resp, _ := httpClient.Do(req)
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("reading response body: %v", err)
		}
		fmt.Println(string(body))
		return nil
	} // TODO: DRY

	DigestAuth.AddAuthHeader(req)

	resp, _ := httpClient.Do(req)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %v", err)
	}
	fmt.Println(string(body))

	return nil
}
