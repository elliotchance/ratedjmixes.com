package main

import (
	"github.com/elliotchance/ratedjmixes.com/ratedjmixes"
)

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
