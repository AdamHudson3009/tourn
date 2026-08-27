package main

import (
	"log"
	"net/http"

	"trnsch/db"
	"trnsch/handler"
)

// withCORS adds the CORS headers and handles OPTIONS preflight
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	http.Handle("/trns/Build", withCORS(http.HandlerFunc(handler.BuildTrnById)))
	http.Handle("/trns/Get", withCORS(http.HandlerFunc(handler.GetTrnById)))
	http.Handle("/trns/GetSch", withCORS(http.HandlerFunc(handler.GetSch)))
	http.Handle("/trns/SwapOrdr", withCORS(http.HandlerFunc(handler.SwapOrdr)))
	http.Handle("/trns/EditSch", withCORS(http.HandlerFunc(handler.EditSch)))
	http.Handle("/trns/SchResByRnd", withCORS(http.HandlerFunc(handler.SchResByRnd)))
	http.Handle("/trns/AddDeleteMsc", withCORS(http.HandlerFunc(handler.AddDeleteMsc)))
	http.Handle("/trns/OrdNum", withCORS(http.HandlerFunc(handler.OrdNum)))
	http.Handle("/trns/Potw", withCORS(http.HandlerFunc(handler.Potw)))
	http.Handle("/trns/Import", withCORS(http.HandlerFunc(handler.Import)))
	http.Handle("/trns/Build/", withCORS(http.HandlerFunc(handler.BuildTrnById)))
	http.Handle("/trns/Get/", withCORS(http.HandlerFunc(handler.GetTrnById)))
	http.Handle("/trns/GetSch/", withCORS(http.HandlerFunc(handler.GetSch)))
	http.Handle("/trns/SwapOrdr/", withCORS(http.HandlerFunc(handler.SwapOrdr)))
	http.Handle("/trns/EditSch/", withCORS(http.HandlerFunc(handler.EditSch)))
	http.Handle("/trns/SchResByRnd/", withCORS(http.HandlerFunc(handler.SchResByRnd)))
	http.Handle("/trns/AddDeleteMsc/", withCORS(http.HandlerFunc(handler.AddDeleteMsc)))
	http.Handle("/trns/OrdNum/", withCORS(http.HandlerFunc(handler.OrdNum)))
	http.Handle("/trns/Potw/", withCORS(http.HandlerFunc(handler.Potw)))
	http.Handle("/trns/Import/", withCORS(http.HandlerFunc(handler.Import)))

	db.Init() // Initialize DB connection

	log.Fatal(http.ListenAndServe(":8080", nil)) // ← use default mux where your http.Handle calls are registered
}
