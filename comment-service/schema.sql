DROP TABLE IF EXISTS comments;

CREATE TABLE comments (
                          id SERIAL PRIMARY KEY,
                          news_id INT NOT NULL,
                          parent_comment_id INT,
                          content TEXT NOT NULL
);