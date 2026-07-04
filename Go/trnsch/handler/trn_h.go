package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"trnsch/db"
	"trnsch/model"
	"trnsch/trn" // Adjust this import path based on your actual module name
)

func BuildTrnById(w http.ResponseWriter, req *http.Request) {
	var payload model.TrnBuild
	defer req.Body.Close()

	err := json.NewDecoder(req.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err = trn.Start(payload.Id, payload.Build); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // for CORS if needed
	w.WriteHeader(http.StatusOK)
	return
}

func SchResByRnd(w http.ResponseWriter, req *http.Request) {
	var payload model.TrnSch
	defer req.Body.Close()

	err := json.NewDecoder(req.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var jsonOut []byte
	if jsonOut, err = trn.BuildSchres(payload.Id, payload.Rnd); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // for CORS if needed
	w.WriteHeader(http.StatusOK)
	w.Write(jsonOut)
	return
}

func GetSch(w http.ResponseWriter, req *http.Request) {
	var payload model.TrnSch
	defer req.Body.Close()

	err := json.NewDecoder(req.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	var jsonOut []byte
	if jsonOut, err = trn.GetRnd(payload.Id, payload.Rnd); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // for CORS if needed
	w.WriteHeader(http.StatusOK)
	w.Write(jsonOut)
	return
}

func SwapOrdr(w http.ResponseWriter, req *http.Request) {
	var payload model.SchSwap
	defer req.Body.Close()

	err := json.NewDecoder(req.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var jsonOut []byte
	if jsonOut, err = trn.SchSwap(payload.TrnId, payload.Rnd, payload.Gid, payload.Up); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // for CORS if needed
	w.WriteHeader(http.StatusOK)
	w.Write(jsonOut)
	return
}

func EditSch(w http.ResponseWriter, req *http.Request) {
	var payload model.SchEdit
	defer req.Body.Close()

	err := json.NewDecoder(req.Body).Decode(&payload)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var jsonOut []byte
	if jsonOut, err = trn.UpdateRnd(payload.Gid, payload.Plyrwn, payload.Pf, payload.Pa, payload.Notes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // for CORS if needed
	w.WriteHeader(http.StatusOK)
	w.Write(jsonOut)
	return
}

func AddDeleteMsc(w http.ResponseWriter, req *http.Request) {
	var payload model.AddDeleteMsc
	defer req.Body.Close()

	err := json.NewDecoder(req.Body).Decode(&payload)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var jsonOut []byte
	if jsonOut, err = trn.AddDeleteMsc(payload.TrnId, payload.Rnd, payload.Plyrtm, payload.Letters, payload.AddDelete); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // for CORS if needed
	w.WriteHeader(http.StatusOK)
	w.Write(jsonOut)
	return
}

func GetTrnById(w http.ResponseWriter, req *http.Request) {
	var payload model.Trn
	defer req.Body.Close()

	err := json.NewDecoder(req.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var jsonOut []byte
	if jsonOut, err = trn.Get(payload.Id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // for CORS if needed
	w.WriteHeader(http.StatusOK)
	w.Write(jsonOut)
	return
}

func GetLeaguesByParam(w http.ResponseWriter, req *http.Request) {
	var payload model.League
	defer req.Body.Close()

	err := json.NewDecoder(req.Body).Decode(&payload)
	if err != nil || payload.Description == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	rows, err := db.DB.Query("SELECT id, description FROM league where id = ?", payload.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var leagues []model.League
	for rows.Next() {
		var leag model.League
		if err := rows.Scan(&leag.Id, &leag.Description); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		leagues = append(leagues, leag)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(leagues)
}

func GetLeagues(w http.ResponseWriter) {
	rows, err := db.DB.Query("SELECT id, description FROM league order by description")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var leagues []model.League
	for rows.Next() {
		var leag model.League
		if err := rows.Scan(&leag.Id, &leag.Description); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		leagues = append(leagues, leag)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(leagues)
}

func Import(w http.ResponseWriter, req *http.Request) {
	var payload model.Import
	defer req.Body.Close()

	err := json.NewDecoder(req.Body).Decode(&payload)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var jsonOut []byte
	if jsonOut, err = trn.Import(payload.TrnId, payload.Filename); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // for CORS if needed
	w.WriteHeader(http.StatusOK)
	w.Write(jsonOut)
	return
}
