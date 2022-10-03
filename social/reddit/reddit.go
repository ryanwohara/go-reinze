package reddit

import (
	"database/sql"
)

func CheckPosts(db *sql.DB) {
	checkPosts(db)
}
