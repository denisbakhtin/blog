-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE INDEX posts_search_idx ON posts USING gin(to_tsvector('english', name || ' ' || description));

-- +migrate Down
-- SQL in section 'Down' is executed when this migration is rolled back
DROP INDEX posts_search_idx;


