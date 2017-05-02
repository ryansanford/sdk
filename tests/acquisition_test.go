package tests

import (
	"time"

	. "github.com/smartystreets/assertions"

	"flywheel.io/sdk/api"
)

func (t *F) TestAcquisitions() {
	_, _, sessionId := t.createTestSession()

	acquisitionName := RandString()
	acquisition := &api.Acquisition{
		Name:      acquisitionName,
		SessionId: sessionId,
	}

	// Add
	acquisitionId, _, err := t.AddAcquisition(acquisition)
	t.So(err, ShouldBeNil)

	// Get
	rAcquisition, _, err := t.GetAcquisition(acquisitionId)
	t.So(err, ShouldBeNil)
	t.So(rAcquisition.Id, ShouldEqual, acquisitionId)
	t.So(rAcquisition.Name, ShouldEqual, acquisition.Name)
	now := time.Now()
	t.So(*rAcquisition.Created, ShouldHappenBefore, now)
	t.So(*rAcquisition.Modified, ShouldHappenBefore, now)

	// Get all
	acquisitions, _, err := t.GetAllAcquisitions()
	t.So(err, ShouldBeNil)
	// workaround: all-container endpoints skip some fields, single-container does not. this sets up the equality check
	rAcquisition.Files = nil
	rAcquisition.Notes = nil
	rAcquisition.Tags = nil
	rAcquisition.Info = nil
	t.So(acquisitions, ShouldContain, rAcquisition)

	// Modify
	newName := RandString()
	acquisitionMod := &api.Acquisition{
		Name: newName,
	}
	_, err = t.ModifyAcquisition(acquisitionId, acquisitionMod)
	t.So(err, ShouldBeNil)
	changedAcquisition, _, err := t.GetAcquisition(acquisitionId)
	t.So(changedAcquisition.Name, ShouldEqual, newName)
	t.So(*changedAcquisition.Created, ShouldBeSameTimeAs, *rAcquisition.Created)
	t.So(*changedAcquisition.Modified, ShouldHappenAfter, *rAcquisition.Modified)

	// Notes, tags
	message := "This is a note"
	_, err = t.AddAcquisitionNote(acquisitionId, message)
	t.So(err, ShouldBeNil)
	tag := "example-tag"
	_, err = t.AddAcquisitionTag(acquisitionId, tag)
	t.So(err, ShouldBeNil)

	// Check
	rAcquisition, _, err = t.GetAcquisition(acquisitionId)
	t.So(err, ShouldBeNil)
	t.So(rAcquisition.Notes, ShouldHaveLength, 1)
	t.So(rAcquisition.Notes[0].Text, ShouldEqual, message)
	t.So(rAcquisition.Tags, ShouldHaveLength, 1)
	t.So(rAcquisition.Tags[0], ShouldEqual, tag)

	// Delete
	_, err = t.DeleteAcquisition(acquisitionId)
	t.So(err, ShouldBeNil)
	acquisitions, _, err = t.GetAllAcquisitions()
	t.So(err, ShouldBeNil)
	t.So(acquisitions, ShouldNotContain, rAcquisition)
}

func (t *F) SkipTestAcquisitionUpload() {
	_, _, sessionId := t.createTestSession()

	acquisition := &api.Acquisition{Name: RandString(), SessionId: sessionId}
	acquisitionId, _, err := t.AddAcquisition(acquisition)
	t.So(err, ShouldBeNil)

	src := UploadSourceFromString("yeats.txt", "Things fall apart; the centre cannot hold;")
	progress, resultChan := t.UploadToAcquisition(acquisitionId, src)

	// Last update should be the full string length.
	t.checkProgressChanEndsWith(progress, 42)
	t.So(<-resultChan, ShouldBeNil)

	rAcquisition, _, err := t.GetAcquisition(acquisitionId)
	t.So(err, ShouldBeNil)
	t.So(rAcquisition.Files, ShouldHaveLength, 1)
	t.So(rAcquisition.Files[0].Name, ShouldEqual, "yeats.txt")
	t.So(rAcquisition.Files[0].Size, ShouldEqual, 42)
	t.So(rAcquisition.Files[0].Mimetype, ShouldEqual, "text/plain")
}

func (t *F) createTestAcquisition() (string, string, string, string) {
	groupId, projectId, sessionId := t.createTestSession()

	acquisitionName := RandString()
	acquisition := &api.Acquisition{
		Name:      acquisitionName,
		SessionId: sessionId,
	}
	acquisitionId, _, err := t.AddAcquisition(acquisition)
	t.So(err, ShouldBeNil)

	return groupId, projectId, sessionId, acquisitionId
}
