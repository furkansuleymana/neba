// This file is heavily inspired by @Reve.
// Source: https://github.com/Reve/httpDigestAuth/blob/master/httpDigestAuth.go
package network

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// DigestAuth represents the necessary fields for performing Digest Authentication.
// Digest Authentication is a method used to confirm the identity of a user before
// allowing access to a resource. It involves a challenge-response mechanism where
// the server sends a nonce value and the client responds with a hashed value that
// includes the nonce, username, password, and other details.
type DigestAuth struct {
	realm      string
	qop        string
	nonce      string
	opaque     string
	algorithm  string
	ha1        string
	ha2        string
	cnonce     string
	uri        string
	nonceCount int
	username   string
	password   string
}

// Authenticate performs an HTTP GET request to the specified URI using the provided client.
// If the server responds with a 401 Unauthorized status and a Digest authentication challenge,
// it attempts to handle Digest authentication. If Digest authentication is not required,
// it falls back to Basic Authentication.
//
// Parameters:
//   - username: The username for authentication.
//   - password: The password for authentication.
//   - uri:      The URI to send the request to.
//   - client:   The HTTP client to use for the request.
//
// Returns:
//   - *http.Request: The authenticated HTTP request.
//   - error:         An error if the request creation or execution fails.
func Authenticate(username, password, uri string, client *http.Client) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body) // Discard response body

	if resp.StatusCode == http.StatusUnauthorized {
		authHeader := resp.Header.Get("WWW-Authenticate")
		if strings.HasPrefix(strings.ToLower(authHeader), "digest ") {
			return (&DigestAuth{}).handleDigestAuth(username, password, uri, authHeader)
		}
	}

	// Fallback to Basic Authentication
	req.SetBasicAuth(username, password)
	return req, nil
}

// handleDigestAuth handles the Digest authentication process by parsing the
// authentication parameters from the provided authHeader, initializing the
// DigestAuth struct with the parsed parameters, and creating a new HTTP GET
// request with the appropriate Digest authentication header.
//
// Parameters:
//   - username:   The username for authentication.
//   - password:   The password for authentication.
//   - uri:        The URI for the HTTP request.
//   - authHeader: The WWW-Authenticate header containing the Digest authentication parameters.
//
// Returns:
//   - *http.Request: The HTTP request with the Digest authentication header.
//   - error:         An error if the authentication parameters could not be parsed or the request could not be created.
func (d *DigestAuth) handleDigestAuth(username, password, uri, authHeader string) (*http.Request, error) {
	authParams := parseDigestAuthParams(authHeader)
	if authParams == nil {
		return nil, fmt.Errorf("failed to parse WWW-Authenticate header")
	}

	d.initialize(authParams, username, password, uri)

	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	d.addDigestAuthHeader(req)
	return req, nil
}

// addDigestAuthHeader adds a Digest authentication header to the provided HTTP request.
// It increments the nonce count, generates a new client nonce (cnonce), and computes the
// digest response using the HTTP method and request URI. The resulting header is then
// set in the request's "Authorization" header.
//
// Parameters:
//   - req: The HTTP request to which the Digest authentication header will be added.
func (d *DigestAuth) addDigestAuthHeader(req *http.Request) {
	d.nonceCount++
	d.cnonce = generateRandomKey()
	d.uri = req.URL.RequestURI()
	d.computeDigest(req.Method)

	hash := fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s:%s:%08x:%s:%s:%s", d.ha1, d.nonce, d.nonceCount, d.cnonce, d.qop, d.ha2))))
	authHeader := fmt.Sprintf(
		`Digest username="%s", realm="%s", nonce="%s", uri="%s", cnonce="%s", nc=%08x, qop=%s, response="%s", algorithm=%s`,
		d.username, d.realm, d.nonce, d.uri, d.cnonce, d.nonceCount, d.qop, hash, d.algorithm,
	)
	if d.opaque != "" {
		authHeader += fmt.Sprintf(", opaque=\"%s\"", d.opaque)
	}

	req.Header.Set("Authorization", authHeader)
}

// parseDigestAuthParams parses the Digest authentication parameters from the given header string.
// It expects the header to start with "Digest " (case insensitive) and then a comma-separated list of key-value pairs.
// Each key-value pair is split by an equals sign, and both the key and value are trimmed of surrounding quotes and spaces.
// The function returns a map where the keys are the parameter names and the values are the corresponding parameter values.
// If the header does not start with "Digest ", the function returns nil.
func parseDigestAuthParams(header string) map[string]string {
	if !strings.HasPrefix(strings.ToLower(header), "digest ") {
		return nil
	}

	params := make(map[string]string)
	for _, kv := range strings.Split(header[len("Digest "):], ",") {
		parts := strings.SplitN(strings.TrimSpace(kv), "=", 2)
		if len(parts) == 2 {
			params[strings.Trim(parts[0], "\" ")] = strings.Trim(parts[1], "\" ")
		}
	}
	return params
}

// computeDigest calculates the HA1 and HA2 hash values used in Digest Authentication.
// HA1 is computed as the MD5 hash of the concatenation of the username, realm, and password.
// HA2 is computed as the MD5 hash of the concatenation of the HTTP method and URI.
// The resulting hashes are stored in the DigestAuth struct fields ha1 and ha2.
//
// Parameters:
//   - method: The HTTP method (e.g., "GET", "POST") used in the request.
func (d *DigestAuth) computeDigest(method string) {
	d.ha1 = fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s:%s:%s", d.username, d.realm, d.password))))
	d.ha2 = fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s:%s", method, d.uri))))
}

// generateRandomKey generates a random 12-byte key and returns it as a base64 encoded string.
// If there is an error during the random key generation, it returns an empty string.
func generateRandomKey() string {
	key := make([]byte, 12)
	if _, err := rand.Read(key); err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(key)
}

// initialize sets up the DigestAuth instance with the provided authentication parameters,
// username, password, and URI. It initializes the realm, qop, nonce, opaque, and algorithm
// fields from the authParams map, sets the URI, and initializes the nonce count to 0.
// It also sets the username and password for the DigestAuth instance.
//
// Parameters:
//   - authParams: a map containing authentication parameters such as realm, qop, nonce, opaque, and algorithm
//   - username:   the username for authentication
//   - password:   the password for authentication
//   - uri:        the URI for the request
func (d *DigestAuth) initialize(authParams map[string]string, username, password, uri string) {
	d.realm = authParams["realm"]
	d.qop = authParams["qop"]
	d.nonce = authParams["nonce"]
	d.opaque = authParams["opaque"]
	d.algorithm = strings.ToUpper(authParams["algorithm"])
	d.uri = uri
	d.nonceCount = 0
	d.username = username
	d.password = password
}

// TODO: Add support for session-based authentication.
