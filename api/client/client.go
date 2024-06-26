// Package client defines a Mercury API client.
package client

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/M-ERCURY/core/api/apiversion"
	"github.com/M-ERCURY/core/api/auth"
	"github.com/M-ERCURY/core/api/nonce"
	"github.com/M-ERCURY/core/api/signer"
	"github.com/M-ERCURY/core/api/status"
)

// Client is an API client type. It exposes the http.Client interface and it is
// safe for concurrent use by multiple goroutines. It also encapsulates the
// keypair used to access the API via its signer.Signer interface.
type Client struct {
	*http.Client
	signer.Signer
	component string

	do func(*http.Request) (*http.Response, error)
}

// New creates a new API client using the given signer to sign API requests.
func New(s signer.Signer, component string) *Client {
	return &Client{&http.Client{}, s, component, nil}
}

// NewMock creates a new API client which uses a given handler for handling all
// requests.
func NewMock(s signer.Signer, h http.Handler, component string) (c *Client) {
	c = New(s, component)

	c.do = func(r *http.Request) (*http.Response, error) {
		rw := httptest.NewRecorder()
		h.ServeHTTP(rw, r)
		return rw.Result(), nil
	}

	return
}

func (c *Client) Do(r *http.Request) (*http.Response, error) {
	if c.do != nil {
		return c.do(r)
	}

	return c.Client.Do(r)
}

func (c *Client) SetComponent(s string)            { c.component = s }
func (c *Client) SetTransport(t http.RoundTripper) { c.Client.Transport = t }

// NewRequest is a convenience function for creating a new http.Request with a
// payload that's JSON-marshaled and signed.
func (c *Client) NewRequest(method string, url string, data interface{}) (*http.Request, error) {
	b, err := json.Marshal(data)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(b))

	if err != nil {
		return nil, err
	}

	err = Refresh(req)

	if err != nil {
		return nil, err
	}

	auth.SetHeader(req.Header, auth.API, auth.Version, apiversion.VERSION_STRING)
	req.Header.Set("Content-Type", "application/json")

	if c.Signer != nil {
		sig := c.Sign(b)

		pks := base64.RawURLEncoding.EncodeToString(c.Public())
		sgs := base64.RawURLEncoding.EncodeToString(sig)

		auth.SetHeader(req.Header, c.component, auth.Pubkey, pks)
		auth.SetHeader(req.Header, c.component, auth.Signature, sgs)
	}

	return req, nil
}

func (c *Client) PerformRequestNoParse(req *http.Request, cs ...string) (res *http.Response, err error) {
	var body io.ReadCloser

	// this should never fail as we're setting io.Reader as the body in
	// NewRequest which provides a GetBody(). but just in case...
	body, err = req.GetBody()

	if err != nil {
		err = fmt.Errorf("could not refresh request body for retry: %w", err)
		return
	}

	defer func() { req.Body = body }()
	return c.Do(req)
}

func (c *Client) PerformRequestOnce(req *http.Request, out interface{}, cs ...string) (err error) {
	var res *http.Response
	res, err = c.PerformRequestNoParse(req, cs...)

	if err != nil {
		return
	}

	return ParseResponse(res, out, cs...)
}

// Perform is a convenience function for creating a new request, performing it
// and parsing the JSON response into the receiving interface.
func (c *Client) Perform(method string, url string, in interface{}, out interface{}, cs ...string) (err error) {
	req, err := c.NewRequest(method, url, in)

	if err != nil {
		return
	}

	tries := 3
	interval := 5 * time.Second

	for i := 1; i <= tries; i++ {
		err = c.PerformRequestOnce(req, out, cs...)

		respErr := &status.T{}
		if errors.As(err, &respErr) {
			// TODO fast solution for now. fix later
			if strings.Contains(respErr.Cause.Error(), string(status.CauseSneakyPof)) {
				return status.SneakyPofErr
			}
		}

		if err == nil || i == tries || !status.IsRetryable(err) {
			// success or max retries hit or no-retry error; return nil or last error
			break
		}

		log.Printf(
			"client: error performing %s %s: %s on try %d of %d, retrying in %s...",
			method,
			url,
			err,
			i,
			tries,
			interval,
		)

		time.Sleep(interval)
	}

	return
}

// ParseResponse extracts the JSON-encoded payload of a request response and
// checks for API errors. It is not a method of the Client type since it uses
// no Client-specific data. Therefore, while low-level, it can be called by
// code which does not use Client if needed.
func ParseResponse(res *http.Response, dst interface{}, cs ...string) (err error) {
	defer res.Body.Close()
	var body []byte

	if len(cs) > 0 {
		body, err = auth.SignedResBody(res, cs...)
	} else {
		body, err = ioutil.ReadAll(res.Body)
	}

	fmt.Println("")
	fmt.Println("RESPONSE BODY:", string(body))
	fmt.Println("")

	if err != nil {
		return fmt.Errorf(
			"error while trying to read response body: %w",
			err,
		)
	}

	// check API version if present
	if len(auth.GetHeader(res.Header, auth.API, auth.Version)) > 0 {
		err = auth.VersionCheck(res.Header, auth.API, &apiversion.VERSION)

		if err != nil {
			return fmt.Errorf(
				"invalid server API version: %w, request body: %s",
				err,
				string(body),
			)
		}
	}

	// check for API error
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		e := &status.T{}
		err = json.Unmarshal(body, e)

		if err != nil {
			return fmt.Errorf(
				"error while trying to parse response body `%s`: %w",
				string(body),
				err,
			)
		}

		return e
	}

	if dst != nil {
		err = json.Unmarshal(body, dst)

		if err != nil {
			return fmt.Errorf(
				"error while trying to unmarshal response body to destination: %w, body='%s'",
				err,
				string(body),
			)
		}
	}

	return nil
}

// Refresh refreshes the idempotency key (if applicable).
func Refresh(req *http.Request) (err error) {
	if req.Method == http.MethodPost {
		var ik string
		ik, err = nonce.New(32)

		if err != nil {
			return
		}

		req.Header.Set("Idempotency-Key", ik)
	}

	return
}
