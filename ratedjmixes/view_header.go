package ratedjmixes

import (
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"
)

type headerOptions struct {
	TlID string
}

func viewHeader(db *sqlx.DB, w http.ResponseWriter, req *http.Request, options headerOptions) error {
	fmt.Fprintf(w, `<form method="GET" action="/search">`)
	fmt.Fprintf(w, "<table>\n")
	fmt.Fprintf(w, "<tr>\n")
	fmt.Fprintf(w, `<td><input type="text" name="url"/></td>`)
	fmt.Fprintf(w, `<td><input type="submit" /></td>`)
	fmt.Fprintf(w, `<td><a href="/elliot/mixes">My Collection</a></td>`)
	if options.TlID != "" {
		fmt.Fprintf(w, `<td><a href="/add/%s">Add to Collection</a></td>`, options.TlID)
	}
	fmt.Fprintf(w, "</tr>\n")
	fmt.Fprintf(w, "</table>\n")

	// fmt.Fprintf(w, `<input type="text" name="url"/>`)
	// fmt.Fprintf(w, `<input type="submit" />`)
	// fmt.Fprintf(w, "</form>\n")

	// fmt.Fprintf(w, `<form method="GET" action="/add">`)
	// fmt.Fprintf(w, `<input type="text" value="%s" />`, "ybc81g9")
	// fmt.Fprintf(w, `<input type="submit" value="Add to Collection" />`)
	fmt.Fprintf(w, "</form>\n")

	return nil
}
