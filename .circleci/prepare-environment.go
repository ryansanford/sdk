package main

import (
	"encoding/json"
	. "log"
	"time"

	"gopkg.in/mgo.v2"

	"flywheel.io/sdk/api"
)

func main() {
	var session *mgo.Session
	var err error

	for i := 1; i <= 15; i++ {
		err = nil
		Println("Connecting to mongo...")
		wait := time.Duration(float64(i) * 0.3 * float64(time.Second))
		session, err = mgo.DialWithTimeout("localhost", wait)
		if err == nil { break }
	}

	if err != nil { Fatalln("Could not connect to mongo:", err) }
	defer session.Close()
	session.SetSafe(&mgo.Safe{})

	root := true
	now := time.Now()
	testUser := &api.User{
		Id: "a@b.c", Email: "a@b.c",
		Firstname: "Test", Lastname: "User",
		Created: &now, Modified: &now,
		RootAccess: &root,
		ApiKey: &api.Key{
			Key: "insecure-key",
			Created: &now, LastUsed: &now,
		},
	}

	// The mongo client does not seem to respect json annotations!
	// Sidestep this by passing it through a json encode, back to a map.
	// Could open an issue on mgo later.
	raw, _ := json.Marshal(testUser)
	var encoded map[string]interface{}
	json.Unmarshal(raw, &encoded)

	err = session.DB("scitran").C("users").Insert(encoded)
	if err != nil { Fatalln(err) }
	Println("Test user inserted.")

	Println("Connecting to API...")
	client := api.NewApiKeyClient("localhost:8080:insecure-key", api.InsecureNoSSLVerification(), api.InsecureUsePlaintext())
	user, _, err := client.GetCurrentUser()
	if err != nil {	Fatalln(err) }

	Println("Environment is ready with user", user.Firstname, user.Lastname + ".")
}
