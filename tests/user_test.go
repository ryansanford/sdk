package tests

import (
	"time"

	. "github.com/smartystreets/assertions"

	"flywheel.io/sdk/api"
)

func (t *F) sanityCheckUser(user *api.User) {
	t.So(user.Id, ShouldNotBeEmpty)
	t.So(user.Email, ShouldContainSubstring, "@")
	t.So(user.Firstname, ShouldNotBeEmpty)
	t.So(user.Lastname, ShouldNotBeEmpty)
	now := time.Now()
	t.So(*user.Created, ShouldHappenBefore, now)
	t.So(*user.Modified, ShouldHappenBefore, now)
}

func (t *F) TestGetCurrentUser() {
	user, _, err := t.GetCurrentUser()
	t.So(err, ShouldBeNil)

	t.sanityCheckUser(user)
	t.So(user.ApiKey.Key, ShouldNotBeEmpty)

	now := time.Now()
	t.So(*user.ApiKey.Created, ShouldHappenBefore, now)
	t.So(*user.ApiKey.LastUsed, ShouldHappenBefore, now)
}

func (t *F) TestGetAllUsers() {
	users, _, err := t.GetAllUsers()
	t.So(err, ShouldBeNil)

	t.So(len(users), ShouldBeGreaterThan, 0)

	for _, user := range users {
		t.sanityCheckUser(user)
	}
}

func (t *F) TestGetUser() {
	user, _, err := t.GetCurrentUser()
	t.So(err, ShouldBeNil)

	user2, _, err := t.GetUser(user.Id)
	t.So(err, ShouldBeNil)

	// Should be a valid user, but not have an api key
	t.sanityCheckUser(user2)
	t.So(user2.ApiKey, ShouldBeNil)
}

func (t *F) TestAddModifyDeleteUser() {
	email := RandString() + "@" + RandString() + ".com"
	user := &api.User{
		Id:        email,
		Email:     email,
		Firstname: RandString(),
		Lastname:  RandString(),
	}

	// Add
	userId, _, err := t.AddUser(user)
	t.So(err, ShouldBeNil)
	t.So(userId, ShouldNotBeEmpty)

	// modify
	newName := RandString()
	modUser := &api.User{
		Firstname: newName,
	}
	_, err = t.ModifyUser(userId, modUser)
	t.So(err, ShouldBeNil)

	//Check
	compare, _, err := t.GetUser(userId)
	t.So(err, ShouldBeNil)
	t.So(compare.Id, ShouldEqual, user.Id)
	t.So(compare.Email, ShouldEqual, user.Email)
	t.So(compare.Firstname, ShouldEqual, newName)
	t.So(compare.Lastname, ShouldEqual, user.Lastname)

	// remove
	_, err = t.DeleteUser(userId)
	t.So(err, ShouldBeNil)
}
