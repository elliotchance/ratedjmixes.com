package ratedjmixes

import (
	"net/http"

	"github.com/jmoiron/sqlx"
)

func handleAddMix(db *sqlx.DB, w http.ResponseWriter, req *http.Request) error {
	tl, err := GetTracklist(db, req.PathValue("mix"))
	if err != nil {
		return err
	}

	err = AddTracklistToCollection(db, UserID, tl.ID)
	if err != nil {
		return err
	}

	http.Redirect(w, req, "/elliot/mixes", http.StatusSeeOther)
	return nil
}
