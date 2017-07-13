package tests

import (
	"time"

	. "github.com/smartystreets/assertions"

	"flywheel.io/sdk/api"
)

func (t *F) TestSessions() {
	_, projectId := t.createTestProject()

	sessionName := RandString()
	session := &api.Session{
		Name:      sessionName,
		ProjectId: projectId,
		Info: map[string]interface{}{
			"some-key": 37,
		},
		Subject: &api.Subject{
			Code:      RandStringLower(),
			Firstname: RandString(),
			Lastname:  RandString(),
			Sex:       "other",
			Age:       56,
			Info: map[string]interface{}{
				"some-subject-key": 37,
			},
		},
	}

	// Add
	sessionId, _, err := t.AddSession(session)
	t.So(err, ShouldBeNil)

	// Get
	rSession, _, err := t.GetSession(sessionId)
	t.So(err, ShouldBeNil)
	t.So(rSession.Id, ShouldEqual, sessionId)
	t.So(rSession.Name, ShouldEqual, session.Name)
	now := time.Now()
	t.So(rSession.Info, ShouldContainKey, "some-key")
	t.So(rSession.Info["some-key"], ShouldEqual, 37)
	t.So(*rSession.Created, ShouldHappenBefore, now)
	t.So(*rSession.Modified, ShouldHappenBefore, now)
	t.So(*rSession.Subject, ShouldNotBeNil)
	t.So(rSession.Subject.Id, ShouldNotBeEmpty)
	t.So(rSession.Subject.Firstname, ShouldResemble, session.Subject.Firstname)

	// Get all
	sessions, _, err := t.GetAllSessions()
	t.So(err, ShouldBeNil)
	// workaround: all-container endpoints skip some fields, single-container does not. this sets up the equality check
	rSession.Files = nil
	rSession.Notes = nil
	rSession.Tags = nil
	rSession.Info = nil
	rSession.Subject = &api.Subject{
		Id:   rSession.Subject.Id,
		Code: rSession.Subject.Code,
	}
	t.So(sessions, ShouldContain, rSession)

	// Get from parent
	sessions, _, err = t.GetProjectSessions(projectId)
	t.So(err, ShouldBeNil)
	t.So(sessions, ShouldContain, rSession)

	// Modify
	newName := RandString()
	sessionMod := &api.Session{
		Name: newName,
	}
	_, err = t.ModifySession(sessionId, sessionMod)
	t.So(err, ShouldBeNil)
	changedSession, _, err := t.GetSession(sessionId)
	t.So(changedSession.Name, ShouldEqual, newName)
	t.So(*changedSession.Created, ShouldBeSameTimeAs, *rSession.Created)
	t.So(*changedSession.Modified, ShouldHappenAfter, *rSession.Modified)

	// Notes, tags
	message := "This is a note"
	_, err = t.AddSessionNote(sessionId, message)
	t.So(err, ShouldBeNil)
	tag := "example-tag"
	_, err = t.AddSessionTag(sessionId, tag)
	t.So(err, ShouldBeNil)

	// Check
	rSession, _, err = t.GetSession(sessionId)
	t.So(err, ShouldBeNil)
	t.So(rSession.Notes, ShouldHaveLength, 1)
	t.So(rSession.Notes[0].Text, ShouldEqual, message)
	t.So(rSession.Tags, ShouldHaveLength, 1)
	t.So(rSession.Tags[0], ShouldEqual, tag)

	// Delete
	_, err = t.DeleteSession(sessionId)
	t.So(err, ShouldBeNil)
	sessions, _, err = t.GetAllSessions()
	t.So(err, ShouldBeNil)
	t.So(sessions, ShouldNotContain, rSession)
}

func (t *F) TestSessionFiles() {
	_, projectId := t.createTestProject()
	session := &api.Session{Name: RandString(), ProjectId: projectId}
	sessionId, _, err := t.AddSession(session)
	t.So(err, ShouldBeNil)

	src := UploadSourceFromString("yeats.txt", "Are full of passionate intensity.")
	progress, resultChan := t.UploadToSession(sessionId, src)

	// Last update should be the full string length.
	t.checkProgressChanEndsWith(progress, 33)
	t.So(<-resultChan, ShouldBeNil)

	rSession, _, err := t.GetSession(sessionId)
	t.So(err, ShouldBeNil)
	t.So(rSession.Files, ShouldHaveLength, 1)
	t.So(rSession.Files[0].Name, ShouldEqual, "yeats.txt")
	t.So(rSession.Files[0].Size, ShouldEqual, 33)
	t.So(rSession.Files[0].Mimetype, ShouldEqual, "text/plain")

	// Download the same file
	buffer, dest := DownloadSourceToBuffer()
	progress, resultChan = t.DownloadFromSession(sessionId, "yeats.txt", dest)

	// Last update should be the full string length.
	t.checkProgressChanEndsWith(progress, 33)
	t.So(<-resultChan, ShouldBeNil)
	t.So(buffer.String(), ShouldEqual, "Are full of passionate intensity.")
}

func (t *F) createTestSession() (string, string, string) {
	groupId, projectId := t.createTestProject()

	sessionName := RandString()
	session := &api.Session{
		Name:      sessionName,
		ProjectId: projectId,
	}
	sessionId, _, err := t.AddSession(session)
	t.So(err, ShouldBeNil)

	return groupId, projectId, sessionId
}
