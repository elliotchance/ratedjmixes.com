package ratedjmixes

import (
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/jmoiron/sqlx"
)

var tlURLRe = regexp.MustCompile(`tracklist/([^/$]+)`)

func handleSearch(db *sqlx.DB, w http.ResponseWriter, req *http.Request) error {
	url := req.URL.Query().Get("url")
	matches := tlURLRe.FindAllStringSubmatch(url, -1)
	tlID := strings.TrimSpace(matches[0][1])

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.182 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	tl, err := ParseTracklist(dat, tlID)
	if err != nil {
		return err
	}

	tl, err = SaveTracklist(db, tl)
	if err != nil {
		return err
	}

	http.Redirect(w, req, "/mix/"+tl.TlID, http.StatusSeeOther)
	return nil
}
