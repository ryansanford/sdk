package tests

import (
	. "github.com/smartystreets/assertions"

	"flywheel.io/sdk/api"
)

func (t *F) TestCreateUploadSourceFromFilenames() {
	sources := api.CreateUploadSourceFromFilenames("one.txt", "two.txt")

	t.So(sources[0].Path, ShouldEqual, "one.txt")
	t.So(sources[1].Path, ShouldEqual, "two.txt")
}

func (t *F) TestBadUploads() {
	// Invalid upload source
	source := &api.UploadSource{}
	_, result := t.UploadSimple("", nil, source)
	t.So((<-result).Error(), ShouldEqual, "Neither file name nor path was set in upload source")

	// Nonexistant upload path
	source = &api.UploadSource{Path: "/dev/null/does-not-exist"}
	_, result = t.UploadSimple("", nil, source)
	t.So((<-result).Error(), ShouldStartWith, "open /dev/null/does-not-exist: ")

	// Bad upload url
	source = UploadSourceFromString("yeats.txt", "Surely some revelation is at hand;")
	_, result = t.UploadSimple("not-an-endpoint", nil, source)
	// Could improve this in the future
	t.So((<-result).Error(), ShouldEqual, "{\"status_code\": 404, \"message\": \"The resource could not be found.\"}")
}
