package tests

import (
	"time"

	. "github.com/smartystreets/assertions"

	"flywheel.io/sdk/api"
)

// Separate context creation out from the test
func (t *F) contextTestSessions() (string, string) {
	groupId := RandStringLower()
	_, _, err := t.AddGroup(&api.Group{Id: groupId})
	t.So(err, ShouldBeNil)

	projectName := RandString()
	project := &api.Project{
		Name:    projectName,
		GroupId: groupId,
	}
	projectId, _, err := t.AddProject(project)
	t.So(err, ShouldBeNil)

	return groupId, projectId
}

func (t *F) TestSessions() {
	_, projectId := t.contextTestSessions()

	sessionName := RandString()
	session := &api.Session{
		Name: sessionName,
		ProjectId: projectId,
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
	t.So(*rSession.Created, ShouldHappenBefore, now)
	t.So(*rSession.Modified, ShouldHappenBefore, now)

	// Get all
	sessions, _, err := t.GetAllSessions()
	t.So(err, ShouldBeNil)
	rSession.Files = nil // workaround: all-container endpoints skip files array, single-container does not. this sets up the equality check
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

	// Delete
	_, err = t.DeleteSession(sessionId)
	t.So(err, ShouldBeNil)
	sessions, _, err = t.GetAllSessions()
	t.So(err, ShouldBeNil)
	t.So(sessions, ShouldNotContain, rSession)
}

func (t *F) TestSessionUpload() {
	_, projectId := t.contextTestSessions()
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
	t.So(rSession.Files[0].Name, ShouldEqual, "yeats.txt")
	t.So(rSession.Files[0].Size, ShouldEqual, 33)
	t.So(rSession.Files[0].Mimetype, ShouldEqual, "text/plain")
}
