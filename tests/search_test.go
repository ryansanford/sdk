package tests

import (
	"time"

	. "github.com/smartystreets/assertions"

	"flywheel.io/sdk/api"
)

// Ref https://github.com/flywheel-io/sdk/issues/31
func (t *F) SkipTestSearch() {
	_, _, sessionId, acquisitionId := t.createTestAcquisition()

	// Get Acquisition
	a, _, err := t.GetAcquisition(acquisitionId)
	t.So(err, ShouldBeNil)
	// Ref https://github.com/flywheel-io/sdk/issues/32
	time.Sleep(1000 * time.Millisecond)

	s := &api.SearchQuery{
		ReturnType:   api.SessionString,
		SearchString: a.Name,
	}

	sR, _, err := t.SearchRaw(s)
	// t.So(a.Name, ShouldBeNil)
	t.So(err, ShouldBeNil)
	t.So(len(sR.Results), ShouldEqual, 1)
	t.So(sR.Results[0].Source.Session.Id, ShouldEqual, sessionId)

	s = &api.SearchQuery{
		ReturnType:   api.AcquisitionString,
		SearchString: a.Name,
	}

	sR, _, err = t.SearchRaw(s)
	// t.So(a.Name, ShouldBeNil)
	t.So(err, ShouldBeNil)
	t.So(len(sR.Results), ShouldEqual, 1)
	t.So(sR.Results[0].Source.Acquisition.Id, ShouldEqual, acquisitionId)


	s = &api.SearchQuery{
		ReturnType:   api.SessionString,
		SearchString: a.Name,
	}

	sC, _, err := t.Search(s)
	// t.So(a.Name, ShouldBeNil)
	t.So(err, ShouldBeNil)
	t.So(len(sC), ShouldEqual, 1)
	t.So(sC[0].Session.Id, ShouldEqual, sessionId)

	s = &api.SearchQuery{
		ReturnType:   api.AcquisitionString,
		SearchString: a.Name,
	}

	sC, _, err = t.Search(s)
	// t.So(a.Name, ShouldBeNil)
	t.So(err, ShouldBeNil)
	t.So(len(sC), ShouldEqual, 1)
	t.So(sC[0].Acquisition.Id, ShouldEqual, acquisitionId)
}
