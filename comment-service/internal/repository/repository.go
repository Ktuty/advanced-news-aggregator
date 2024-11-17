package repository

import (
	"comment-service/internal/models"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Comment interface {
	AddComment(comment models.Comment) error
	CommentByNewsID(newsID int) ([]models.Comment, error)
}

type Repository struct {
	Comment
}

// Создание экземпляра репозитория
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		Comment: NewCommentPostgres(db),
	}
}
