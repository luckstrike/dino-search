-- migrations/001_create_tables.sql

-- +goose Up
CREATE TABLE IF NOT EXISTS scraped_content (
    id SERIAL PRIMARY KEY,
    url TEXT NOT NULL UNIQUE,
    title TEXT,
    main_text TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS content_headlines (
    id SERIAL PRIMARY KEY,
    content_id INTEGER REFERENCES scraped_content(id) ON DELETE CASCADE,
    headline TEXT NOT NULL,
    headline_order INTEGER NOT NULL,
    UNIQUE (content_id, headline_order)
);

CREATE TABLE IF NOT EXISTS content_keywords (
    id SERIAL PRIMARY KEY,
    content_id INTEGER REFERENCES scraped_content(id) ON DELETE CASCADE,
    keyword TEXT NOT NULL,
    UNIQUE (content_id, keyword)
);

CREATE INDEX IF NOT EXISTS idx_scraped_content_text ON scraped_content USING gin(to_tsvector('english', main_text));
CREATE INDEX IF NOT EXISTS idx_scraped_content_title ON scraped_content USING gin(to_tsvector('english', title));

-- +goose Down
DROP TABLE IF EXISTS content_keywords;
DROP TABLE IF EXISTS content_headlines;
DROP TABLE IF EXISTS scraped_content;
