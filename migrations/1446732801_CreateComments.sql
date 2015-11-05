-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE comments(
  id SERIAL PRIMARY KEY,
  post_id INTEGER REFERENCES posts(id) ON DELETE CASCADE,
  parent_id INTEGER REFERENCES comments(id) ON DELETE CASCADE,
  author_name TEXT NOT NULL,
  description TEXT NOT NULL,
  published BOOLEAN NOT NULL DEFAULT true,
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);


-- +migrate Down
-- SQL in section 'Down' is executed when this migration is rolled back
DROP TABLE comments;

