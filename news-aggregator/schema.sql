DROP TABLE IF EXISTS posts;

CREATE TABLE posts (
                       id SERIAL PRIMARY KEY,
                       title TEXT NOT NULL,
                       content TEXT NOT NULL,
                       published_at INTEGER DEFAULT 0,
                       link TEXT NOT NULL UNIQUE
);

CREATE INDEX idx_posts_title ON posts (title);