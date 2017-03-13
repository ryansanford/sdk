package tests

import (
	"os"
	"sync"
	"testing"

	// . "github.com/smartystreets/assertions"
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
	> which uses t.Parallel() under the hood for every test case.
	> - @mdwhatcott
	> https://github.com/smartystreets/goconvey/issues/360

	Downside: we have to put up with gunit's struct model. It doesn't look too bad, and might turn out to be helpful.


	For goal #4, my plan is to incorportate go-vcr:
	https://github.com/dnaeon/go-vcr

	The implementation throws requests in YAML files, which... eh, let's try it maybe.
	There will have to be some setup trickery to transparently hit live or recorded.
	I think the vcr transport should handle that.

*/

// TestSuite fires off gunit.
//
// Gunit will look at various function name prefixes to determine behaviour:
//
//   "Test": Well, it's a test.
//   "Skip": Skipped.
//   "Long": Skipped when `go test` is ran with the `short` flag.
//
//   "Setup":    Executed before each test.
//   "Teardown": Executed after  each test.
//
// Functions without
func TestSuite(t *testing.T) {
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

	// SdkTestHost is the environment variable that sets the test host.
	// Valid values are a host:port combination: "localhost:8443".
	// No affect in unit test mode.
	SdkTestHost = "SdkTestHost"

	// SdkTestKey is the environment variable that sets the test API key.
	// Valid values are an API key: "32334"
	// No affect in unit test mode.
	SdkTestKey = "SdkTestKey"

	DefaultMode = "integration"
	DefaultHost = "localhost:8443"
	DefaultKey  = "change-me"
)

// makeClient reads settings from the environment and returns the corresponding client
func makeClient() *api.Client {
	mode, modeSet := os.LookupEnv(SdkTestMode)

	if !modeSet {
		mode = DefaultMode
	}

	if mode != "integration" && mode != "unit" {
		panic("Unsupported test mode " + mode)
	}

	if mode == "unit" {
		panic("Unit test mode is not supported yet")
	}

	if mode == "integration" {
		host, hostSet := os.LookupEnv(SdkTestHost)
		key, keySet := os.LookupEnv(SdkTestKey)

		if !hostSet {
			host = DefaultHost
		}

		if !keySet {
			key = DefaultKey
		}

		return api.NewApiKeyClient(host, key, true)
	}

	return nil
}

// Re-use state: clients are safe for concurrent use and are stateless.
var once sync.Once
var client *api.Client

// Setup prepares the fixture with SDK client state. Runs once per test.
func (t *F) Setup() {
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
