package responder

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Err is an error
type Err struct {
	Code int
	Err  error
}

func (r Err) Error() string {
	return r.Err.Error()
}

// MarshalJSON turns the struct into a error we can return outside safely
func (r Err) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{})
	m["err"] = make(map[string]interface{})
	m["res"] = nil

	m["err"].(map[string]interface{})["msg"] = r.Err.Error()
	if r.Code >= 500 && r.Code < 600 {
		m["err"].(map[string]interface{})["msg"] = "internal server error"
	}

	if len(m["err"].(map[string]interface{})["msg"].(string)) < 1 {
		m["err"].(map[string]interface{})["msg"] = "something went wrong, please try again later"
	}

	return json.Marshal(m)
}

// Send sends the error over the responsewriter
func (r Err) Send(w http.ResponseWriter) {
	fmt.Println(r)
	w.WriteHeader(r.Code)
	json, err := json.Marshal(r)
	if err != nil {
		fmt.Printf("error marshalling request error: %v \n", err)
		w.WriteHeader(500)
		return
	}
	w.Write(json)
}

// CastAndSend sends a requestErr back
func CastAndSend(e error, w http.ResponseWriter) {
	if rErr, ok := e.(Err); ok {
		rErr.Send(w)
	} else {
		fmt.Printf("error is not a request error %v \n", e)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
