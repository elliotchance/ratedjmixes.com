package ratedjmixes

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const UserID = 1 // FIXME

var schema = `
CREATE TABLE IF NOT EXISTS artist (
	artist_id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT,
	tl_id TEXT

	-- links, last_updated
);

CREATE TABLE IF NOT EXISTS track (
	track_id INTEGER PRIMARY KEY AUTOINCREMENT,
	title TEXT,
	artist_id INTEGER NOT NULL,
	tl_id TEXT

	-- label, links?, last_updated
);

-- SOURCE

CREATE TABLE IF NOT EXISTS tracklist (
	tracklist_id INTEGER PRIMARY KEY AUTOINCREMENT,
	tl_id TEXT,
	artist_id INTEGER NOT NULL,
	title TEXT,
	date DATE,
	episode INT,
	UNIQUE(tl_id)

	-- track_total, tracks_id, genres, duration, duration_estimate, shortlink?
	-- media links, source, image, last_updated
);

CREATE TABLE IF NOT EXISTS tracklist_track (
	tracklist_id INTEGER NOT NULL,
	track_id INTEGER NOT NULL,
	number INTEGER NOT NULL,
	UNIQUE(tracklist_id, number, track_id)

	-- time_code
);

CREATE TABLE IF NOT EXISTS users (
	user_id INTEGER PRIMARY KEY AUTOINCREMENT,
	username TEXT,
	email TEXT,
	UNIQUE(username),
	UNIQUE(email)
);

CREATE TABLE IF NOT EXISTS tracklist_collection (
	tracklist_collection_id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id INTEGER NOT NULL,
	tracklist_id INTEGER NOT NULL,
	rating SMALLINT,
	UNIQUE(user_id, tracklist_id)

	-- calculated_rating/score, notes, tags
);

-- track_collection

-- artist_collection

INSERT OR IGNORE INTO users (username, email) VALUES ('elliot', 'elliotchance@gmail.com');
`

func OpenDB() (*sqlx.DB, error) {
	db, err := sqlx.Connect("sqlite3", "db.sqlite3")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(schema)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func SaveArtist(db *sqlx.DB, artist *Artist) (*Artist, error) {
	existingArtist := &Artist{}
	err := db.Get(existingArtist, "SELECT * FROM artist WHERE name=$1", artist.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			_, err = db.Exec("INSERT INTO artist (name, tl_id) VALUES ($1, '')", artist.Name)
			if err != nil {
				return nil, err
			}

			return SaveArtist(db, artist)
		}

		return nil, err
	}

	return existingArtist, nil
}

func GetTracklist(db *sqlx.DB, tlID string) (*Tracklist, error) {
	tracklist := &Tracklist{}
	err := db.Get(tracklist, "SELECT * FROM tracklist WHERE tl_id=$1", tlID)
	if err != nil {
		return nil, err
	}

	err = db.Select(&tracklist.Tracks, `
	SELECT track.*, artist.name as artist_name
	FROM tracklist_track
	JOIN track USING (track_id)
	JOIN artist USING (artist_id)
	WHERE tracklist_id = $1
	ORDER BY number`,
		tracklist.ID)
	if err != nil {
		return nil, err
	}

	return tracklist, nil
}

func GetTrack(db *sqlx.DB, artistID int, title string) (*Track, error) {
	track := &Track{}
	err := db.Get(track, "SELECT * FROM track WHERE artist_id = $1 AND title = $2",
		artistID, title)
	if err != nil {
		return nil, err
	}

	return track, nil
}

func SaveTracklist(db *sqlx.DB, tracklist *Tracklist) (*Tracklist, error) {
	artist, err := SaveArtist(db, &Artist{
		Name: tracklist.ArtistName,
	})
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("INSERT OR IGNORE INTO tracklist (tl_id, artist_id, title, date, episode) VALUES ($1, $2, $3, $4, $5)",
		tracklist.TlID, artist.ID, tracklist.Title, tracklist.Date, tracklist.Episode)
	if err != nil {
		return nil, err
	}

	tl, err := GetTracklist(db, tracklist.TlID)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("DELETE FROM tracklist_track WHERE tracklist_id = $1", tl.ID)
	if err != nil {
		return nil, err
	}

	for i, track := range tracklist.Tracks {
		newTrack, err := SaveTrack(db, track)
		if err != nil {
			return nil, err
		}

		_, err = db.Exec("INSERT INTO tracklist_track (tracklist_id, number, track_id) VALUES ($1, $2, $3)",
			tl.ID, i+1, newTrack.ID)
		if err != nil {
			return nil, err
		}
	}

	return tl, nil
}

func SaveTrack(db *sqlx.DB, track *Track) (*Track, error) {
	artist, err := SaveArtist(db, &Artist{
		Name: track.ArtistName,
	})
	if err != nil {
		return nil, err
	}

	newTrack := &Track{}
	err = db.Get(newTrack, "SELECT * FROM track WHERE artist_id = $1 AND title = $2",
		artist.ID, track.Title)
	if err != nil {
		if err == sql.ErrNoRows {
			_, err = db.Exec("INSERT INTO track (artist_id, title, tl_id) VALUES ($1, $2, '')",
				artist.ID, track.Title)
			if err != nil {
				return nil, err
			}

			return GetTrack(db, artist.ID, track.Title)
		}

		return nil, err
	}

	return GetTrack(db, artist.ID, track.Title)
}

func AddTracklistToCollection(db *sqlx.DB, userId, tracklistID int) error {
	_, err := db.Exec("INSERT OR IGNORE INTO tracklist_collection (user_id, tracklist_id) VALUES ($1, $2)",
		userId, tracklistID)
	if err != nil {
		return err
	}

	return nil
}

func SetTracklistRating(db *sqlx.DB, userId, tracklistID int, rating float64) error {
	err := AddTracklistToCollection(db, userId, tracklistID)
	if err != nil {
		return err
	}

	_, err = db.Exec("UPDATE tracklist_collection SET rating = $1 WHERE user_id = $2 AND tracklist_id = $3",
		rating, userId, tracklistID)
	if err != nil {
		return err
	}

	return nil
}

func GetUserMixes(db *sqlx.DB, userID int) ([]*Tracklist, error) {
	var tracklists []*Tracklist
	err := db.Select(&tracklists, `
	SELECT tracklist.*, artist.name as artist_name
	FROM tracklist
	JOIN tracklist_collection USING (tracklist_id)
	JOIN artist USING (artist_id)
	WHERE user_id = $1`,
		userID)
	if err != nil {
		return nil, err
	}

	return tracklists, nil
}
