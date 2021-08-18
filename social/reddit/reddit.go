package reddit

import (
	"database/sql"
)

func CheckPosts(db *sql.DB) {
	checkPosts(db)

	return
}

func maybePanic(err error) {
	if err != nil {
		panic(err)
	}
}
