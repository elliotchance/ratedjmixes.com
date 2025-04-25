package ratedjmixes

import (
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"
)

func handleIndex(db *sqlx.DB, w http.ResponseWriter, req *http.Request) error {
	fmt.Fprintf(w, "<table border=1>\n")
	fmt.Fprintf(w, "<tr>\n")
	fmt.Fprintf(w, "<th>#</th>\n")
	fmt.Fprintf(w, "<th>Artist</th>\n")
	fmt.Fprintf(w, "<th>Title</th>\n")
	fmt.Fprintf(w, "</tr>\n")

	tl, err := GetTracklist(db, "1jyrlttk")
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

	return nil
}

func StartServer(db *sqlx.DB) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := handleIndex(db, w, r)
		if err != nil {
			panic(err)
		}
	})
	http.ListenAndServe(":8080", nil)
}
