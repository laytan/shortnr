package responder

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Res is a valid response
type Res struct {
	Code    int
	Message string
	Data    interface{}
}

// MarshalJSON turns the struct into a error we can return outside safely
func (r Res) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{})
	m["err"] = nil
	m["res"] = make(map[string]interface{})

	m["res"].(map[string]interface{})["data"] = r.Data
	m["res"].(map[string]interface{})["msg"] = r.Message

	if len(m["res"].(map[string]interface{})["msg"].(string)) < 1 {
		m["res"].(map[string]interface{})["msg"] = "succes"
	}

	return json.Marshal(m)
}

// Send sends the error over the responsewriter
func (r Res) Send(w http.ResponseWriter) {
	if r.Code != 0 {
		w.WriteHeader(r.Code)
	}
	json, err := json.Marshal(r)
	if err != nil {
		fmt.Printf("error marshalling request error: %v \n", err)
		w.WriteHeader(500)
		return
	}

	w.Write(json)
}
