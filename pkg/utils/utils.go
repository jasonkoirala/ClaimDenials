package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func ParseRequestBody(r *http.Request, X interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, X); err != nil {
		return err
	}
	return nil
}
