package tests

import (
	"time"

	. "github.com/smartystreets/assertions"

	"flywheel.io/sdk/api"
)

func (t *F) TestProjects() {
	groupId := t.createTestGroup()

	projectName := RandString()
	project := &api.Project{
		Name:        projectName,
		GroupId:     groupId,
		Description: "This is a description",
		Info: map[string]interface{}{
			"some-key": 37,
		},
	}

	// Add
	projectId, _, err := t.AddProject(project)
	t.So(err, ShouldBeNil)

	// Get
	rProject, _, err := t.GetProject(projectId)
	t.So(err, ShouldBeNil)
	t.So(rProject.Id, ShouldEqual, projectId)
	t.So(rProject.Name, ShouldEqual, project.Name)
	t.So(rProject.Description, ShouldEqual, project.Description)
	t.So(rProject.Info, ShouldContainKey, "some-key")
	t.So(rProject.Info["some-key"], ShouldEqual, 37)
	now := time.Now()
	t.So(*rProject.Created, ShouldHappenBefore, now)
	t.So(*rProject.Modified, ShouldHappenBefore, now)

	// Get all
	projects, _, err := t.GetAllProjects()
	t.So(err, ShouldBeNil)
	// workaround: all-container endpoints skip some fields, single-container does not. this sets up the equality check
	rProject.Files = nil
	rProject.Notes = nil
	rProject.Tags = nil
	rProject.Info = nil
	t.So(projects, ShouldContain, rProject)

	// Modify
	newName := RandString()
	projectMod := &api.Project{
		Name: newName,
		Info: map[string]interface{}{
			"another-key": 52,
		},
	}
	_, err = t.ModifyProject(projectId, projectMod)
	t.So(err, ShouldBeNil)
	changedProject, _, err := t.GetProject(projectId)
	t.So(changedProject.Name, ShouldEqual, newName)
	t.So(changedProject.Info, ShouldContainKey, "some-key")
	t.So(changedProject.Info, ShouldContainKey, "another-key")
	t.So(changedProject.Info["another-key"], ShouldEqual, 52)
	t.So(*changedProject.Created, ShouldBeSameTimeAs, *rProject.Created)
	t.So(*changedProject.Modified, ShouldHappenAfter, *rProject.Modified)

	// Notes, tags
	message := "This is a note"
	_, err = t.AddProjectNote(projectId, message)
	t.So(err, ShouldBeNil)
	tag := "example-tag"
	_, err = t.AddProjectTag(projectId, tag)
	t.So(err, ShouldBeNil)

	// Check
	rProject, _, err = t.GetProject(projectId)
	t.So(err, ShouldBeNil)
	t.So(rProject.Notes, ShouldHaveLength, 1)
	t.So(rProject.Notes[0].Text, ShouldEqual, message)
	t.So(rProject.Tags, ShouldHaveLength, 1)
	t.So(rProject.Tags[0], ShouldEqual, tag)

	// Delete
	_, err = t.DeleteProject(projectId)
	t.So(err, ShouldBeNil)
	projects, _, err = t.GetAllProjects()
	t.So(err, ShouldBeNil)
	t.So(projects, ShouldNotContain, rProject)
}

func (t *F) TestProjectFiles() {
	groupId := t.createTestGroup()

	project := &api.Project{Name: RandString(), GroupId: groupId}
	projectId, _, err := t.AddProject(project)
	t.So(err, ShouldBeNil)

	src := UploadSourceFromString("yeats.txt", "The best lack all conviction, while the worst")
	progress, resultChan := t.UploadToProject(projectId, src)

	// Last update should be the full string length.
	t.checkProgressChanEndsWith(progress, 45)
	t.So(<-resultChan, ShouldBeNil)

	rProject, _, err := t.GetProject(projectId)
	t.So(err, ShouldBeNil)
	t.So(rProject.Files, ShouldHaveLength, 1)
	t.So(rProject.Files[0].Name, ShouldEqual, "yeats.txt")
	t.So(rProject.Files[0].Size, ShouldEqual, 45)
	t.So(rProject.Files[0].Mimetype, ShouldEqual, "text/plain")

	// Download the same file
	buffer, dest := DownloadSourceToBuffer()
	progress, resultChan = t.DownloadFromProject(projectId, "yeats.txt", dest)

	// Last update should be the full string length.
	t.checkProgressChanEndsWith(progress, 45)
	t.So(<-resultChan, ShouldBeNil)
	t.So(buffer.String(), ShouldEqual, "The best lack all conviction, while the worst")
}

func (t *F) createTestProject() (string, string) {
	groupId := t.createTestGroup()

	project := &api.Project{
		Name:    RandString(),
		GroupId: groupId,
	}
	projectId, _, err := t.AddProject(project)
	t.So(err, ShouldBeNil)

	return groupId, projectId
}
