package tests

import (
	"bytes"

	. "github.com/smartystreets/assertions"

	"flywheel.io/sdk/api"
)

func (t *F) TestClient() {

	invalidKey := func() {
		api.NewApiKeyClient("bad-format")
	}

	invalidKeyWithBadPort := func() {
		api.NewApiKeyClient("hostname.example:not-a-port:my-key")
	}

	validKeyWithNoPort := func() {
		api.NewApiKeyClient("hostname.example:my-key")
	}

	validKeyWithPort := func() {
		api.NewApiKeyClient("hostname.example:443:my-key")
	}

	validKeyWithInsecure := func() {
		api.NewApiKeyClient("hostname.example:80:my-key", api.InsecureNoSSLVerification)
	}

	validKeyWithPlaintext := func() {
		api.NewApiKeyClient("hostname.example:80:my-key", api.InsecureUsePlaintext)
	}

	t.So(invalidKey, ShouldPanic)
	t.So(invalidKeyWithBadPort, ShouldPanic)

	t.So(validKeyWithNoPort, ShouldNotPanic)
	t.So(validKeyWithPort, ShouldNotPanic)
	t.So(validKeyWithInsecure, ShouldNotPanic)
	t.So(validKeyWithPlaintext, ShouldNotPanic)
}

func (t *F) TestDebugTransport() {
	// This could test response content if we exposed the key or allowed regeneration of clients
	var buffer bytes.Buffer
	client := api.NewApiKeyClient("hostname.example:80:my-key", api.DebugLogRequests(&buffer))

	_, _, err := client.GetCurrentUser()
	t.So(err, ShouldNotBeNil)

	message := buffer.String()
	t.So(message, ShouldContainSubstring, "BEGIN HTTP")
	t.So(message, ShouldContainSubstring, "BEGIN RESPONSE")
	t.So(message, ShouldContainSubstring, "Request failed: ")
	t.So(message, ShouldContainSubstring, "END HTTP")
	t.So(message, ShouldContainSubstring, "GET /api")
	t.So(message, ShouldContainSubstring, "Host: ")
	t.So(message, ShouldContainSubstring, "User-Agent: ")
	t.So(message, ShouldContainSubstring, "Authorization: ")
}
