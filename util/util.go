package util

import (
	"encoding/json"
	. "fmt"
	"os"
)

// :/
func Check(err error) {
	if err != nil {
		Println(err)
		os.Exit(1)
	}
}

func Format(x interface{}) string {
	y, err := json.MarshalIndent(x, "", "\t")
	if err != nil {
		panic(err)
	}
	return string(y)
}

func PrintFormat(x interface{}) {
	y, err := json.MarshalIndent(x, "", "\t")
	if err != nil {
		panic(err)
	}
	Println(string(y))
}
