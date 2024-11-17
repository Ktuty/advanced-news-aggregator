package repository

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"news/internal/models"
)

type PostPostgres struct {
	db *pgxpool.Pool
}

// конструктор для создания экземпляра БД
func NewPostPostgres(db *pgxpool.Pool) *PostPostgres {
	return &PostPostgres{db: db}
}

// Добавление поста в БД
func (p *PostPostgres) CreatePost(ctx context.Context, post models.Post) error {
	requestID := ctx.Value("X-Request-ID")
	log.Printf("Request ID: %s, Creating post: %v", requestID, post)

	_, err := p.db.Exec(ctx, `
	INSERT INTO posts (title, content, published_at, link) VALUES ($1, $2, $3, $4) `,
		post.Title, post.Content, post.PubTime, post.Link)
	if err != nil {
		log.Printf("Request ID: %s, Error inserting post: %v", requestID, err)
		return err
	}
	return nil
}

// Получение постов из БД с постраничной навигацией
func (p *PostPostgres) Posts(ctx context.Context, page, pageSize, postID int) ([]models.Post, int, error) {
	requestID := ctx.Value("X-Request-ID")
	log.Printf("Request ID: %s, Getting posts", requestID)

	var query string
	var args []interface{}

	if postID != 0 {
		query = `SELECT id, title, content, published_at, link FROM posts WHERE id = $1`
		args = append(args, postID)
	} else {
		offset := (page - 1) * pageSize
		query = `SELECT id, title, content, published_at, link FROM posts ORDER BY published_at DESC LIMIT $1 OFFSET $2`
		args = append(args, pageSize, offset)
	}

	rows, err := p.db.Query(ctx, query, args...)
	if err != nil {
		log.Printf("Request ID: %s, Error getting posts: %v", requestID, err)
		return nil, 0, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.Id, &post.Title, &post.Content, &post.PubTime, &post.Link); err != nil {
			log.Printf("Request ID: %s, Error scanning posts: %v", requestID, err)
			return nil, 0, err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Request ID: %s, Error iterating over rows: %v", requestID, err)
		return nil, 0, err
	}

	var totalRows int
	if err := p.db.QueryRow(ctx, `SELECT COUNT(*) FROM posts`).Scan(&totalRows); err != nil {
		log.Printf("Request ID: %s, Error counting posts: %v", requestID, err)
		return nil, 0, err
	}

	totalPages := (totalRows + pageSize - 1) / pageSize
	return posts, totalPages, nil
}

// Поиск постов по названию с постраничной навигацией
func (p *PostPostgres) SearchPostsByTitle(ctx context.Context, title string, page, pageSize int) ([]models.Post, int, error) {
	requestID := ctx.Value("X-Request-ID")
	log.Printf("Request ID: %s, Searching posts by title: %s", requestID, title)

	offset := (page - 1) * pageSize

	rows, err := p.db.Query(ctx, `
	SELECT id, title, content, published_at, link FROM posts WHERE title ILIKE $1 ORDER BY published_at DESC LIMIT $2 OFFSET $3`, "%"+title+"%", pageSize, offset)
	if err != nil {
		log.Printf("Request ID: %s, Error searching posts by title: %v", requestID, err)
		return nil, 0, err
	}

	defer rows.Close()

	var posts []models.Post

	for rows.Next() {
		var post models.Post
		err = rows.Scan(&post.Id, &post.Title, &post.Content, &post.PubTime, &post.Link)
		if err != nil {
			log.Printf("Request ID: %s, Error searching posts by title: %v", requestID, err)
			return nil, 0, err
		}
		posts = append(posts, post)
	}

	var totalRows int
	err = p.db.QueryRow(ctx, `SELECT COUNT(*) FROM posts WHERE title ILIKE $1`, "%"+title+"%").Scan(&totalRows)
	if err != nil {
		log.Printf("Request ID: %s, Error counting posts by title: %v", requestID, err)
		return nil, 0, err
	}

	totalPages := (totalRows + pageSize - 1) / pageSize

	return posts, totalPages, nil
}
