package reddit

import (
	"database/sql"
	"fmt"
)

type Row struct {
	Id           int          `db:"id" json:"id"`
	Path         string       `db:"path" json:"path"`
	Scheduled_At sql.NullTime `db:"scheduled_at" json:"scheduled_at"`
	Created_At   sql.NullTime `db:"created_at" json:"created_at"`
	Posted_At    sql.NullTime `db:"posted_at" json:"posted_at"`
}

func checkPosts(db *sql.DB) {
	if db != nil {
		rows, err := db.Query("select * from `to_post` where `scheduled_at` < now() and `posted_at` is null")

		if err != nil {
			return
		}

		defer rows.Close()

		for rows.Next() {
			var row Row

			if err := rows.Scan(&row.Id, &row.Path, &row.Scheduled_At, &row.Created_At, &row.Posted_At); err != nil {
				fmt.Println(err)
			}

			fmt.Println(row.Path)
		}

		err = rows.Close()

		if err != nil {
			return
		}

		if err = rows.Err(); err != nil {
			fmt.Println(err)
		}

		fmt.Println("[reddit] Opened database successfully")
	}
}
