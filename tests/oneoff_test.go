package tests

import (
	"time"

	. "github.com/smartystreets/assertions"

	"flywheel.io/sdk/api"
)

func (t *F) TestGetConfig() {
	config, _, err := t.GetConfig()
	t.So(err, ShouldBeNil)

	now := time.Now()
	t.So(config.Auth, ShouldNotBeEmpty)
	t.So(config.Site, ShouldNotBeEmpty)
	t.So(config.Created, ShouldHappenBefore, now)
	t.So(config.Modified, ShouldHappenBefore, now)
}

func (t *F) TestGetVersion() {
	version, _, err := t.GetVersion()
	t.So(err, ShouldBeNil)

	// Eh, a vaguely-modern database number
	t.So(version.Database, ShouldBeGreaterThan, 20)
}

func (t *F) TestError() {
	var aerr *api.Error
	var unused interface{}
	resp, err := t.New().Get("does-not-exist").Receive(&unused, &aerr)

	t.So(err, ShouldBeNil)
	t.So(aerr, ShouldNotBeNil)
	t.So(resp.StatusCode, ShouldEqual, 404)
	t.So(aerr.StatusCode, ShouldEqual, 404)
	t.So(aerr.Message, ShouldEqual, "The resource could not be found.")
}
