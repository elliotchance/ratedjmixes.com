package ratedjmixes

import (
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"
)

func handleViewMix(db *sqlx.DB, w http.ResponseWriter, req *http.Request) error {
	fmt.Fprintf(w, "<body>")

	err := viewHeader(db, w, req, headerOptions{TlID: req.PathValue("mix")})
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "<table border=1>\n")
	fmt.Fprintf(w, "<tr>\n")
	fmt.Fprintf(w, "<th>#</th>\n")
	fmt.Fprintf(w, "<th>Artist</th>\n")
	fmt.Fprintf(w, "<th>Title</th>\n")
	fmt.Fprintf(w, "</tr>\n")

	tl, err := GetTracklist(db, req.PathValue("mix"))
	if err != nil {
		return err
	}

	for i, track := range tl.Tracks {
		fmt.Fprintf(w, "<tr>\n")
		fmt.Fprintf(w, "<td>%d</td>\n", i+1)
		fmt.Fprintf(w, "<td>%s</td>\n", track.ArtistName)
		fmt.Fprintf(w, "<td>%s</td>\n", track.Title)
		fmt.Fprintf(w, "</tr>\n")
	}

	fmt.Fprintf(w, "</table>\n")
	fmt.Fprintf(w, "</body>")

	return nil
}
