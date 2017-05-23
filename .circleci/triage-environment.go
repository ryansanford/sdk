package main

import (
	"encoding/json"
	. "fmt"
	"log"
	"time"

	"gopkg.in/mgo.v2"

	"flywheel.io/sdk/api"
)

func check(err error) {
	if err != nil { log.Fatalln(err) }
}

func main() {
	var session *mgo.Session
	var err error

	for i := 1; i <= 3; i++ {
		log.Println("Connecting to mongo...")
		wait := time.Duration(float64(i) * 5.0 * float64(time.Second))
		session, err = mgo.DialWithTimeout("localhost", wait)
		if err == nil { break }
	}

	if err != nil { log.Fatalln("Could not connect to mongo:", err) }
	defer session.Close()
	session.SetSafe(&mgo.Safe{})

	database := session.DB("scitran")
	tables, err := database.CollectionNames()
	check(err)

	Println()
	Println("There are", len(tables), "tables:", tables)

	for _, tableName := range tables {
		table := database.C(tableName)
		count, err := table.Count()
		check(err)
		if count <= 0 { continue }

		Println()
		Println("There are", count, "documents in table", tableName + ":")
		Println()

		var item interface{}
		cursor := table.Find(nil).Iter()

		for cursor.Next(&item) {
			y, err := json.MarshalIndent(item, "", "\t")
			check(err)
			Println(string(y))
		}

		Println()
		check(cursor.Close())
	}

	var client *api.Client
	var user *api.User

	for i := 1; i <= 15; i++ {
		err = nil
		Println("Checking API...")
		client = api.NewApiKeyClient("localhost:8080:insecure-key", api.InsecureNoSSLVerification, api.InsecureUsePlaintext)
		user, _, err = client.GetCurrentUser()
		if err == nil { break }
		time.Sleep(1000 * time.Millisecond)
	}
	if err != nil {	log.Fatalln("Could not connect to API:", err) }

	Println("Environment still up with user", user.Firstname, user.Lastname + ".")
}
