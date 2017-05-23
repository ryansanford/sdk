package api

import (
	"crypto/tls"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"

	"github.com/dghubble/sling"
)

// NewApiKeyClient creates a Client with the given API key and options.
// Passing a key with an invalid format will panic.
func NewApiKeyClient(apiKey string, options ...ApiKeyClientOption) *Client {

	// If the debug environment variable is set, add the debug transport.
	// This is added to the beginning of the options array, in case it is overridden later.
	_, debug := os.LookupEnv(SdkDebugKey)
	if debug {
		options = append([]ApiKeyClientOption{DebugLogRequests(os.Stderr)}, options...)
	}

	// Load all configuration options into a config struct
	config := DefaultApiKeyClientOptions
	for _, x := range options {
		x(&config)
	}

	host, port, key, err := ParseApiKey(apiKey)
	if err != nil {
		panic(err)
	}

	// Load TLS configuration into a transport
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: config.InsecureSkipVerify,
		},
	}

	var hc *http.Client

	// Create an HTTP client, adding the debug transport if specified
	if config.DebugWriter == nil {
		hc = &http.Client{
			Transport: tr,
		}
	} else {
		debugT := &DebugTransport{
			Transport: tr,
			Writer:    config.DebugWriter,
		}
		hc = debugT.Client()
	}

	protocol := "https"
	if config.InsecureUsePlaintext {
		protocol = "http"
	}

	// Create a sling client, which is used for most server interactions
	sc := sling.New().
		Base(protocol+"://"+host+":"+strconv.Itoa(port)+"/").
		Set("Authorization", "scitran-user "+key).
		Path("api/").
		Client(hc)

	return &Client{
		hc,
		sc,
	}
}

// DebugTransport prints its raw request bodies to a writer.
type DebugTransport struct {

	// All requests made with this transport will be written to Writer.
	// It will default to os.Stderr if nil.
	Writer io.Writer

	// Transport is the underlying HTTP transport to use when making requests.
	// It will default to http.DefaultTransport if nil.
	Transport http.RoundTripper
}

var header = []byte("\n--------------------\n---- BEGIN HTTP ----\n--------------------\n")
var footer = []byte("\n--------------------\n----- END HTTP -----\n--------------------\n")

// RoundTrip implements the RoundTripper interface.
func (t *DebugTransport) RoundTrip(req *http.Request) (*http.Response, error) {

	msg, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		return nil, err
	}

	msg = append(header, msg...)
	msg = append(msg, footer...)

	if t.Writer == nil {
		_, err = os.Stderr.Write(msg)
	} else {
		_, err = t.Writer.Write(msg)
	}

	if err != nil {
		return nil, err
	}

	if t.Transport != nil {
		return t.Transport.RoundTrip(req)
	} else {
		return http.DefaultTransport.RoundTrip(req)
	}
}

// Client returns an *http.Client which uses the DebugTransport.
func (t *DebugTransport) Client() *http.Client {
	return &http.Client{Transport: t}
}
