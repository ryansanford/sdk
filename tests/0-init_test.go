package tests

import (
	"os"
	"sync"
	"testing"

	. "github.com/smartystreets/assertions"
	"github.com/smartystreets/gunit"

	"flywheel.io/sdk/api"
)

/*
	Testing goals, in order:

	1) automatically parallel on the terminal
	2) zero-ish overhead
	3) tolerable syntax
	4) able to both hit a testing infra, or replay locally.


	For goal #1, Gunit seems to be the best-equipped to handle things, and there's some agreement on that:

	> I know that I'm the creator of GoConvey and all, but I've actually moved to gunit,
	> which uses t.Parallel() under the hood for every test case. - @mdwhatcott
	> https://github.com/smartystreets/goconvey/issues/360

	For goal #4, my plan is to incorporate go-vcr:
	https://github.com/dnaeon/go-vcr

	The implementation throws requests in YAML files, which... eh, let's try it maybe.
	There will have to be some setup trickery to transparently hit live or recorded.
	I think the vcr transport should handle that.


	Test requirements:

	1) Each test works independent of any preexisting state, or lack thereof. Only a working, root API key is required.
	2) Ideally tests can clean up after themselves, but this is not required.

	Please keep these goals and requirements in mind when modifying this package.


	There are instances where context creation could/should be handed off to a struct setup - right now, those are instead handled by "context" functions. This isn't a perfect layout if we had a larger test suite, but seems to work fine for now. Let's leave it until & unless it becomes unbearable.
*/

// TestSuite fires off gunit.
//
// Gunit will look at various function name prefixes to determine behavior:
//
//   "Test": Well, it's a test.
//   "Skip": Skipped.
//   "Long": Skipped when `go test` is ran with the `short` flag.
//
//   "Setup":    Executed before each test.
//   "Teardown": Executed after  each test.
//
// Functions without these prefixes are ignored.
func TestSuite(t *testing.T) {

	// Workaround:
	//
	// In CI, the client can lag for awhile waiting for the API to become available.
	// This inflates the timing of the first GO_MAXPROCS unit tests, making them look artificially slow.
	//
	// In an attempt to combat this, let's query to the API once, ignoring results and any errors.
	// Hopefully, this moves timing to the suite and out of the individual tests.
	workAroundClient := makeClient()
	// begin := time.Now()
	workAroundClient.GetCurrentUser()
	// duration := time.Since(begin)
	// PrintFormat("Connecting to the API took " + duration.String())

	gunit.Run(new(F), t)
}

// F is the default fixture, so-named for convenience.
type F struct {
	*gunit.Fixture

	*api.Client
}

const (
	// SdkTestMode is the environment variable that sets the test mode.
	// Valid values are "unit" and "integration".
	SdkTestMode = "SdkTestMode"

	// SdkTestKey is the environment variable that sets the test API key.
	// Valid values are an API key: "localhost:8443:change-me"
	// No affect in unit test mode.
	SdkTestKey = "SdkTestKey"

	// SdkTestKey is the environment variable that sets the API protocol.
	// Valid values are "https" and "http".
	// No affect in unit test mode.
	SdkProtocolKey = "SdkTestProtocol"

	DefaultMode     = "integration"
	DefaultKey      = "localhost:8443:change-me"
	DefaultProtocol = "https"
)

// makeClient reads settings from the environment and returns the corresponding client
func makeClient() *api.Client {
	mode, modeSet := os.LookupEnv(SdkTestMode)

	if !modeSet {
		mode = DefaultMode
	}

	var client *api.Client

	if mode == "unit" {
		panic("Unit test mode is not supported yet")

	} else if mode == "integration" {
		key, keySet := os.LookupEnv(SdkTestKey)
		protocol, protocolSet := os.LookupEnv(SdkProtocolKey)

		if !keySet {
			key = DefaultKey
		}

		if !protocolSet {
			protocol = DefaultProtocol
		}

		if protocol == "https" {
			client = api.NewApiKeyClient(key, api.InsecureNoSSLVerification)
		} else if protocol == "http" {
			client = api.NewApiKeyClient(key, api.InsecureNoSSLVerification, api.InsecureUsePlaintext)
		} else {
			panic("Protocol must be http or https, was " + protocol)
		}

	} else {
		panic("Unsupported test mode " + mode)
	}

	return client
}

// Re-use state: clients are safe for concurrent use and are stateless.
var once sync.Once
var client *api.Client

// Setup prepares the fixture with SDK client state. Runs once per test.
func (t *F) Setup() {

	// Use custom fork of gunit to specify which assertions should result in the test failing immediately.
	// This prevents predictable, useless stack traces from trying to access bad data during an assertion.
	t.AddFatalAssertion(ShouldBeNil)
	t.AddFatalAssertion(ShouldNotBeNil)
	t.AddFatalAssertion(ShouldHaveLength)

	once.Do(func() {
		client = makeClient()
	})

	t.Client = client
}

/*
// An example test:
func (t *F) SkipTestExample() {
	t.So(42, ShouldEqual, 42)
	t.So("Hello, World!", ShouldContainSubstring, "World")
}
*/
