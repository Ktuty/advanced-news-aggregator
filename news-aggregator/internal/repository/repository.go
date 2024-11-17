package repository

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"news/internal/models"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=Post
type Post interface {
	CreatePost(ctx context.Context, post models.Post) error
	Posts(ctx context.Context, page, pageSize, postID int) ([]models.Post, int, error)
	SearchPostsByTitle(ctx context.Context, title string, page, pageSize int) ([]models.Post, int, error)
}

type Repository struct {
	Post
}

// Создание экземпляра репозитория
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		Post: NewPostPostgres(db),
	}
}
