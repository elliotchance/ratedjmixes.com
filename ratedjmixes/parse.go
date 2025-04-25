package ratedjmixes

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Artist struct {
	ID   int    `db:"artist_id"`
	Name string `db:"name"`
	TlID string `db:"tl_id"`
}

type Track struct {
	ID       int    `db:"track_id"`
	ArtistID int    `db:"artist_id"`
	Title    string `db:"title"`
	TlID     string `db:"tl_id"`

	// Fake fields
	ArtistName string `db:"artist_name"`
}

type Tracklist struct {
	ID       int    `db:"tracklist_id"`
	ArtistID int    `db:"artist_id"`
	Title    string `db:"title"`
	Date     string `db:"date"`
	Episode  int    `db:"episode"`
	TlID     string `db:"tl_id"`

	// Fake fields
	ArtistName string  `db:"artist_name"`
	Rating     float64 `db:"rating"`
	Tracks     []*Track
}

func (m *Tracklist) RatingString() string {
	if m.Rating == 0 {
		return ""
	}

	return fmt.Sprintf("%.1f", m.Rating)
}

type TracklistCollection struct {
	ID          int     `db:"tracklist_collection_id"`
	UserID      int     `db:"user_id"`
	TracklistID int     `db:"tracklist_id"`
	Rating      float64 `db:"rating"`
}

var titleRe = regexp.MustCompile(`(.*) - (.*) ([\d-]{10})`)
var episodeRe = regexp.MustCompile(` (\d{3,})`)
var trackRe = regexp.MustCompile(`(.*) - (.*)`)

func ParseTracklist(dat []byte, tlID string) (*Tracklist, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(dat))
	if err != nil {
		return nil, err
	}

	tl := &Tracklist{
		TlID: tlID,
	}

	// Find the review items
	doc.Find("#pageTitle").Each(func(i int, s *goquery.Selection) {
		matches := titleRe.FindAllStringSubmatch(s.Text(), -1)
		tl.ArtistName = strings.TrimSpace(matches[0][1])
		tl.Title = strings.TrimSpace(matches[0][2])
		tl.Date = strings.TrimSpace(matches[0][3])

		matches = episodeRe.FindAllStringSubmatch(tl.Title, -1)
		if len(matches) > 0 {
			tl.Episode, _ = strconv.Atoi(matches[0][1])
		}
	})

	doc.Find(".bCont.tl meta[itemprop='name']").Each(func(i int, s *goquery.Selection) {
		matches := trackRe.FindAllStringSubmatch(s.AttrOr("content", ""), -1)
		tl.Tracks = append(tl.Tracks, &Track{
			ArtistName: strings.TrimSpace(matches[0][1]),
			Title:      strings.TrimSpace(matches[0][2]),
		})
	})

	return tl, nil
}
