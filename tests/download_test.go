package tests

import (
	. "github.com/smartystreets/assertions"

	"flywheel.io/sdk/api"
)

func (t *F) TestCreateDownloadSourceFromFilenames() {
	source := api.CreateDownloadSourceFromFilename("one.txt")

	t.So(source.Path, ShouldEqual, "one.txt")
}

func (t *F) TestBadDownloads() {
	// Invalid download source
	source := &api.DownloadSource{}
	_, result := t.DownloadSimple("", source)
	t.So((<-result).Error(), ShouldEqual, "Neither destination path nor writer was set in download source")

	// Nonexistant download path
	source = &api.DownloadSource{Path: "/dev/null/does-not-exist"}
	_, result = t.DownloadSimple("", source)
	t.So((<-result).Error(), ShouldStartWith, "open /dev/null/does-not-exist: ")

	// Bad download url
	buffer, source := DownloadSourceToBuffer()
	_, result = t.DownloadSimple("not-an-endpoint", source)

	// Could improve this in the future
	err := <-result
	t.So(err.Error(), ShouldEqual, "{\"status_code\": 404, \"message\": \"The resource could not be found.\"}")
	t.So(buffer.String(), ShouldEqual, "")
}

// Given an download function, container ID, filename, and content - download & check content
func (t *F) downloadText(fn func(string, string, *api.DownloadSource) (chan int64, chan error), id, filename, text string) {
	buffer, dest := DownloadSourceToBuffer()
	progress, resultChan := fn(id, filename, dest)

	// Last update should be the full string length.
	t.checkProgressChanEndsWith(progress, int64(len(text)))
	t.So(<-resultChan, ShouldBeNil)
	t.So(buffer.String(), ShouldEqual, text)
}
