package tests

import (
	"errors"
	"io/ioutil"

	. "github.com/smartystreets/assertions"

	"flywheel.io/sdk/api"
)

func (t *F) TestCoalesce() {
	// No error
	res := api.Coalesce(nil, nil)
	t.So(res, ShouldBeNil)

	// Http error
	err := errors.New("This is an error")
	res = api.Coalesce(err, nil)
	t.So(res, ShouldEqual, err)

	// Both http error and api error
	aErr := &api.Error{Message: "This is an api error", StatusCode: 500}
	res = api.Coalesce(err, aErr)
	t.So(res, ShouldEqual, err)

	// Api error
	res = api.Coalesce(nil, aErr)
	t.So(res.Error(), ShouldEqual, "(500) This is an api error")

	// Invalid api error
	aErr = &api.Error{Message: "", StatusCode: 500}
	res = api.Coalesce(nil, aErr)
	t.So(res.Error(), ShouldEqual, "(500) Unknown server error")
}

func (t *F) TestFormat() {
	wat := api.Format(&api.File{Name: "yeats.txt"})
	t.So(wat, ShouldEqual, "{\n\t\"name\": \"yeats.txt\"\n}")

	t.So(func() {
		api.Format(ioutil.NopCloser)
	}, ShouldPanic)
}
