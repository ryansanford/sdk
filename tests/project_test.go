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
		Name:    projectName,
		GroupId: groupId,
	}

	// Add
	projectId, _, err := t.AddProject(project)
	t.So(err, ShouldBeNil)

	// Get
	rProject, _, err := t.GetProject(projectId)
	t.So(err, ShouldBeNil)
	t.So(rProject.Id, ShouldEqual, projectId)
	t.So(rProject.Name, ShouldEqual, project.Name)
	now := time.Now()
	t.So(*rProject.Created, ShouldHappenBefore, now)
	t.So(*rProject.Modified, ShouldHappenBefore, now)

	// Get all
	projects, _, err := t.GetAllProjects()
	t.So(err, ShouldBeNil)
	rProject.Files = nil // workaround: all-container endpoints skip files array, single-container does not. this sets up the equality check
	t.So(projects, ShouldContain, rProject)

	// Modify
	newName := RandString()
	projectMod := &api.Project{
		Name: newName,
	}
	_, err = t.ModifyProject(projectId, projectMod)
	t.So(err, ShouldBeNil)
	changedProject, _, err := t.GetProject(projectId)
	t.So(changedProject.Name, ShouldEqual, newName)
	t.So(*changedProject.Created, ShouldBeSameTimeAs, *rProject.Created)
	t.So(*changedProject.Modified, ShouldHappenAfter, *rProject.Modified)

	// Delete
	_, err = t.DeleteProject(projectId)
	t.So(err, ShouldBeNil)
	projects, _, err = t.GetAllProjects()
	t.So(err, ShouldBeNil)
	t.So(projects, ShouldNotContain, rProject)
}

func (t *F) TestProjectUpload() {
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
	t.So(rProject.Files[0].Name, ShouldEqual, "yeats.txt")
	t.So(rProject.Files[0].Size, ShouldEqual, 45)
	t.So(rProject.Files[0].Mimetype, ShouldEqual, "text/plain")
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
