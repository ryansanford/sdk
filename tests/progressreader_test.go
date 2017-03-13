package tests

import (
	"bytes"
	"io"
	"io/ioutil"

	. "github.com/smartystreets/assertions"

	"flywheel.io/sdk/api"
)

func (t *F) checkProgressChanEndsWith(progress chan int64, total int64) {
	// ProgressReader contract is to only send updates on change.
	last := int64(0)
	for x := range progress {
		t.So(x, ShouldBeGreaterThan, last)
		last = x
	}

	// Last update should be the full string length.
	t.So(last, ShouldEqual, total)
}

func (t *F) TestProgressReaderSimple() {

	dataStr := "Things fall apart; the centre cannot hold;"

	// Create a reader that closes (for coverage's sake)
	src := bytes.NewBufferString(dataStr)
	srcCloser := ioutil.NopCloser(src)

	progress := make(chan int64, 10)
	pr := api.NewProgressReader(srcCloser, progress)
	dest := &bytes.Buffer{}

	_, err1 := io.Copy(dest, pr)
	err2 := pr.Close()
	t.So(err1, ShouldBeNil)
	t.So(err2, ShouldBeNil)

	// Last update should be the full string length.
	t.checkProgressChanEndsWith(progress, 42)

	t.So(dest.String(), ShouldEqual, dataStr)
}

func (t *F) TestProgressReaderAdvanced() {

	dataStr := "Mere anarchy is loosed upon the world,"

	src1 := bytes.NewBufferString(dataStr)
	progress := make(chan int64, 10)
	pr := api.NewProgressReader(src1, progress)
	dest := &bytes.Buffer{}

	// Copy the first two words, then wait so a progress report will generate.
	written, err := io.CopyN(dest, pr, 13)
	t.So(err, ShouldBeNil)
	t.So(written, ShouldEqual, 13)

	// First progress update should be of the initial bytes.
	x := <-progress
	t.So(x, ShouldEqual, 13)

	// Set a new reader from the same string
	src2 := bytes.NewBufferString(dataStr)
	pr.SetReader(src2)

	// Copy the rest of the bytes
	_, err1 := io.Copy(dest, pr)
	err2 := pr.Close()
	t.So(err1, ShouldBeNil)
	t.So(err2, ShouldBeNil)

	// Last update should be the full string length.
	t.checkProgressChanEndsWith(progress, 51)

	// Result should have copied the first 13 bytes
	t.So(dest.String(), ShouldEqual, "Mere anarchy Mere anarchy is loosed upon the world,")

	// Initial buffer should have some content left; last buffer should have none
	t.So(src1.Len(), ShouldEqual, 25)
	t.So(src2.Len(), ShouldEqual, 0)
}
