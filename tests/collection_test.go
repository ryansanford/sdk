package tests

import (
	"time"

	. "github.com/smartystreets/assertions"

	"flywheel.io/sdk/api"
)

func (t *F) TestCollections() {
	collectionName := RandString()

	collection := &api.Collection{
		Name:        collectionName,
		Description: RandString(),
	}

	// Add
	cId, _, err := t.AddCollection(collection)
	t.So(err, ShouldBeNil)

	// Get
	savedCollection, _, err := t.GetCollection(cId)
	t.So(err, ShouldBeNil)
	t.So(savedCollection.Id, ShouldEqual, cId)
	t.So(savedCollection.Name, ShouldEqual, collection.Name)
	now := time.Now()
	t.So(*savedCollection.Created, ShouldHappenBefore, now)
	t.So(*savedCollection.Modified, ShouldHappenBefore, now)

	// Get all
	collections, _, err := t.GetAllCollections()
	t.So(err, ShouldBeNil)
	// workaround: all-container endpoints skip some fields, single-container does not. this sets up the equality check
	savedCollection.Files = nil
	savedCollection.Notes = nil
	savedCollection.Tags = nil
	savedCollection.Info = nil
	t.So(collections, ShouldContain, savedCollection)

	// Modify
	newName := RandString()
	collectionMod := &api.Collection{
		Name: newName,
	}
	_, err = t.ModifyCollection(cId, collectionMod)
	t.So(err, ShouldBeNil)

	// Check
	changedCollection, _, err := t.GetCollection(cId)
	t.So(changedCollection.Name, ShouldEqual, newName)
	t.So(*changedCollection.Created, ShouldBeSameTimeAs, *savedCollection.Created)
	t.So(*changedCollection.Modified, ShouldHappenAfter, *savedCollection.Modified)

	// Delete
	_, err = t.DeleteCollection(cId)
	t.So(err, ShouldBeNil)
	collections, _, err = t.GetAllCollections()
	t.So(err, ShouldBeNil)
	t.So(collections, ShouldNotContain, savedCollection)
}
