package api

import (
	"crypto/tls"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
	"strings"

	"github.com/dghubble/sling"
)

// Client is an http and sling client capable of making flywheel requests.
type Client struct {
	*http.Client
	*sling.Sling
}

type ApiKeyClientOption func(*ApiKeyClientOptions)
type ApiKeyClientOptions struct {
	// Skip SSL verification
	InsecureSkipVerify bool

	// Use plaintext (HTTP) transport
	InsecureUsePlaintext bool
}

var DefaultApiKeyClientOptions = ApiKeyClientOptions{
	InsecureSkipVerify:   false,
	InsecureUsePlaintext: false,
}

// Specify that the ApiKeyClient should not verify SSL connections.
// Should only be used for development.
var InsecureNoSSLVerification ApiKeyClientOption

// Specify that the ApiKeyClient should use a plaintext HTTP transport.
// Should only be used for development.
var InsecureUsePlaintext ApiKeyClientOption

func init() {
	InsecureNoSSLVerification = func(o *ApiKeyClientOptions) {
		o.InsecureSkipVerify = true
	}

	InsecureUsePlaintext = func(o *ApiKeyClientOptions) {
		o.InsecureUsePlaintext = true
	}
}

// NewApiKeyClient creates a Client with the given API key and options.
// Passing a key with an invalid format will panic.
func NewApiKeyClient(apiKey string, options ...ApiKeyClientOption) *Client {
	config := DefaultApiKeyClientOptions
	for _, x := range options {
		x(&config)
	}

	splits := strings.Split(apiKey, ":")
	if len(splits) < 2 {
		panic("Invalid API key")
	}

	var err error
	host := ""
	port := 443
	key := ""

	if len(splits) == 2 {
		host = splits[0]
		key = splits[1]
	} else {
		host = splits[0]
		port, err = strconv.Atoi(splits[1])
		key = splits[len(splits)-1]
	}

	if err != nil {
		panic(err)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: config.InsecureSkipVerify,
		},
	}

	hc := &http.Client{
		Transport: tr,
	}

	protocol := "https"
	if config.InsecureUsePlaintext {
		protocol = "http"
	}

	sc := sling.New().
		Base(protocol+"://"+host+":"+strconv.Itoa(port)+"/").
		Set("Authorization", "scitran-user "+key).
		Path("api/").
		Client(hc)

	client := &Client{
		hc,
		sc,
	}

	return client
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

var header = []byte("\n--------------------\n-----BEGIN HTTP-----\n--------------------\n")
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
