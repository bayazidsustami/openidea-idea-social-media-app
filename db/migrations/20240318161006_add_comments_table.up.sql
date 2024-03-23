CREATE TABLE IF NOT EXISTS comments (
  comment_id SERIAL PRIMARY KEY NOT NULL,
  post_id INTEGER NOT NULL,
  user_id UUID NOT NULL,
  comment TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);