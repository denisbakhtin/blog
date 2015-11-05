package models

import (
	"time"

	"gopkg.in/guregu/null.v3"
)

//Comment type contains post comments
type Comment struct {
	ID          int64     `json:"id" db:"id"`
	PostID      int64     `json:"post_id" db:"id"`
	ParentID    null.Int  `json:"parent_id" db:"parent_id"`
	AuthorName  string    `json:"name" db:"author_name"`
	Description string    `json:"description"`
	Published   bool      `json:"published"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	//calculated fields
	Children []Comment `json:"children" db:"-"`
}

//Insert stores Comment  in db
func (comment *Comment) Insert() error {
	err := db.QueryRow(
		`INSERT INTO comments(post_id, parent_id, author_name, description, published, created_at, updated_at) 
		VALUES($1,$2,$3,$4,$5,$6,$6) RETURNING id`,
		comment.PostID,
		comment.ParentID,
		comment.AuthorName,
		comment.Description,
		comment.Published,
		time.Now(),
	).Scan(&comment.ID)
	return err
}

//Update updates Comment record in db
func (comment *Comment) Update() error {
	_, err := db.Exec(
		`UPDATE comments 
		SET description=$2, published=$3, updated_at=$4 
		WHERE id=$1`,
		comment.ID,
		comment.Description,
		comment.Published,
		time.Now(),
	)
	return err
}

//Delete removes Comment from db.
func (comment *Comment) Delete() error {
	_, err := db.Exec("DELETE FROM comments WHERE id=$1", comment.ID)
	return err
}

//GetComment returns Comment record by its ID.
func GetComment(id interface{}) (*Comment, error) {
	comment := &Comment{}
	err := db.Get(comment, "SELECT * FROM comments WHERE id=$1", id)
	return comment, err
}

//GetComments returns a slice of comments
func GetComments() ([]Comment, error) {
	var list []Comment
	err := db.Select(&list, "SELECT * FROM comments ORDER BY comments.id DESC")
	return list, err
}

//GetPublishedComments returns a slice published of comments
func GetPublishedComments() ([]Comment, error) {
	var list []Comment
	err := db.Select(&list, "SELECT * FROM comments WHERE published=$1 ORDER BY comments.id DESC", true)
	return list, err
}

//GetCommentsByPostID returns a slice of published comments, associated with given post id
func GetCommentsByPostID(postID int64) ([]Comment, error) {
	var list []Comment
	//TODO: load child comments
	err := db.Select(
		&list,
		`SELECT * FROM comments 
		WHERE published=$1 AND post_id=$2 
		ORDER BY created_at DESC`,
		true,
		postID,
	)
	return list, err
}
