package request

import (
	"encoding/json"
	"net/http"
)

func ReadBody(r *http.Request, dst interface{}) error {

	err := json.NewDecoder(r.Body).Decode(dst)

	if err != nil {
		return err
	}

	return nil
}
