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

	// Add acquisition to the collection
	_, _, sessionId, acquisitionId := t.createTestAcquisition()
	_, err = t.AddAcquisitionsToCollection(cId, []string{acquisitionId})
	t.So(err, ShouldBeNil)

	// Get Sessions
	savedSessions, _, err := t.GetCollectionSessions(cId)
	t.So(savedSessions, ShouldHaveLength, 1)
	t.So(savedSessions[0].Id, ShouldEqual, sessionId)

	// Get Acquisitions
	savedAcquisitions, _, err := t.GetCollectionAcquisitions(cId)
	t.So(savedAcquisitions, ShouldHaveLength, 1)
	t.So(savedAcquisitions[0].Id, ShouldEqual, acquisitionId)

	// Get Session Acquisitions
	savedSessionAcquisitions, _, err := t.GetCollectionSessionAcquisitions(cId, savedSessions[0].Id)
	t.So(savedSessionAcquisitions, ShouldHaveLength, 1)
	t.So(savedSessionAcquisitions[0].Id, ShouldEqual, acquisitionId)

	// Add session to the collection
	_, _, sessionId, acquisitionId = t.createTestAcquisition()
	_, err = t.AddSessionsToCollection(cId, []string{sessionId})
	t.So(err, ShouldBeNil)

	// Get Sessions
	savedSessions, _, err = t.GetCollectionSessions(cId)
	t.So(savedSessions, ShouldHaveLength, 2)
	// Could add contains check

	// Get Acquisitions
	savedAcquisitions, _, err = t.GetCollectionAcquisitions(cId)
	t.So(savedAcquisitions, ShouldHaveLength, 2)
	// Could add contains check

	// Get Session Acquisitions
	savedSessionAcquisitions, _, err = t.GetCollectionSessionAcquisitions(cId, savedSessions[0].Id)
	t.So(savedSessionAcquisitions, ShouldHaveLength, 1)
	// Could add contains check

	// Get all
	collections, _, err := t.GetAllCollections()
	t.So(err, ShouldBeNil)
	// workaround: all-container endpoints skip some fields, single-container does not. this sets up the equality check
	savedCollection.Files = nil
	savedCollection.Notes = nil
	savedCollection.Info = nil
	// t.So(collections, ShouldContain, savedCollection)

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

	// Add note
	_, err = t.AddCollectionNote(cId, "This is a note")
	t.So(err, ShouldBeNil)
	changedCollection, _, err = t.GetCollection(cId)
	t.So(changedCollection.Notes, ShouldHaveLength, 1)
	t.So(changedCollection.Notes[0].Text, ShouldEqual, "This is a note")

	// Delete
	_, err = t.DeleteCollection(cId)
	t.So(err, ShouldBeNil)
	collections, _, err = t.GetAllCollections()
	t.So(err, ShouldBeNil)
	t.So(collections, ShouldNotContain, savedCollection)
}

func (t *F) TestCollectionFiles() {
	collection := &api.Collection{Name: RandString()}
	collectionId, _, err := t.AddCollection(collection)
	t.So(err, ShouldBeNil)

	src := UploadSourceFromString("yeats.txt", "Surely some revelation is at hand;")
	progress, resultChan := t.UploadToCollection(collectionId, src)

	// Last update should be the full string length.
	t.checkProgressChanEndsWith(progress, 34)
	t.So(<-resultChan, ShouldBeNil)

	rCollection, _, err := t.GetCollection(collectionId)
	t.So(err, ShouldBeNil)
	t.So(rCollection.Files, ShouldHaveLength, 1)
	t.So(rCollection.Files[0].Name, ShouldEqual, "yeats.txt")
	t.So(rCollection.Files[0].Size, ShouldEqual, 34)
	t.So(rCollection.Files[0].Mimetype, ShouldEqual, "text/plain")

	// Download the same file
	buffer, dest := DownloadSourceToBuffer()
	progress, resultChan = t.DownloadFromCollection(collectionId, "yeats.txt", dest)

	// Last update should be the full string length.
	t.checkProgressChanEndsWith(progress, 34)
	t.So(<-resultChan, ShouldBeNil)
	t.So(buffer.String(), ShouldEqual, "Surely some revelation is at hand;")
}
