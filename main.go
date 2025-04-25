package main

import (
	"github.com/elliotchance/ratedjmixes.com/ratedjmixes"
)

func loadTracklist() error {
	db, err := ratedjmixes.OpenDB()
	if err != nil {
		return err
	}

	tl, err := ratedjmixes.ParseTracklist("ratedjmixes/test_data/tiesto-tiestos-club-life-942-2025-04-19.html", "1jyrlttk")
	if err != nil {
		return err
	}

	tl, err = ratedjmixes.SaveTracklist(db, tl)
	if err != nil {
		return err
	}

	userId := 1 // FIXME
	err = ratedjmixes.SetTracklistRating(db, userId, tl.ID, 3.5)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	db, err := ratedjmixes.OpenDB()
	if err != nil {
		panic(db)
	}

	ratedjmixes.StartServer(db)

	// err := loadTracklist()
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Printf("%+#v\n", tl)
}
