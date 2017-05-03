package main

import (
	"C"
	"encoding/json"
	"strings"

	"flywheel.io/sdk/api"
)

// Stub to keep go happy. Ignored in c-shared mode.
func main() {}

// CallResult holds the response to any call to the C bridge.
// Most calls will return a JSON marshaled CallResult.
type CallResult struct {

	// Success specifies whether the call succeeded.
	Success bool `json:"success"`

	// Message contains an error message. Valid IFF success is false.
	Message string `json:"message"`

	// Data contains the result of the call. Can be null.
	Data interface{} `json:"data"`
}

func makeClient(key *C.char) *api.Client {
	apiKey := C.GoString(key)
	splits := strings.Split(apiKey, ":")

	// If port number is specified, it's a non-production key; disable SSL verification.
	// Otherwise, require it.
	if len(splits) == 2 {
		return api.NewApiKeyClient(apiKey)
	} else if len(splits) > 2 {
		return api.NewApiKeyClient(apiKey, api.InsecureNoSSLVerification())
	} else {
		return nil
	}
}

func handleError(err error, status *C.int) *C.char {
	result := CallResult{
		Success: false,
		Message: err.Error(),
		Data:    nil,
	}

	// Error ignored because no unknown types to marshal
	raw, _ := json.Marshal(result)
	return C.CString(string(raw))
}

func handleSuccess(data interface{}, status *C.int) *C.char {
	result := CallResult{
		Success: true,
		Message: "",
		Data:    data,
	}

	raw, encodeErr := json.Marshal(result)

	// Should never happen; if triggered, the helper was called with bad data
	if encodeErr != nil {
		*status = -1
		return handleError(encodeErr, status)
	}

	return C.CString(string(raw))
}

// Given a normal API result, set a success pointer and return either the data or the error contents.
func format(data interface{}, err error, status *C.int) *C.char {

	if err != nil {
		*status = -1
		return handleError(err, status)
	}

	*status = 0
	return handleSuccess(data, status)
}

//export Free
func Free(pointer *C.char) {
	// C.free(unsafe.Pointer(pointer))
}

//export TestBridge
func TestBridge(name *C.char) *C.char {
	nameGo := C.GoString(name)
	return C.CString("Hello " + nameGo)
}

//
// -- AUTO GENERATED CODE FOLLOWS --
//
