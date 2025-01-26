package tools

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type DigestAuth struct {
	realm      string
	qop        string
	method     string
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

// Authenticate performs Digest Authentication for the given username and password
// against the specified URI using the provided HTTP client.
//
// It sends an initial request to the URI to retrieve the WWW-Authenticate header
// and then parses the digest authentication parameters from the response.
//
// If the response status code is 401 (Unauthorized), it extracts the necessary
// authentication parameters and returns a DigestAuth instance populated with
// these parameters. If the status code is 200 (OK), it indicates that
// authentication is not required and returns nil.
//
// Parameters:
//   - username: The username for authentication.
//   - password: The password for authentication.
//   - uri: The URI to authenticate against.
//   - client: The HTTP client to use for making requests.
//
// Returns:
//   - A pointer to a DigestAuth instance populated with the authentication parameters.
//   - An error if the request fails, the response status code is unexpected, or
//     the WWW-Authenticate header cannot be parsed.
func (d *DigestAuth) Authenticate(username, password, uri string, client *http.Client) (*DigestAuth, error) {
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body) // Ensure the response body is read and discarded

	if resp.StatusCode != http.StatusUnauthorized {
		if resp.StatusCode == http.StatusOK {
			return nil, nil // Auth not required
		}
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	authParams := parseDigestAuthParams(resp)
	if authParams == nil {
		return nil, fmt.Errorf("failed to parse WWW-Authenticate header")
	}

	return &DigestAuth{
		realm:      authParams["realm"],
		qop:        authParams["qop"],
		nonce:      authParams["nonce"],
		opaque:     authParams["opaque"],
		algorithm:  strings.ToUpper(authParams["algorithm"]),
		nonceCount: 0,
		username:   username,
		password:   password,
		uri:        uri,
	}, nil
}

// AddAuthHeader adds the Digest authentication header to the provided HTTP request.
// It increments the nonce count, generates a new client nonce, sets the HTTP method and URI,
// computes the digest response, and constructs the Authorization header.
//
// The Authorization header is formatted as per the Digest authentication scheme, including
// the username, realm, nonce, URI, client nonce, nonce count, quality of protection (qop),
// response, and algorithm. If the opaque value is present, it is also included in the header.
//
// Parameters:
//
//	req (*http.Request): The HTTP request to which the Authorization header will be added.
func (d *DigestAuth) AddAuthHeader(req *http.Request) {
	d.nonceCount++
	d.cnonce = generateRandomKey()
	d.method = req.Method
	d.uri = req.URL.RequestURI()
	d.computeDigest()

	response := fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s:%s:%08x:%s:%s:%s", d.ha1, d.nonce, d.nonceCount, d.cnonce, d.qop, d.ha2))))

	authHeader := fmt.Sprintf(
		`Digest username="%s", realm="%s", nonce="%s", uri="%s", cnonce="%s", nc=%08x, qop=%s, response="%s", algorithm=%s`,
		d.username, d.realm, d.nonce, d.uri, d.cnonce, d.nonceCount, d.qop, response, d.algorithm,
	)

	if d.opaque != "" {
		authHeader += fmt.Sprintf(", opaque=\"%s\"", d.opaque)
	}

	req.Header.Set("Authorization", authHeader)
}

// parseDigestAuthParams extracts and parses the parameters from a Digest
// authentication challenge in an HTTP response header. It returns a map
// containing the key-value pairs of the parameters.
//
// Parameters:
//   - resp: A pointer to an http.Response object containing the HTTP response.
//
// Returns:
//   - A map[string]string containing the parsed Digest authentication parameters.
//     If the "WWW-Authenticate" header does not start with "Digest ", it returns nil.
func parseDigestAuthParams(resp *http.Response) map[string]string {
	header := resp.Header.Get("WWW-Authenticate")
	if !strings.HasPrefix(strings.ToLower(header), "digest ") {
		return nil
	}

	params := map[string]string{}
	for _, kv := range strings.Split(header[len("Digest "):], ",") {
		parts := strings.SplitN(strings.TrimSpace(kv), "=", 2)
		if len(parts) == 2 {
			params[strings.Trim(parts[0], "\" ")] = strings.Trim(parts[1], "\" ")
		}
	}
	return params
}

// computeDigest calculates the HA1 and HA2 hashes used in Digest Authentication.
// HA1 is computed using the username, realm, and password.
// HA2 is computed using the HTTP method and URI.
// The resulting hashes are stored in the DigestAuth struct.
func (d *DigestAuth) computeDigest() {
	ha1Input := fmt.Sprintf("%s:%s:%s", d.username, d.realm, d.password)
	ha1Hash := md5.Sum([]byte(ha1Input))
	d.ha1 = fmt.Sprintf("%x", ha1Hash)

	ha2Input := fmt.Sprintf("%s:%s", d.method, d.uri)
	ha2Hash := md5.Sum([]byte(ha2Input))
	d.ha2 = fmt.Sprintf("%x", ha2Hash)
}

// generateRandomKey generates a random key of 12 bytes and returns it as a base64 encoded string.
// It uses the crypto/rand package to generate cryptographically secure random bytes.
// If there is an error during the random byte generation, it returns an empty string.
//
// Returns:
//   - A base64 encoded string representation of the random key.
//   - An empty string if there is an error during random byte generation.
func generateRandomKey() string {
	key := make([]byte, 12)
	if _, err := rand.Read(key); err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(key)
}
