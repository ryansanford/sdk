package api

import (
	"encoding/json"
	"errors"
	. "fmt"
	"strconv"
)

// Coalesce will extract an API error message into a golang error, if applicable.
func Coalesce(err error, aerr *Error) error {

	// This should eventually be replaced with something that takes the http response as well and looks at that. Expect200? Etc.

	if err != nil {
		return err
	} else if aerr != nil {
		if aerr.Message == "" {
			aerr.Message = "Unknown server error"
		}
		aerr.Message = "(" + strconv.Itoa(aerr.StatusCode) + ") " + aerr.Message
		return errors.New(aerr.Message)
	} else {
		return nil
	}
}

// Convenience functions for development

func Format(x interface{}) string {
	y, err := json.MarshalIndent(x, "", "\t")
	if err != nil {
		panic(err)
	}
	return string(y)
}

func PrintFormat(x interface{}) {
	Println(Format(x))
}
