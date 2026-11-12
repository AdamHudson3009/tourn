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

func OrdNum(w http.ResponseWriter, req *http.Request) {
	var payload model.OrdNum
	defer req.Body.Close()

	err := json.NewDecoder(req.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var jsonOut []byte
	if jsonOut, err = trn.OrdNum(payload.TrnId, payload.Rnd, payload.OldNum, payload.NewNum); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // for CORS if needed
	w.WriteHeader(http.StatusOK)
	w.Write(jsonOut)
	return
}

func Potw(w http.ResponseWriter, req *http.Request) {
	var payload model.Trn
	defer req.Body.Close()

	err := json.NewDecoder(req.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var jsonOut []byte
	if jsonOut, err = trn.Potw(payload.Id); err != nil {
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

func EditSchs(w http.ResponseWriter, req *http.Request) {
	var payload []model.SchEdit
	defer req.Body.Close()
	
	err := json.NewDecoder(req.Body).Decode(&payload)
	
	if err != nil {
		fmt.Println("Invalid input: " + err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(payload) == 0 {
		http.Error(w, "No edits provided", http.StatusBadRequest)
		return
	}

	results := make([]json.RawMessage, 0, len(payload))
	for _, p := range payload {
		jsonOut, err := trn.UpdateRnd(p.Gid, p.Plyrwn, p.Pf, p.Pa, p.Notes)
		if err != nil {
			http.Error(w, fmt.Sprintf("gid %d: %s", p.Gid, err.Error()), http.StatusInternalServerError)
			return
		}
		results = append(results, json.RawMessage(jsonOut))
	}

	var Ordernums [] int
	var Gids [] int
	call_ordnum := false
	for _, p := range payload {
		if (p.Ordnum == 0) {
			break
		}
		call_ordnum = true
		Ordernums = append(Ordernums, p.Ordnum)
		Gids = append(Gids, p.Gid)
	}
	
	for _, p := range payload {
		var jsonOut []byte
		var err error
		if (p.OldMscAddAll1 != p.MscAddAll1) {
			jsonOut, err = trn.AddRemoveMscEach(p.Gid, p.Plyr1, p.MscAddAll1, p.OldMscAddAll1)
		}
		if (p.OldMscAddAll2 != p.MscAddAll2) {
			jsonOut, err = trn.AddRemoveMscEach(p.Gid, p.Plyr2, p.MscAddAll2, p.OldMscAddAll2)
		}
		if err != nil {
			http.Error(w, fmt.Sprintf("AddRemoveMsc %s %s %s %s", p.Plyr1, p.MscAddAll1, p.OldMscAddAll1, err.Error()), http.StatusInternalServerError)
			return
		}
		results = append(results, json.RawMessage(jsonOut))
	}

	if (call_ordnum) {
		jsonOut, err := trn.UpdateOrdnums(Gids, Ordernums)
		if err != nil {
			http.Error(w, fmt.Sprintf("gid %s", Gids, Ordernums, err.Error()), http.StatusInternalServerError)
			return
		}
		results = append(results, json.RawMessage(jsonOut))
	}

	out, err := json.Marshal(results)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}



	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // for CORS if needed
	w.WriteHeader(http.StatusOK)
	w.Write(out)
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
