package repository

import (
	"comment-service/internal/models"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type CommentPostgres struct {
	db *pgxpool.Pool
}

func NewCommentPostgres(db *pgxpool.Pool) *CommentPostgres {
	return &CommentPostgres{db: db}
}

func (c *CommentPostgres) AddComment(comment models.Comment) error {
	_, err := c.db.Exec(context.Background(), `
	INSERT INTO comments (news_id, parent_comment_id, content) VALUES ($1, $2, $3)`,
		comment.NewsID, comment.ParentCommentID, comment.Content)
	if err != nil {
		log.Println("Error inserting comment:", err)
	}
	return err
}

func (c *CommentPostgres) CommentByNewsID(newsID int) ([]models.Comment, error) {
	query := `SELECT id, news_id, parent_comment_id, content FROM comments WHERE news_id = $1`
	rows, err := c.db.Query(context.Background(), query, newsID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(&comment.ID, &comment.NewsID, &comment.ParentCommentID, &comment.Content)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}
