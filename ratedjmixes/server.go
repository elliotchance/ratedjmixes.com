package ratedjmixes

import (
	"net/http"

	"github.com/jmoiron/sqlx"
)

func handler(db *sqlx.DB, fn func(db *sqlx.DB, w http.ResponseWriter, r *http.Request) error) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := fn(db, w, r)
		if err != nil {
			panic(err)
		}
	}
}

func StartServer(db *sqlx.DB) {
	http.HandleFunc("/elliot/mixes", handler(db, handleMyMixes))
	http.HandleFunc("/mix/{mix}", handler(db, handleViewMix))
	http.HandleFunc("/add/{mix}", handler(db, handleAddMix))
	http.HandleFunc("/search", handler(db, handleSearch))
	http.HandleFunc("/", handler(db, handleIndex))
	http.ListenAndServe(":8080", nil)
}
