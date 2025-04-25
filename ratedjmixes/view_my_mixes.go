package ratedjmixes

import (
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"
)

func handleMyMixes(db *sqlx.DB, w http.ResponseWriter, req *http.Request) error {
	fmt.Fprintf(w, "<body>")

	err := viewHeader(db, w, req, headerOptions{})
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "<table border=1>\n")
	fmt.Fprintf(w, "<tr>\n")
	fmt.Fprintf(w, "<th>Artist</th>\n")
	fmt.Fprintf(w, "<th>Title</th>\n")
	fmt.Fprintf(w, "<th>Date</th>\n")
	fmt.Fprintf(w, "<th>Rating</th>\n")
	fmt.Fprintf(w, "</tr>\n")

	mixes, err := GetUserMixes(db, UserID)
	if err != nil {
		return err
	}

	for _, mix := range mixes {
		fmt.Fprintf(w, "<tr>\n")
		fmt.Fprintf(w, "<td>%s</td>\n", mix.ArtistName)
		fmt.Fprintf(w, "<td>%s</td>\n", mix.Title)
		fmt.Fprintf(w, "<td>%s</td>\n", mix.Date[:10])
		fmt.Fprintf(w, "<td>%s</td>\n", mix.RatingString())
		fmt.Fprintf(w, "</tr>\n")
	}

	fmt.Fprintf(w, "</table>\n")
	fmt.Fprintf(w, "</body>")

	return nil
}
