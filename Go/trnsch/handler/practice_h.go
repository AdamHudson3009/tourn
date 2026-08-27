package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"trnsch/model"
)

func CreateUserHandler(w http.ResponseWriter, req *http.Request) {
	var payload model.TrnBuild
	defer req.Body.Close()

	err := json.NewDecoder(req.Body).Decode(&payload)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid input %s", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"message": fmt.Sprintf("Id retrieved %v created", payload.Id)}
	json.NewEncoder(w).Encode(response)
}
