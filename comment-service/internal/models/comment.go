package models

type Comment struct {
	ID              int    `json:"id"`
	NewsID          int    `json:"news_id"`
	ParentCommentID int    `json:"parent_comment_id"`
	Content         string `json:"content"`
}
