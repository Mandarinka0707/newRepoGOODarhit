package entity

import "time"

type Post struct {
	ID        int64     `db:"id"`
	Title     string    `db:"title"`
	Content   string    `db:"content"`
	AuthorID  int64     `db:"author_id"`
	CreatedAt time.Time `db:"created_at"`
}
