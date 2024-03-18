CREATE TABLE IF EXISTS comments (
  comment_id SERIAL PRIMARY KEY NOT NULL,
  post_id INTEGER NOT NULL,
  user_id INTEGER NOT NULL,
  comment TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);