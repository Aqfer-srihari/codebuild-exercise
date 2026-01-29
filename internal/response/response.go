package response

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Status int
	Data   any
	Err    error
}

func Write(w http.ResponseWriter, resp APIResponse) {
	if resp.Err != nil {
		http.Error(w, resp.Err.Error(), resp.Status)
		return
	}

	if resp.Data != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.Status)
		json.NewEncoder(w).Encode(resp.Data)
		return
	}

	w.WriteHeader(resp.Status)
}
