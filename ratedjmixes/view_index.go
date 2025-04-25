package ratedjmixes

import (
	"net/http"

	"github.com/jmoiron/sqlx"
)

func handleIndex(db *sqlx.DB, w http.ResponseWriter, req *http.Request) error {
	http.Redirect(w, req, "/elliot/mixes", http.StatusSeeOther)
	return nil
}
